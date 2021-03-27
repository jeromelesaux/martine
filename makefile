CC=go
RM=rm
MV=mv


SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

VERSION:=$(shell grep -m1 "appVersion" *.go | sed 's/[", ]//g' | cut -d= -f2)
suffix=$(shell grep -m1 "version" *.go | sed 's/[", ]//g' | cut -d= -f2 | sed 's/[0-9.]//g')
snapshot=$(shell date +%FT%T)

ifeq ($(suffix),rc)
	appversion=$(VERSION)$(snapshot)
else 
	appversion=$(VERSION)
endif 

.DEFAULT_GOAL:=build


build: 
	rm -fr martine*
	@echo "Compilation for linux"
	mkdir martine-linux-64bits
	export GOOS=linux && go build ${LDFLAGS} -o martine-linux-64bits/martine $(SOURCEDIR)/main.go $(SOURCEDIR)/process.go
	export GOOS=linux && go build ${LDFLAGS} -o martine-linux-64bits/prepare_delta $(SOURCEDIR)/resources/formatter/delta/prepare_delta.go
	export GOOS=linux && go build ${LDFLAGS} -o martine-linux-64bits/format_sprite $(SOURCEDIR)/resources/formatter/sprites/format_sprite.go
	export GOOS=linux && go build ${LDFLAGS} -o martine-linux-64bits/format_data $(SOURCEDIR)/resources/formatter/data/format_data.go
	zip martine-$(appversion)-linux.zip martine-linux-64bits/* ./resources/*
	@echo "Compilation for windows"
	mkdir martine-windows-64bits
	export GOOS=windows && go build ${LDFLAGS} -o martine-windows-64bits/martine.exe $(SOURCEDIR)/main.go $(SOURCEDIR)/process.go
	export GOOS=windows && go build ${LDFLAGS} -o martine-windows-64bits/prepare_delta.exe $(SOURCEDIR)/resources/formatter/delta/prepare_delta.go
	export GOOS=windows && go build ${LDFLAGS} -o martine-windows-64bits/format_sprite.exe $(SOURCEDIR)/resources/formatter/sprites/format_sprite.go
	export GOOS=windows && go build ${LDFLAGS} -o martine-windows-64bits/format_data.exe $(SOURCEDIR)/resources/formatter/data/format_data.go
	zip martine-$(appversion)-windows.zip  martine-windows-64bits/* ./resources/*
	@echo "Compilation for macos"
	mkdir martine-darwin-64bits 
	export GOOS=darwin && go build ${LDFLAGS} -o martine-darwin-64bits/martine $(SOURCEDIR)/main.go $(SOURCEDIR)/process.go
	export GOOS=darwin && go build ${LDFLAGS} -o martine-darwin-64bits/prepare_delta $(SOURCEDIR)/resources/formatter/delta/prepare_delta.go
	export GOOS=darwin && go build ${LDFLAGS} -o martine-darwin-64bits/format_sprite $(SOURCEDIR)/resources/formatter/sprites/format_sprite.go
	export GOOS=darwin && go build ${LDFLAGS} -o martine-darwin-64bits/format_data $(SOURCEDIR)/resources/formatter/data/format_data.go
	zip martine-$(appversion)-macos.zip martine-darwin-64bits/* ./resources/*
	@echo "Compilation for raspberry pi Raspbian 64 bits"
	mkdir martine-linux-arm64bits
	export GOOS=linux && export GOARCH=arm64 && go build ${LDFLAGS} -o martine-linux-arm64bits/martine $(SOURCEDIR)/main.go $(SOURCEDIR)/process.go
	export GOOS=linux && export GOARCH=arm64 && go build ${LDFLAGS} -o martine-linux-arm64bits/prepare_delta $(SOURCEDIR)/resources/formatter/delta/prepare_delta.go
	export GOOS=linux && export GOARCH=arm64 && go build ${LDFLAGS} -o martine-linux-arm64bits/format_sprite $(SOURCEDIR)/resources/formatter/sprites/format_sprite.go
	export GOOS=linux && export GOARCH=arm64 && go build ${LDFLAGS} -o martine-linux-arm64bits/format_data $(SOURCEDIR)/resources/formatter/data/format_data.go
	zip martine-$(appversion)-arm-64bits.zip martine-linux-arm64bits/* ./resources/*
	@echo "Compilation for windows 32bits"
	mkdir martine-windows-32bits
	export GOOS=windows && export GOARCH=386 && go build ${LDFLAGS} -o martine-windows-32bits/martine.exe $(SOURCEDIR)/main.go $(SOURCEDIR)/process.go
	export GOOS=windows && export GOARCH=386 && go build ${LDFLAGS} -o martine-windows-32bits/prepare_delta.exe $(SOURCEDIR)/resources/formatter/delta/prepare_delta.go
	export GOOS=windows && export GOARCH=386 && go build ${LDFLAGS} -o martine-windows-32bits/format_sprite.exe $(SOURCEDIR)/resources/formatter/sprites/format_sprite.go
	export GOOS=windows && export GOARCH=386 && go build ${LDFLAGS} -o martine-windows-32bits/format_data.exe $(SOURCEDIR)/resources/formatter/data/format_data.go
	zip martine-$(appversion)-windows-32bits.zip martine-windows-32bits/* ./resources/*
	@echo "Compilation for raspberry pi Raspbian 32 bits"
	mkdir martine-linux-arm32bits
	export GOOS=linux && GOARCH=arm && GOARM=5 && go build ${LDFLAGS} -o martine-linux-arm32bits/martine $(SOURCEDIR)/main.go $(SOURCEDIR)/process.go
	export GOOS=linux && GOARCH=arm && GOARM=5 && go build ${LDFLAGS} -o martine-linux-arm32bits/prepare_delta $(SOURCEDIR)/resources/formatter/delta/prepare_delta.go
	export GOOS=linux && GOARCH=arm && GOARM=5 && go build ${LDFLAGS} -o martine-linux-arm32bits/format_sprite $(SOURCEDIR)/resources/formatter/sprites/format_sprite.go
	export GOOS=linux && GOARCH=arm && GOARM=5 && go build ${LDFLAGS} -o martine-linux-arm32bits/format_data $(SOURCEDIR)/resources/formatter/data/format_data.go
	zip martine-$(appversion)-arm.zip martine-linux-arm32bits/* ./resources/*
