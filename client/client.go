package client

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"

	"golang.org/x/net/websocket"
)

type Client struct {
	OutAddr string `json:"out_addr,omitempty"`
	TlsIp   string `json:"tls_ip,omitempty"`
	TlsHost string `json:"tls_host,omitempty"`

	Port string `json:"port,omitempty"`
	UUID string `json:"uuid,omitempty"`
}

func (self *Client) Start() {
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

func (self *Client) HandleClientRequest(client net.Conn) {
	if client == nil {
		return
	}
	var config *websocket.Config
	var err error
	if config, err = websocket.NewConfig("wss://"+self.OutAddr+"/http_proxy", "https://"+self.OutAddr); err != nil {
		fmt.Println("ws new err:", err)
		return
	}
	config.TlsConfig = &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         self.TlsHost,
	}
	config.Header.Add("UUID", self.UUID)
	var server *websocket.Conn

	if self.TlsIp == "" {
		if server, err = websocket.DialConfig(config); err != nil {
			fmt.Println("dial err:", err)
			return
		}
	} else {
		var serverTls net.Conn
		if config.Location == nil {
			fmt.Println("config.Location:nil")
			return
		}
		if config.Origin == nil {
			fmt.Println("config.Origin:nil")
			return
		}
		dialer := config.Dialer
		if dialer == nil {
			dialer = &net.Dialer{}
		}
		switch config.Location.Scheme {
		case "wss":
			if serverTls, err = tls.DialWithDialer(dialer, "tcp", self.TlsIp, config.TlsConfig); err != nil {
				fmt.Println("tcp dial err:", err)
				return
			}
		default:
			fmt.Println("scheme err:", config.Location.Scheme)
			return
		}
		if server, err = websocket.NewClient(config, serverTls); err != nil {
			server.Close()
			fmt.Println("ws client dial err:", err)
			return
		}
	}

	go self.CopyIO(client, server)
	go self.CopyIO(server, client)
}

func (self *Client) CopyIO(src, dest io.ReadWriteCloser) {
	defer src.Close()
	defer dest.Close()
	io.Copy(src, dest)
}
