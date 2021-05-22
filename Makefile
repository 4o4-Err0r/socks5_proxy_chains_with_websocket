export GOPROXY=https://goproxy.cn

all:
	go build -o ./bin/client ./cmd/client.go
	go build -o ./bin/server ./cmd/server.go
	go build -o ./bin/forward ./cmd/forward.go
	go build -o ./bin/http_server ./cmd/http_server.go

windows:
	CGO_ENABLED=0 GOOS=windows go build -o ./bin/client.exe ./cmd/client.go
	CGO_ENABLED=0 GOOS=windows go build -o ./bin/server.exe ./cmd/server.go
	CGO_ENABLED=0 GOOS=windows go build -o ./bin/forward.exe ./cmd/forward.go
	CGO_ENABLED=0 GOOS=windows go build -o ./bin/http_server_windows.exe ./cmd/http_server_windows.go

mac:
	CGO_ENABLED=0 GOOS=darwin go build -o ./bin/client ./cmd/client.go
	CGO_ENABLED=0 GOOS=darwin go build -o ./bin/server ./cmd/server.go
	CGO_ENABLED=0 GOOS=darwin go build -o ./bin/forward ./cmd/forward.go
	CGO_ENABLED=0 GOOS=darwin go build -o ./bin/http_server ./cmd/http_server.go

linux_mips:
	CGO_ENABLED=0 GOOS=linux GOARCH=mips go build -o ./bin/client ./cmd/client.go
	CGO_ENABLED=0 GOOS=linux GOARCH=mips go build -o ./bin/server ./cmd/server.go
	CGO_ENABLED=0 GOOS=linux GOARCH=mips go build -o ./bin/forward ./cmd/forward.go
	CGO_ENABLED=0 GOOS=linux GOARCH=mips go build -o ./bin/http_server ./cmd/http_server.go

linux_armv5:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=5 go build -o ./bin/client ./cmd/client.go
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=5 go build -o ./bin/server ./cmd/server.go
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=5 go build -o ./bin/forward ./cmd/forward.go
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=5 go build -o ./bin/http_server ./cmd/http_server.go

linux_armv7:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -o ./bin/client ./cmd/client.go
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -o ./bin/server ./cmd/server.go
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -o ./bin/forward ./cmd/forward.go
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -o ./bin/http_server ./cmd/http_server.go

linux_armv8:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./bin/client ./cmd/client.go
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./bin/server ./cmd/server.go
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./bin/forward ./cmd/forward.go
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./bin/http_server ./cmd/http_server.go

openssl_key:
	openssl genrsa -out ./bin/server.key 2048
	openssl req -nodes -new -key ./bin/server.key -subj "/CN=localhost" -out ./bin/server.csr
	openssl x509 -req -sha256 -days 365 -in ./bin/server.csr -signkey ./bin/server.key -out ./bin/server.crt

releases:
	make clean; make openssl_key; make all; tar cvfz ./socks5_proxy_linux.tar.gz bin res; mv ./socks5_proxy_linux.tar.gz ../releases/socks5_proxy
	make clean; make openssl_key; make linux_armv5; tar cvfz ./socks5_proxy_linux_armv5.tar.gz bin res; mv ./socks5_proxy_linux_armv5.tar.gz ../releases/socks5_proxy
	make clean; make openssl_key; make linux_armv7; tar cvfz ./socks5_proxy_linux_armv7.tar.gz bin res; mv ./socks5_proxy_linux_armv7.tar.gz ../releases/socks5_proxy
	make clean; make openssl_key; make linux_armv8; tar cvfz ./socks5_proxy_linux_armv8.tar.gz bin res; mv ./socks5_proxy_linux_armv8.tar.gz ../releases/socks5_proxy
	make clean; make openssl_key; make linux_mips; tar cvfz ./socks5_proxy_linux_mips.tar.gz bin res; mv ./socks5_proxy_linux_mips.tar.gz ../releases/socks5_proxy
	make clean; make openssl_key; make mac; tar cvfz ./socks5_proxy_macos.tar.gz bin res; mv ./socks5_proxy_macos.tar.gz ../releases/socks5_proxy
	make clean; make openssl_key; make windows; tar cvfz ./socks5_proxy_windows.tar.gz bin res; mv ./socks5_proxy_windows.tar.gz ../releases/socks5_proxy

clean:
	rm -rfd ./bin/*
