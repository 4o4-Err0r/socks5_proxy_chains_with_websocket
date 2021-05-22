package ws_jsonrpc_v1

import (
	"errors"
	"fmt"
	"os/exec"
)

func (self *WsJsonrpcV1) Check(request ConfigReq, response *string) (err error) {
	if err = self.Validate.Struct(request); err != nil {
		return
	}

	if request.SessionId != self.SessionId {
		err = errors.New("no authority")
		return
	}

	var oldCmd *exec.Cmd
	var ok bool
	if oldCmd, ok = self.CmdMap[request.Filename]; !ok {
		err = errors.New("未启动")
		return
	}

	fmt.Println(oldCmd.Process.Pid, oldCmd.Process)
	*response = oldCmd.String()

	return
}
