package common

import "github.com/nosyliam/revolution/pkg/config"

type MessageKind int

const (
	RegistrationMessageKind MessageKind = iota
	AckRegistrationMessageKind
	ConnectedIdentitiesMessageKind
	SetRoleMessageKind
	AckSetRoleMessageKind
	QueryMainAccountMessageKind
	MainAccountMessageKind
	VicDetectMessageKind
	NightDetectMessageKind
	UnknownMessageKind
)

type ClientRole int

const (
	MainClientRole     ClientRole = iota
	SearcherClientRole ClientRole = iota
	PassiveClientRole  ClientRole = iota
)

type Message struct {
	Kind     MessageKind
	Sender   string
	Receiver string
	Content  string

	Data string `json:"-"`
}

type Network struct {
	Client NetworkClient
	Relay  NetworkRelay
}

type NetworkClient interface {
	Start()
	Close()
	Subscribe(kind MessageKind) <-chan *Message
	SubscribeOnce(kind MessageKind) <-chan *Message
	Send(receiver string, content interface{})
	Broadcast(content interface{})
	Connect(relay config.NetworkIdentity)
	Disconnect() error
}

type NetworkRelay interface {
}
