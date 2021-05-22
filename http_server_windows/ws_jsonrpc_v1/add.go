package ws_jsonrpc_v1

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
)

func (self *WsJsonrpcV1) Add(request ConfigReq, response *map[string]string) (err error) {
	if err = self.Validate.Struct(request); err != nil {
		return
	}

	if request.SessionId != self.SessionId {
		err = errors.New("no authority")
		return
	}

	if strings.Contains("http_server", request.Filename) {
		err = errors.New("不要使用名称http_server")
		return
	}

	var tmpByte []byte
	switch request.Type {
	case "client":
		tmpByte, err = json.Marshal(request.Client)
	case "server":
		request.Server.CertFile = self.CertFile
		request.Server.KeyFile = self.KeyFile
		tmpByte, err = json.Marshal(request.Server)
	case "forward":
		tmpByte, err = json.Marshal(request.Forward)
	}
	if err != nil || len(tmpByte) == 0 {
		return
	}

	var file *os.File
	if file, err = os.Create(self.ConfigPath + "/" + request.Type + "/" + request.Filename + ".json"); err != nil {
		return
	}
	defer file.Close()
	if _, err = file.Write(tmpByte); err != nil {
		return
	}

	*response = map[string]string{
		"status": "ok",
	}
	return
}
