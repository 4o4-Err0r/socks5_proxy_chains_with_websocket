package ws_jsonrpc_v1

import (
	"errors"
	"io/ioutil"
	"os"
)

func (self *WsJsonrpcV1) Cat(request ConfigReq, response *string) (err error) {
	if err = self.Validate.Struct(request); err != nil {
		return
	}

	if request.SessionId != self.SessionId {
		err = errors.New("no authority")
		return
	}

	var file *os.File
	if file, err = os.Open(self.ConfigPath + "/" + request.Type + "/" + request.Filename + ".json"); err != nil {
		return
	}
	var fileData []byte
	if fileData, err = ioutil.ReadAll(file); err != nil {
		return
	}

	*response = string(fileData)
	return
}
