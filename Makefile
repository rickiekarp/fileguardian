GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BUILD_PATH=build

clean:
	rm -rf build

buildARM64v7:
	CGO_ENABLED=1 GOOS=linux GOARCH=arm64 GOARM=7 CC=aarch64-linux-gnu-gcc \
	go build \
		-ldflags="-X git.rickiekarp.net/rickie/fileguardian/storage.Version=$(shell git rev-parse HEAD)" \
		-o $(BUILD_PATH)/fileguardian \
		main.go

deploy:
	rsync -rlvpt --delete build/fileguardian pi:~/tools/
