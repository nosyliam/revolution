package networking

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/grandcat/zeroconf"
	. "github.com/nosyliam/revolution/pkg/common"
	"github.com/nosyliam/revolution/pkg/config"
	"github.com/nosyliam/revolution/pkg/logging"
	"github.com/pkg/errors"
	"net"
	"net/url"
	"strings"
	"sync"
)

const BroadcastReceiver = "!BROADCAST"
const RelayReceiver = "!RELAY"

type Relay struct {
	mu         sync.Mutex
	client     *Client
	server     *zeroconf.Server
	listener   net.Listener
	port       int
	identities map[string]net.Conn
	roles      map[string]string
	state      *config.Object[config.MacroState]
	logger     *logging.Logger
	stop       chan struct{}
}

func NewRelay(client *Client, state *config.Object[config.MacroState], logger *logging.Logger) *Relay {
	return &Relay{
		client:     client,
		port:       45645, // squid game?
		identities: make(map[string]net.Conn),
		state:      state,
		logger:     logger,
	}
}

func (r *Relay) Identity() string {
	return getIdentity() + "/" + r.state.Object().AccountName
}

func (r *Relay) Start() error {
	var err error
	defer r.state.SetPath("networking.relayStarting", false)
	r.state.SetPath("networking.relayStarting", true)
	r.listener, err = net.Listen("tcp", fmt.Sprintf(":45645"))
	if err != nil {
		return err
	}

	if err := r.client.Connect(r.listener.Addr().String()); err != nil {
		return errors.Wrap(err, "failed to connect to local relay")
	}

	txtRecords := []string{fmt.Sprintf("identity=%s", url.QueryEscape(r.Identity()))}
	r.server, err = zeroconf.Register("RevolutionMacro", "_revolution._tcp", "local.", r.port, txtRecords, nil)
	if err != nil {
		return errors.Wrap(err, "failed to start zeroconf")
	}

	r.state.SetPath("networking.relayActive", true)

	go func() {
		for {
			select {
			case <-r.stop:
				return
			default:
				conn, err := r.listener.Accept()
				if err != nil {
					if errors.Is(err, net.ErrClosed) || err.Error() == "use of closed network connection" {
						return
					}
					r.logger.Log(0, logging.Error, fmt.Sprintf("[Relay]: failed to accept new connection: %v", err))
					continue
				}
				go r.handleConnection(conn)
			}

		}
	}()

	return nil
}

func (r *Relay) Stop() {
	var message = Message{
		Kind:     ShutdownMessageKind,
		Receiver: BroadcastReceiver,
		Sender:   RelayReceiver,
		Content:  "{}",
	}
	r.handleMessage(&message)
	r.state.SetPath("networking.relayActive", false)
	r.listener.Close()
	r.server.Shutdown()
	for _, conn := range r.identities {
		conn.Close()
	}
	clear(r.identities)
	clear(r.roles)
}

func (r *Relay) handleConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)
	data, err := reader.ReadString('\n')
	if err != nil {
		r.logger.Log(0, logging.Warning, fmt.Sprintf("[Relay]: Failed to receive registration message from client %s", conn.RemoteAddr().String()))
		conn.Close()
		return
	}

	var message Message
	if err := json.Unmarshal([]byte(data), &message); err != nil || message.Kind != RegistrationMessageKind || message.Receiver != RelayReceiver {
		r.logger.Log(0, logging.Warning, fmt.Sprintf("[Relay]: Failed to decode message from client %s", conn.RemoteAddr().String()))
		conn.Close()
		return
	}

	var registrationMessage RegistrationMessage
	if err = json.Unmarshal([]byte(message.Content), &registrationMessage); err != nil {
		r.logger.Log(0, logging.Warning, fmt.Sprintf("[Relay]: Failed to decode registration message from client %s", conn.RemoteAddr().String()))
		conn.Close()
		return
	}

	identity := registrationMessage.Identity
	ack := Message{
		Kind:     AckRegistrationMessageKind,
		Receiver: identity,
		Sender:   RelayReceiver,
		Content:  "{}",
	}

	if _, ok := r.identities[identity]; ok {
		byteData, _ := json.Marshal(AckRegistrationMessage{Error: fmt.Sprintf("The identity \"%s\" is already connected to this relay!", identity)})
		ack.Content = string(byteData)
		msg, _ := json.Marshal(ack)
		conn.Write(append(msg, "\r\n"...))
		conn.Close()
		return
	}

	r.mu.Lock()
	r.identities[identity] = conn
	r.mu.Unlock()

	r.handleMessage(&ack)
	r.broadcastIdentities()

	defer func() {
		r.mu.Lock()
		delete(r.identities, identity)
		delete(r.roles, identity)
		r.mu.Unlock()
		conn.Close()
	}()

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		var message Message
		if err = json.Unmarshal([]byte(line), &message); err != nil {
			r.logger.Log(0, logging.Warning, fmt.Sprintf("[Relay]: Failed to decode message from client %s: %v", identity, err))
			continue
		}
		message.Data = []byte(line)
		r.mu.Lock()
		r.handleMessage(&message)
		r.mu.Unlock()
	}
}

func (r *Relay) broadcastIdentities() {
	r.mu.Lock()
	defer r.mu.Unlock()
	var identities ConnectedIdentitiesMessage
	for identity, conn := range r.identities {
		role, _ := r.roles[identity]
		identities.Identities = append(identities.Identities, config.NetworkIdentity{
			Address:  conn.RemoteAddr().String(),
			Identity: identity,
			Role:     role,
		})
	}
	data, _ := json.Marshal(identities)
	var message = Message{
		Kind:     ConnectedIdentitiesMessageKind,
		Receiver: BroadcastReceiver,
		Sender:   RelayReceiver,
		Content:  string(data),
	}
	r.handleMessage(&message)
}

func (r *Relay) handleMessage(message *Message) {
	if string(message.Data) == "" {
		message.Data, _ = json.Marshal(message)
	}
	switch message.Receiver {
	case RelayReceiver:
		if message.Kind == SetRoleMessageKind {

		}
	case BroadcastReceiver:
		for id, conn := range r.identities {
			if id != message.Sender {
				conn.Write(append(message.Data, "\r\n"...))
			}
		}
		return
	default:
		if strings.HasPrefix(message.Receiver, "!") {
			role := strings.TrimPrefix(message.Receiver, "!")
			for identity, idRole := range r.roles {
				if idRole == role && identity != message.Sender {
					r.identities[identity].Write(append(message.Data, "\r\n"...))
				}
			}
		} else {
			if conn, ok := r.identities[message.Receiver]; !ok {
				r.logger.Log(0, logging.Warning,
					fmt.Sprintf("[Relay]: Failed to forward message from %s->%s: invalid receiver", message.Sender, message.Receiver))
				return
			} else {
				fmt.Println(conn.Write(append(message.Data, "\r\n"...)))
			}
		}

	}

}
