package ws_jsonrpc_v1

import (
	"github.com/n454149301/http_proxy/client"
	"github.com/n454149301/http_proxy/forward"
	"github.com/n454149301/http_proxy/server"
)

type ConfigReq struct {
	Type      string `json:"type,omitempty" validate:"oneof=client server forward"`
	Filename  string `json:"filename,omitempty"`
	SessionId string `json:"session_id,omitempty"`

	Client  client.Client   `json:"client,omitempty"`
	Server  server.Server   `json:"server,omitempty"`
	Forward forward.Forward `json:"forward,omitempty"`
}
