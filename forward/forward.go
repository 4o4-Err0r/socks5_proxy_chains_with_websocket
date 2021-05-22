package forward

import (
	"fmt"
	"io"
	"net"
)

type Forward struct {
	OutAddr string `json:"out_addr,omitempty"`

	Port string `json:"port,omitempty"`
}

func (self *Forward) Start() {
	l, err := net.Listen("tcp", ":"+self.Port)
	if err != nil {
		panic(err)
	}
	for {
		client, err := l.Accept()
		if err != nil {
			panic(err)
		}
		go self.HandleClientRequest(client)
	}
}

func (self *Forward) HandleClientRequest(client net.Conn) {
	if client == nil {
		return
	}
	server, err := net.Dial("tcp", self.OutAddr)
	if err != nil {
		fmt.Println("dial err:", err)
		return
	}

	go self.CopyIO(client, server)
	go self.CopyIO(server, client)
}

func (self *Forward) CopyIO(src, dest io.ReadWriteCloser) {
	defer src.Close()
	defer dest.Close()
	io.Copy(src, dest)
}
