package http_server_windows

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os/exec"

	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/websocket"

	"github.com/n454149301/http_proxy/http_server_windows/ws_jsonrpc_v1"
)

type HttpServer struct {
	App      *iris.Application `json:"-"`
	Port     string            `json:"port,omitempty"`
	CertFile string            `json:"cert_file,omitempty"`
	KeyFile  string            `json:"key_file,omitempty"`

	Username   string `json:"username,omitempty"`
	Password   string `json:"password,omitempty"`
	ConfigPath string `json:"config_path,omitempty"`
	WebPath    string `json:"web_path,omitempty"`

	ClientPath  string `json:"client_path,omitempty"`
	ServerPath  string `json:"server_path,omitempty"`
	ForwardPath string `json:"forward_path,omitempty"`
}

type TmpWriterS struct {
	NsConn *websocket.NSConn
	App    *iris.Application
}

func (self *TmpWriterS) Write(p []byte) (n int, err error) {
	self.App.Logger().Debug("ws_jsonrpc write: ", string(p))
	n = len(p)
	if ok := self.NsConn.Conn.Write(websocket.Message{
		Body:     p,
		IsNative: true,
	}); !ok {
		err = errors.New("write iris websocket error")
	}

	return
}

func (self *HttpServer) Start() {
	self.App = iris.New()
	self.App.Use(recover.New())
	self.App.Use(logger.New(logger.Config{
		Status: true,
		IP:     true,
		Method: true,
		Path:   true,
		Query:  true,
	}))

	ws := websocket.New(websocket.DefaultGorillaUpgrader, websocket.Events{
		websocket.OnNativeMessage: func(nsConn *websocket.NSConn, msg websocket.Message) error {
			self.App.Logger().Debug("Server got: ", string(msg.Body), " from [", nsConn.Conn.ID(), "]")

			var conn io.ReadWriteCloser = struct {
				io.Writer
				io.ReadCloser
			}{
				Writer: &TmpWriterS{
					NsConn: nsConn,
					App:    self.App,
				},
				ReadCloser: ioutil.NopCloser(bytes.NewReader(msg.Body)),
			}
			rpc.ServeRequest(jsonrpc.NewServerCodec(conn))
			return nil
		},
	})

	ws.OnConnect = func(c *websocket.Conn) error {
		self.App.Logger().Debug("[%s] Connected to server!", c.ID())
		return nil
	}

	ws.OnDisconnect = func(c *websocket.Conn) {
		self.App.Logger().Debug("[%s] Disconnected from server", c.ID())
	}

	// APIv1注册
	tmpWsJsonrpcServe := (&ws_jsonrpc_v1.WsJsonrpcV1{
		App:      self.App,
		Validate: validator.New(),

		CertFile: self.CertFile,
		KeyFile:  self.KeyFile,

		Username:   self.Username,
		Password:   self.Password,
		ConfigPath: self.ConfigPath,

		ClientPath:  self.ClientPath,
		ServerPath:  self.ServerPath,
		ForwardPath: self.ForwardPath,

		CmdMap:     map[string]*exec.Cmd{},
		KillCmdMap: map[string]bool{},
	}).RegisterV1()
	defer func() {
		for _, tmpCmd := range tmpWsJsonrpcServe.CmdMap {
			if tmpCmd == nil || tmpCmd.Process == nil {
				continue
			}
			tmpCmd.Process.Release()
			tmpCmd.Process.Kill()
		}
	}()
	self.App.Get("/http_proxy/api/ws-jsonrpc/v1", websocket.Handler(ws))

	self.App.HandleDir("/", self.WebPath, iris.DirOptions{
		// IndexName: "index.html",
		Gzip:     true,
		ShowList: false,
	})

	self.App.Run(iris.TLS(":"+self.Port, self.CertFile, self.KeyFile))
}
