package networking

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	. "github.com/nosyliam/revolution/pkg/common"
	"github.com/nosyliam/revolution/pkg/config"
	"github.com/nosyliam/revolution/pkg/logging"
	"github.com/sqweek/dialog"
	"io"
	"net"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/grandcat/zeroconf"
)

type subscriber struct {
	ch   chan *Message
	once bool
}

type Client struct {
	stop       chan struct{}
	disconnect chan struct{}
	mu         sync.Mutex
	conn       net.Conn
	watchers   map[MessageKind]map[subscriber]bool
	state      *config.Object[config.MacroState]
	logger     *logging.Logger
}

func NewClient(state *config.Object[config.MacroState], logger *logging.Logger) *Client {
	client := &Client{
		stop:       make(chan struct{}),
		disconnect: make(chan struct{}),
		watchers:   make(map[MessageKind]map[subscriber]bool),
		state:      state,
		logger:     logger,
	}
	for _, kind := range MessageKinds {
		client.watchers[kind] = make(map[subscriber]bool)
	}
	state.SetPath("networking.identity", client.Identity())
	return client
}

func (c *Client) Identity() string {
	return getIdentity() + "/" + c.state.Object().AccountName
}

func (c *Client) Subscribe(kind MessageKind) <-chan *Message {
	c.mu.Lock()
	defer c.mu.Unlock()
	sub := subscriber{ch: make(chan *Message), once: false}
	c.watchers[kind][sub] = true
	return sub.ch
}

func (c *Client) SubscribeOnce(kind MessageKind) <-chan *Message {
	c.mu.Lock()
	defer c.mu.Unlock()
	sub := subscriber{ch: make(chan *Message), once: true}
	c.watchers[kind][sub] = true
	return sub.ch
}

func (c *Client) Send(receiver string, content interface{}) {
	kind := MessageKinds.Determine(content)
	if kind == UnknownMessageKind {
		panic("received a call to send an unknown type of message")
	}
	data, err := json.Marshal(content)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal message of kind: %v", err))
	}
	var message = &Message{
		Sender:   c.Identity(),
		Receiver: receiver,
		Content:  string(data),
		Kind:     kind,
	}
	serialized, err := json.Marshal(message)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal wrapped message of kind: %v", err))
	}
	c.mu.Lock()
	if _, err = c.conn.Write(append(serialized, '\n')); err != nil {
		c.logger.Log(0, logging.Error, fmt.Sprintf("[Client]: failed to write to relay: %v", err))
		defer c.Disconnect()
	}
	c.mu.Unlock()
}

func (c *Client) Broadcast(content interface{}) {
	c.Send(BroadcastReceiver, content)
}

func (c *Client) Connect(address string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		return nil
	}

	_ = c.state.SetPath("networking.connectingAddress", address)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		dialog.Message(fmt.Sprintf("Failed to connect to %s: %v", address, err))
		_ = c.state.SetPath("networking.connectingAddress", "")
		return err
	}

	c.conn = conn
	return nil
}

func (c *Client) Start() {
	go c.discoverRelays()

	for {
		if c.conn == nil {
			select {
			case <-time.After(1 * time.Second):
				continue
			case <-c.stop:
				return
			}
		}

		c.Send(RelayReceiver, RegistrationMessage{Identity: c.Identity()})
		go c.listenForMessages()
		select {
		case message := <-c.SubscribeOnce(AckRegistrationMessageKind):
			var ack AckRegistrationMessage
			if err := json.Unmarshal([]byte(message.Content), &ack); err != nil {
				c.logger.Log(0, logging.Error, fmt.Sprintf("[Client]: failed to unmarshal registration acknowledgement: %v", err))
				c.Disconnect()
				continue
			}
			if ack.Error != "" {
				dialog.Message(fmt.Sprintf("Failed to connect to relay: %s", ack.Error))
				c.Disconnect()
				continue
			}
			break
		case <-time.After(10 * time.Second):
			c.logger.Log(0, logging.Error, "[Client]: failed to connect to relay: no registration acknowledgement received!")
			c.Disconnect()
			continue
		case <-c.stop:
			return
		}
		fmt.Println("connected to relay")
		c.state.SetPath("networking.connectedAddress", c.conn.RemoteAddr().String())
		c.state.SetPath("networking.connectingAddress", "")
		<-c.disconnect
	}
}

func (c *Client) Close() {
	c.Disconnect()
	c.stop <- struct{}{}
}

func (c *Client) Disconnect() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn == nil {
		return
	}

	conn := c.conn
	c.conn = nil
	for _, watcher := range c.watchers {
		for sub := range watcher {
			sub.ch <- nil
			delete(watcher, sub)
		}
	}

	c.state.SetPath("networking.connectedAddress", "")
	if err := conn.Close(); err != nil {
		c.logger.Log(0, logging.Error, fmt.Sprintf("[Client]: error when disconnecting from relay: %v", err))
	}
}

func (c *Client) discoverRelays() {
	for {
		if *config.Concrete[bool](c.state, "networking.relayActive") {
			<-time.After(1 * time.Second)
		}

		resolver, err := zeroconf.NewResolver(nil)
		if err != nil {
			c.logger.Log(0, logging.Error, fmt.Sprintf("[Client]: failed to create resolver: %v", err))
			<-time.After(5 * time.Second)
			continue
		}

		entries := make(chan *zeroconf.ServiceEntry)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		go func(results <-chan *zeroconf.ServiceEntry) {
			var identities = make(map[string]*config.NetworkIdentity)
			for entry := range results {
				for _, txt := range entry.Text {
					if strings.HasPrefix(txt, "identity=") {
						identity, err := url.QueryUnescape(strings.TrimPrefix(txt, "identity="))
						if err != nil {
							continue
						}
						if identity == c.Identity() {
							continue
						}
						address := fmt.Sprintf("%s:%d", entry.AddrIPv4[0].String(), entry.Port)
						identities[address] = &config.NetworkIdentity{
							Identity: identity,
							Address:  fmt.Sprintf("%s:%d", entry.AddrIPv4[0].String(), entry.Port),
						}
					}
				}
			}
			if len(c.stop) > 0 {
				cancel()
				return
			}
			status := c.state.Object().Networking.Object()
			var removedIdentities []string
			var existingIdentities = make(map[string]bool)
			status.AvailableRelays.ForEach(func(id *config.NetworkIdentity) {
				if _, ok := identities[id.Address]; !ok {
					removedIdentities = append(removedIdentities, id.Address)
				} else {
					existingIdentities[id.Address] = true
				}
			})
			for _, id := range removedIdentities {
				_ = c.state.DeletePathf("networking.availableRelays[%s]", id)
			}
			for address, id := range identities {
				if _, ok := existingIdentities[address]; !ok {
					_ = c.state.AppendPathf("networking.availableRelays[%s]", address)
					_ = c.state.SetPathf(id.Identity, "networking.availableRelays[%s].identity", address)
				}
			}
			cancel()
		}(entries)

		err = resolver.Browse(ctx, "_revolution._tcp", "local.", entries)
		if err != nil {
			c.logger.Log(0, logging.Error, fmt.Sprintf("[Client]: failed to browse for relays: %v", err))
			<-time.After(5 * time.Second)
		}

		<-ctx.Done()

		select {
		case <-time.After(5 * time.Second):
			continue
		case <-c.stop:
			return
		}
	}
}

func (c *Client) handleConnectedIdentities(message Message) {
	var data ConnectedIdentitiesMessage
	if err := json.Unmarshal([]byte(message.Content), &data); err != nil {
		c.logger.Log(0, logging.Warning, fmt.Sprintf("[Client]: failed to unserialize connected identities message: %v", err))
		return
	}
	var identities = make(map[string]*config.NetworkIdentity)
	for _, identity := range data.Identities {
		identities[identity.Address] = &identity
	}
	var removedIdentities []string
	var existingIdentities = make(map[string]bool)
	status := c.state.Object().Networking.Object()
	status.ConnectedIdentities.ForEach(func(id *config.NetworkIdentity) {
		if _, ok := identities[id.Address]; !ok {
			removedIdentities = append(removedIdentities, id.Address)
		} else {
			existingIdentities[id.Address] = true
		}
	})
	for _, id := range removedIdentities {
		c.state.DeletePathf("networking.connectedIdentities[%s]", id)
	}
	for address, id := range identities {
		if _, ok := existingIdentities[address]; !ok {
			c.state.AppendPathf("networking.connectedIdentities[%s]", address)
			c.state.SetPathf(id.Identity, "networking.connectedIdentities[%s].identity", address)
		}
	}
}

func (c *Client) handleShutdown() {
	status := c.state.Object().Networking.Object()
	status.ConnectedIdentities.ForEach(func(id *config.NetworkIdentity) {
		c.state.DeletePathf("networking.connectedIdentities[%s]", id.Address)
	})
}

func (c *Client) listenForMessages() {
	scanner := bufio.NewScanner(c.conn)
	for scanner.Scan() {
		var msg Message
		if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {
			c.logger.Log(0, logging.Warning, fmt.Sprintf("[Client]: received invalid message from relay: %v", err))
			continue
		}
		switch msg.Kind {
		case ConnectedIdentitiesMessageKind:
			c.handleConnectedIdentities(msg)
			continue
		case ShutdownMessageKind:
			c.handleShutdown()
			break
		default:
			if _, ok := c.watchers[msg.Kind]; !ok {
				c.logger.Log(0, logging.Warning, "[Client]: received invalid message type from relay!")
				continue
			}
			for sub := range c.watchers[msg.Kind] {
				sub.ch <- &msg
				if sub.once {
					delete(c.watchers[msg.Kind], sub)
				}
			}
		}
	}
	if err := scanner.Err(); err != nil && err != io.EOF {
		c.logger.Log(0, logging.Error, fmt.Sprintf("Relay connection error: %v", err))
	}
	c.Disconnect()
	fmt.Println("disconnecting")
	c.disconnect <- struct{}{}
}
