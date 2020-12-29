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
	@echo "Compilation for linux"
	GOOS=linux go build ${LDFLAGS} -o martine $(SOURCEDIR)/main.go $(SOURCEDIR)/process.go
	GOOS=linux go build ${LDFLAGS} -o prepare_delta $(SOURCEDIR)/resources/formatter/delta/prepare_delta.go
	GOOS=linux go build ${LDFLAGS} -o format_sprite $(SOURCEDIR)/resources/formatter/sprites/format_sprite.go
	GOOS=linux go build ${LDFLAGS} -o format_data $(SOURCEDIR)/resources/formatter/data/format_data.go
	zip martine-$(appversion)-linux.zip martine prepare_delta format_sprite format_data ./resources/*
	@echo "Compilation for windows"
	GOOS=windows go build ${LDFLAGS} -o martine.exe $(SOURCEDIR)/main.go $(SOURCEDIR)/process.go
	GOOS=windows go build ${LDFLAGS} -o prepare_delta.exe $(SOURCEDIR)/resources/formatter/delta/prepare_delta.go
	GOOS=windows go build ${LDFLAGS} -o format_sprite.exe $(SOURCEDIR)/resources/formatter/sprites/format_sprite.go
	GOOS=windows go build ${LDFLAGS} -o format_data.exe $(SOURCEDIR)/resources/formatter/data/format_data.go
	zip martine-$(appversion)-windows.zip martine.exe  prepare_delta.exe format_sprite.exe format_data.exe ./resources/*
	@echo "Compilation for macos"
	GOOS=darwin go build ${LDFLAGS} -o martine $(SOURCEDIR)/main.go $(SOURCEDIR)/process.go
	GOOS=darwin go build ${LDFLAGS} -o prepare_delta $(SOURCEDIR)/resources/formatter/delta/prepare_delta.go
	GOOS=darwin go build ${LDFLAGS} -o format_sprite $(SOURCEDIR)/resources/formatter/sprites/format_sprite.go
	GOOS=darwin go build ${LDFLAGS} -o format_data $(SOURCEDIR)/resources/formatter/data/format_data.go
	zip martine-$(appversion)-macos.zip martine  prepare_delta format_sprite format_data ./resources/*
	@echo "Compilation for raspberry pi Raspbian"
	GOOS=linux ARCH=arm GOARM=5 go build ${LDFLAGS} -o martine $(SOURCEDIR)/main.go $(SOURCEDIR)/process.go
	GOOS=linux ARCH=arm GOARM=5 go build ${LDFLAGS} -o prepare_delta $(SOURCEDIR)/resources/formatter/delta/prepare_delta.go
	GOOS=linux ARCH=arm GOARM=5 go build ${LDFLAGS} -o format_sprite $(SOURCEDIR)/resources/formatter/sprites/format_sprite.go
	GOOS=linux ARCH=arm GOARM=5 go build ${LDFLAGS} -o format_data $(SOURCEDIR)/resources/formatter/data/format_data.go
	zip martine-$(appversion)-arm.zip martine  prepare_delta format_sprite format_data ./resources/*
	@echo "Compilation for windows 32bits"
	GOOS=windows GOARCH=386 go build ${LDFLAGS} -o martine.exe $(SOURCEDIR)/main.go $(SOURCEDIR)/process.go
	GOOS=windows GOARCH=386 go build ${LDFLAGS} -o prepare_delta.exe $(SOURCEDIR)/resources/formatter/delta/prepare_delta.go
	GOOS=windows GOARCH=386 go build ${LDFLAGS} -o format_sprite.exe $(SOURCEDIR)/resources/formatter/sprites/format_sprite.go
	GOOS=windows GOARCH=386 go build ${LDFLAGS} -o format_data.exe $(SOURCEDIR)/resources/formatter/data/format_data.go
	zip martine-$(appversion)-windows-32bits.zip martine.exe  prepare_delta.exe format_sprite.exe format_data.exe ./resources/*