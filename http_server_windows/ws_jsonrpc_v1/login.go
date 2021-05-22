package ws_jsonrpc_v1

import (
	"errors"
)

type LoginRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func (self *WsJsonrpcV1) Login(request LoginRequest, response *map[string]string) (err error) {
	if request.Username != self.Username || request.Password != self.Password {
		err = errors.New("no authroity")
		return
	}
	*response = map[string]string{
		"session_id": self.SessionId,
	}
	return
}
