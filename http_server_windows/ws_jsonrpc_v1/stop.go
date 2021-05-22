package ws_jsonrpc_v1

import (
	"errors"
	"os/exec"
	"strconv"
)

func (self *WsJsonrpcV1) Stop(request ConfigReq, response *map[string]string) (err error) {
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
		*response = map[string]string{
			"status": "ok",
		}
		return
	}

	// 写入进程kill名单
	self.KillCmdMap[request.Filename] = true
	// 进程没有退出
	if err = exec.Command("taskkill", "/F", "/T", "/PID", strconv.Itoa(oldCmd.Process.Pid)).Run(); err != nil {
		return
	}

	// 进程已经退出
	*response = map[string]string{
		"status": "ok",
	}

	return
}
