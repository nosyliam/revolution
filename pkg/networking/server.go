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
	"strings"
	"sync"
)

const BroadcastReceiver = "!BROADCAST"
const RelayReceiver = "!RELAY"

type Relay struct {
	mu         sync.Mutex
	client     *Client
	server     *zeroconf.Server
	port       int
	identity   string
	identities map[string]net.Conn
	roles      map[string]string
	state      *config.Object[config.MacroState]
	logger     *logging.Logger
	stop       chan struct{}
}

func NewRelay(client *Client, state *config.Object[config.MacroState], logger *logging.Logger) *Relay {
	return &Relay{
		identity:   state.Object().AccountName,
		port:       45645, // squid game?
		identities: make(map[string]net.Conn),
		state:      state,
		logger:     logger,
	}
}

func (r *Relay) Identity() string {
	return getIdentity()
}

func (r *Relay) Start() error {
	defer r.state.SetPath("networking.relayStarting", false)
	r.state.SetPath("networking.relayStarting", true)
	listener, err := net.Listen("tcp", fmt.Sprintf(":45645"))
	if err != nil {
		return err
	}
	defer listener.Close()

	if err := r.client.Connect(listener.Addr().String()); err != nil {
		return errors.Wrap(err, "failed to connect to local relay")
	}

	txtRecords := []string{fmt.Sprintf("identity=%s", r.identity)}
	r.server, err = zeroconf.Register("RelayService", "_relay._tcp", "local.", r.port, txtRecords, nil)
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
				conn, err := listener.Accept()
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
	if err = json.Unmarshal([]byte(message.Content), &message); err != nil {
		r.logger.Log(0, logging.Warning, fmt.Sprintf("[Relay]: Failed to decode registration message from client %s", conn.RemoteAddr().String()))
		conn.Close()
		return
	}

	identity := registrationMessage.Identity
	if _, ok := r.identities[identity]; ok {
		data, _ := json.Marshal(AckRegistrationMessage{Error: fmt.Sprintf("The identity %s is already connected to this relay!", identity)})
		conn.Write(data)
		conn.Close()
		return
	}

	r.mu.Lock()
	r.identities[identity] = conn
	r.mu.Unlock()

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
		message.Data = line
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
	switch message.Receiver {
	case RelayReceiver:
		if message.Kind == SetRoleMessageKind {

		}
	case BroadcastReceiver:
		for id, conn := range r.identities {
			if id != message.Sender {
				conn.Write([]byte(message.Data))
			}
		}
		return
	default:
		if strings.HasPrefix(message.Receiver, "!") {
			role := strings.TrimPrefix(message.Receiver, "!")
			for identity, idRole := range r.roles {
				if idRole == role && identity != message.Sender {
					r.identities[identity].Write([]byte(message.Data))
				}
			}
		} else {
			if conn, ok := r.identities[message.Receiver]; !ok {
				r.logger.Log(0, logging.Warning,
					fmt.Sprintf("[Relay]: Failed to forward message from %s->%s: invalid receiver", message.Sender, message.Receiver))
				return
			} else {
				conn.Write([]byte(message.Data))
			}
		}

	}

}
