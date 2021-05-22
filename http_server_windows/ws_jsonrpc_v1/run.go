package ws_jsonrpc_v1

import (
	"errors"
	"fmt"
	"os/exec"
)

func (self *WsJsonrpcV1) Run(request ConfigReq, response *map[string]string) (err error) {
	if err = self.Validate.Struct(request); err != nil {
		return
	}

	if request.SessionId != self.SessionId {
		err = errors.New("no authority")
		return
	}

	if oldCmd, ok := self.CmdMap[request.Filename]; ok {
		if oldCmd.ProcessState == nil {
			*response = map[string]string{
				"status": "ok",
			}
			return
		}
	}

	var cmdPath string
	var cfgPath string
	switch request.Type {
	case "client":
		cmdPath = self.ClientPath
		cfgPath = self.ConfigPath + "/client/" + request.Filename + ".json"
	case "server":
		cmdPath = self.ServerPath
		cfgPath = self.ConfigPath + "/server/" + request.Filename + ".json"
	case "forward":
		cmdPath = self.ForwardPath
		cfgPath = self.ConfigPath + "/forward/" + request.Filename + ".json"
	}

	fmt.Println(cmdPath, cfgPath)
	cmd := exec.Command(cmdPath, cfgPath)
	if err = cmd.Start(); err != nil {
		return
	}

	defer func() {
		filenameKill := request.Filename
		cmdPathKill := cmdPath
		cfgPathKill := cfgPath
		go func() {
			cmdKill := cmd

			for {
				cmdKill.Wait()
				fmt.Println(filenameKill, ":is dead")
				if _, ok := self.KillCmdMap[filenameKill]; ok {
					// 如果在kill名单里，让进程退出
					fmt.Println(filenameKill, ":is kill")
					delete(self.CmdMap, request.Filename)
					delete(self.KillCmdMap, request.Filename)
					return
				}

				// 如果不在kill名单里，让进程重启
				fmt.Println(filenameKill, ":is restart")
				cmdKill = exec.Command(cmdPathKill, cfgPathKill)
				if err = cmdKill.Start(); err != nil {
					// 如果重启失败，进程退出
					fmt.Println(filenameKill, ":restart error:", err)
					delete(self.CmdMap, request.Filename)
					delete(self.KillCmdMap, request.Filename)
					return
				}
				self.CmdMap[filenameKill] = cmdKill
			}
		}()
	}()
	self.CmdMap[request.Filename] = cmd

	*response = map[string]string{
		"status": "ok",
	}
	return
}
