package networking

import (
	"encoding/json"
	"fmt"
	. "github.com/nosyliam/revolution/pkg/common"
	"github.com/nosyliam/revolution/pkg/config"
)

type MessageKindEnumerator []MessageKind

var MessageKinds = MessageKindEnumerator{
	RegistrationMessageKind,
	AckRegistrationMessageKind,
	AckSetRoleMessageKind,
	SetRoleMessageKind,
	VicDetectMessageKind,
	NightDetectMessageKind,
}

func (m *MessageKindEnumerator) Determine(data interface{}) MessageKind {
	switch data.(type) {
	case RegistrationMessage:
		return RegistrationMessageKind
	case AckRegistrationMessage:
		return AckRegistrationMessageKind
	case ConnectedIdentitiesMessage:
		return ConnectedIdentitiesMessageKind
	case SetRoleMessage:
		return SetRoleMessageKind
	case AckSetRoleMessage:
		return AckSetRoleMessageKind
	case VicDetectMessage:
		return VicDetectMessageKind
	case NightDetectMessage:
		return NightDetectMessageKind
	}
	return UnknownMessageKind
}

type RegistrationMessage struct {
	Identity string
}

type AckRegistrationMessage struct {
	Error string
}

type ConnectedIdentitiesMessage struct {
	Identities []config.NetworkIdentity
}

type SetRoleMessage struct {
	Role string
}

type AckSetRoleMessage struct {
	Role string
}

type VicDetectMessage struct {
	AccessCode string
	Field      string
	TileX      int
	TileY      int
}

type NightDetectMessage struct {
	AccessCode string
}

type EmptyMessage struct{}

func SubscribeMessage[T any](macro *Macro, kind MessageKind, callback func(message *T)) {
	watcher := macro.Network.Client.Subscribe(kind)
	for {
		message := <-watcher
		if message == nil {
			return
		}
		var msg T
		if err := json.Unmarshal([]byte(message.Content), &msg); err != nil {
			macro.Error <- fmt.Sprintf("Unable to decode message from %s: %v", message.Sender, err)
			continue
		}
		callback(&msg)
	}
}

func UnsubscribeMessage[T any](macro *Macro, kind MessageKind) {

}
