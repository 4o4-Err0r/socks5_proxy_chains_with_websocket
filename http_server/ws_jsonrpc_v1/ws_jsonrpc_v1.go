package ws_jsonrpc_v1

import (
	"crypto/sha512"
	"encoding/hex"
	"net/rpc"
	"os/exec"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
)

type WsJsonrpcV1 struct {
	App      *iris.Application
	Validate *validator.Validate

	CertFile string
	KeyFile  string

	Username   string
	Password   string
	ConfigPath string

	PsExParm string

	ClientPath  string
	ServerPath  string
	ForwardPath string

	// 自动生成session_id
	SessionId string

	// 启动的cmd map
	CmdMap map[string]*exec.Cmd
	// 准备杀死的cmd map
	KillCmdMap map[string]bool
}

func (self *WsJsonrpcV1) RegisterV1() *WsJsonrpcV1 {
	rpc.RegisterName("Api", self)
	go self.RefreshSessionId()
	return self
}

func (self *WsJsonrpcV1) RefreshSessionId() {
	for {
		hashStr := sha512.Sum512([]byte(time.Now().String()))
		self.SessionId = hex.EncodeToString(hashStr[:])

		time.Sleep(20 * time.Minute)
	}
}
