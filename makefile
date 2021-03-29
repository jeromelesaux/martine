CC=go
RM=rm
MV=mv


SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

VERSION:=$(shell grep -m1 "appVersion" *.go | sed 's/[", ]//g' | cut -d= -f2)
suffix=$(shell grep -m1 "version" *.go | sed 's/[", ]//g' | cut -d= -f2 | sed 's/[0-9.]//g')
snapshot=$(shell date +%FT%T)
BINARY=binaries

ifeq ($(suffix),rc)
	appversion=$(VERSION)$(snapshot)
else 
	appversion=$(VERSION)
endif 

.DEFAULT_GOAL:=build

clean:
	rm -fr ${BINARY}
	mkdir ${BINARY}

build: 
	rm -fr ${BINARY}
	mkdir ${BINARY}
	@echo "Compilation for linux"
	mkdir ${BINARY}/martine-linux-64bits
	GOOS=linux go build ${LDFLAGS} -o ${BINARY}/martine-linux-64bits/martine $(SOURCEDIR)/main.go $(SOURCEDIR)/process.go
	GOOS=linux go build ${LDFLAGS} -o ${BINARY}/martine-linux-64bits/prepare_delta $(SOURCEDIR)/resources/formatter/delta/prepare_delta.go
	GOOS=linux go build ${LDFLAGS} -o ${BINARY}/martine-linux-64bits/format_sprite $(SOURCEDIR)/resources/formatter/sprites/format_sprite.go
	GOOS=linux go build ${LDFLAGS} -o ${BINARY}/martine-linux-64bits/format_data $(SOURCEDIR)/resources/formatter/data/format_data.go
	zip ${BINARY}/martine-$(appversion)-linux.zip ${BINARY}/martine-linux-64bits/* ./resources/*
	@echo "Compilation for windows"
	mkdir ${BINARY}/martine-windows-64bits
	GOOS=windows go build ${LDFLAGS} -o ${BINARY}/martine-windows-64bits/martine.exe $(SOURCEDIR)/main.go $(SOURCEDIR)/process.go
	GOOS=windows go build ${LDFLAGS} -o ${BINARY}/martine-windows-64bits/prepare_delta.exe $(SOURCEDIR)/resources/formatter/delta/prepare_delta.go
	GOOS=windows go build ${LDFLAGS} -o ${BINARY}/martine-windows-64bits/format_sprite.exe $(SOURCEDIR)/resources/formatter/sprites/format_sprite.go
	GOOS=windows go build ${LDFLAGS} -o ${BINARY}/martine-windows-64bits/format_data.exe $(SOURCEDIR)/resources/formatter/data/format_data.go
	zip ${BINARY}/martine-$(appversion)-windows.zip  ${BINARY}/martine-windows-64bits/* ./resources/*
	@echo "Compilation for macos"
	mkdir ${BINARY}/martine-darwin-64bits 
	GOOS=darwin go build ${LDFLAGS} -o ${BINARY}/martine-darwin-64bits/martine $(SOURCEDIR)/main.go $(SOURCEDIR)/process.go
	GOOS=darwin go build ${LDFLAGS} -o ${BINARY}/martine-darwin-64bits/prepare_delta $(SOURCEDIR)/resources/formatter/delta/prepare_delta.go
	GOOS=darwin go build ${LDFLAGS} -o ${BINARY}/martine-darwin-64bits/format_sprite $(SOURCEDIR)/resources/formatter/sprites/format_sprite.go
	GOOS=darwin go build ${LDFLAGS} -o ${BINARY}/martine-darwin-64bits/format_data $(SOURCEDIR)/resources/formatter/data/format_data.go
	zip ${BINARY}/martine-$(appversion)-macos.zip ${BINARY}/martine-darwin-64bits/* ./resources/*
	@echo "Compilation for raspberry pi Raspbian 64 bits"
	mkdir ${BINARY}/martine-linux-arm64bits
	GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o ${BINARY}/martine-linux-arm64bits/martine $(SOURCEDIR)/main.go $(SOURCEDIR)/process.go
	GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o ${BINARY}/martine-linux-arm64bits/prepare_delta $(SOURCEDIR)/resources/formatter/delta/prepare_delta.go
	GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o ${BINARY}/martine-linux-arm64bits/format_sprite $(SOURCEDIR)/resources/formatter/sprites/format_sprite.go
	GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o ${BINARY}/martine-linux-arm64bits/format_data $(SOURCEDIR)/resources/formatter/data/format_data.go
	zip ${BINARY}/martine-$(appversion)-arm-64bits.zip ${BINARY}/martine-linux-arm64bits/* ./resources/*
	@echo "Compilation for windows 32bits"
	mkdir ${BINARY}/martine-windows-32bits
	GOOS=windows GOARCH=386 go build ${LDFLAGS} -o ${BINARY}/martine-windows-32bits/martine.exe $(SOURCEDIR)/main.go $(SOURCEDIR)/process.go
	GOOS=windows GOARCH=386 go build ${LDFLAGS} -o ${BINARY}/martine-windows-32bits/prepare_delta.exe $(SOURCEDIR)/resources/formatter/delta/prepare_delta.go
	GOOS=windows GOARCH=386 go build ${LDFLAGS} -o ${BINARY}/martine-windows-32bits/format_sprite.exe $(SOURCEDIR)/resources/formatter/sprites/format_sprite.go
	GOOS=windows GOARCH=386 go build ${LDFLAGS} -o ${BINARY}/martine-windows-32bits/format_data.exe $(SOURCEDIR)/resources/formatter/data/format_data.go
	zip ${BINARY}/martine-$(appversion)-windows-32bits.zip ${BINARY}/martine-windows-32bits/* ./resources/*
	@echo "Compilation for raspberry pi Raspbian 32 bits"
	mkdir ${BINARY}/martine-linux-arm32bits
	GOOS=linux GOARCH=arm GOARM=5 go build ${LDFLAGS} -o ${BINARY}/martine-linux-arm32bits/martine $(SOURCEDIR)/main.go $(SOURCEDIR)/process.go
	GOOS=linux GOARCH=arm GOARM=5 go build ${LDFLAGS} -o ${BINARY}/martine-linux-arm32bits/prepare_delta $(SOURCEDIR)/resources/formatter/delta/prepare_delta.go
	GOOS=linux GOARCH=arm GOARM=5 go build ${LDFLAGS} -o ${BINARY}/martine-linux-arm32bits/format_sprite $(SOURCEDIR)/resources/formatter/sprites/format_sprite.go
	GOOS=linux GOARCH=arm GOARM=5 go build ${LDFLAGS} -o ${BINARY}/martine-linux-arm32bits/format_data $(SOURCEDIR)/resources/formatter/data/format_data.go
	zip ${BINARY}/martine-$(appversion)-arm.zip ${BINARY}/martine-linux-arm32bits/* ./resources/*
