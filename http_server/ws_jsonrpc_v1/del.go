package ws_jsonrpc_v1

import (
	"errors"
	"os"
	"syscall"
	"time"
)

func (self *WsJsonrpcV1) Del(request ConfigReq, response *map[string]string) (err error) {
	if err = self.Validate.Struct(request); err != nil {
		return
	}

	if request.SessionId != self.SessionId {
		err = errors.New("no authority")
		return
	}

	i := 0
	for i = 0; i < 3; i++ {
		if oldCmd, ok := self.CmdMap[request.Filename]; ok {
			// 写入进程kill名单
			self.KillCmdMap[request.Filename] = true
			// 进程没有退出
			syscall.Kill(oldCmd.Process.Pid, syscall.SIGKILL)
			time.Sleep(time.Second)
		}
	}
	if i == 2 {
		err = errors.New("进程无法停止")
		return
	}

	if err = os.Remove(self.ConfigPath + "/" + request.Type + "/" + request.Filename + ".json"); err != nil {
		return
	}

	*response = map[string]string{
		"status": "ok",
	}
	return
}
