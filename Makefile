GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BUILD_PATH=build

clean:
	rm -rf build

buildAMD64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build \
		-ldflags="-X git.rickiekarp.net/rickie/fileguardian/config.Version=$(shell git rev-parse HEAD) \
		-X git.rickiekarp.net/rickie/fileguardian/config.ApiProtocol=https \
		-X git.rickiekarp.net/rickie/fileguardian/config.ApiHost=api.rickiekarp.net" \
		-o $(BUILD_PATH)/fileguardian \
		main.go

buildARM64v7:
	CGO_ENABLED=1 GOOS=linux GOARCH=arm64 GOARM=7 CC=aarch64-linux-gnu-gcc \
	go build \
		-ldflags="-X git.rickiekarp.net/rickie/fileguardian/config.Version=$(shell git rev-parse HEAD) \
		-X git.rickiekarp.net/rickie/fileguardian/config.ApiProtocol=https \
		-X git.rickiekarp.net/rickie/fileguardian/config.ApiHost=api.rickiekarp.net" \
		-o $(BUILD_PATH)/fileguardian \
		main.go

deploy:
	rsync -rlvpt --delete build/fileguardian pi:~/tools/
