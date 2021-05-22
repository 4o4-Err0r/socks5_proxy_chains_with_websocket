package ws_jsonrpc_v1

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func (self *WsJsonrpcV1) Ls(request ConfigReq, response *[]string) (err error) {
	if err = self.Validate.Struct(request); err != nil {
		return
	}

	if request.SessionId != self.SessionId {
		err = errors.New("no authority")
		return
	}

	*response = []string{}
	filepath.Walk(self.ConfigPath+"/"+request.Type, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(info.Name(), ".json") {
			*response = append(*response, info.Name())
		}
		return nil
	})
	return
}
