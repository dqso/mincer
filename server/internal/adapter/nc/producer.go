package nc_adapter

import (
	"github.com/dqso/mincer/server/internal/api"
	"github.com/wirepair/netcode"
	"google.golang.org/protobuf/proto"
)

type Producer struct {
	server *netcode.Server
}

func NewProducer(server *netcode.Server) *Producer {
	return &Producer{
		server: server,
	}
}

func (p *Producer) marshalMessage(code api.Code, payload proto.Message) ([]byte, error) {
	msg := api.Message{Code: code}
	var err error
	if msg.Payload, err = proto.Marshal(payload); err != nil {
		return nil, err
	}
	return proto.Marshal(&msg)
}
