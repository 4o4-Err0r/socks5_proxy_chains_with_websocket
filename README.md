# socks5_proxy_chains_with_websocket
基于websocket连接的socks5代理链路工具。


使用方式：

1、生成server.crt和server.key：

make openssl_key

2、在linux环境下编译指定目标环境的可执行程序

3、配置server.json之后在服务器上执行


# server.json
```json
{
	"cert_file": "./bin/server.crt",
	"key_file": "./bin/server.key",
	"port": "12345",
	"uuid": "3F2504E0-4F89-11D3-9A0C-0305E82C3301"
}
```

在服务器上执行：

./server ./server.json

uuid是自定义的通信密钥，server和client配置一致的情况才能通信。

4、配置client.json之后在本地执行，可以放在路由器上执行。

假设服务器工作在192.168.11.1上，监听的端口是12345。

# client.json
```json
{
	"out_addr": "192.168.11.1:12345",
	"port": "55555",
	"uuid": "3F2504E0-4F89-11D3-9A0C-0305E82C3301"
}
```

在本机或者局域网某台设备，或者路由器上执行。

./client ./client.json

client也会监听一个端口，55555。假设client的ip是192.168.11.5

局域网内的socks5代理配置成，主机名：192.168.11.5，端口：55555,就可以使用代理上网了。

uuid是自定义的通信密钥，server和client配置一致的情况才能通信。

5、多级跳中间节点，在中间跳服务器上执行：

./forward ./forward.json

client的目的地址可以跳到forward监听的服务器端口上，数据会跳到forward.json中out_addr指向的下一跳。

# forward.json
```json
{
	"out_addr": "192.168.1.1:12345",
	"port": "12345"
}
```

6、代理集群链路配置web服务器启动。

./http_server ./http_server.json

# http_server.json
```json
{
	"cert_file": "./bin/server.crt",
	"key_file": "./bin/server.key",
	"port": "22233",

	"username": "admin",
	"password": "123456",
	"config_path": "./config",
	"web_path": "./www",

	"ps_ex_parm": "aux",

	"client_path": "./bin/client",
	"server_path": "./bin/server",
	"forward_path": "./bin/forward"
}
```

```shell
cert_file: make openssl_key生成的crt证书路径
key_file: make openssl_key生成的key证书路径
port: web服务器监听的端口

username: web登陆时候的帐号
password: web登陆时候的密码
config_path: 各种代理配置保存的路径
web_path: web静态资源所在文件夹
ps_ex_parm: 当前系统查看进程的ps附加配置参数

client_path: client程序所在路径
server_path: server程序所在路径
forward_path: forward程序所在路径
```
