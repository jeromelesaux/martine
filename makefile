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
	zip martine-$(appversion)-linux.zip martine ./resources/*
	@echo "Compilation for windows"
	GOOS=windows go build ${LDFLAGS} -o martine.exe $(SOURCEDIR)/main.go $(SOURCEDIR)/process.go
	zip martine-$(appversion)-windows.zip martine.exe ./resources/*
	@echo "Compilation for macos"
	GOOS=darwin go build ${LDFLAGS} -o martine $(SOURCEDIR)/main.go $(SOURCEDIR)/process.go
	zip martine-$(appversion)-macos.zip martine  ./resources/*
	@echo "Compilation for raspberry pi Raspbian"
	GOOS=linux ARCH=arm GOARM=5 go build ${LDFLAGS} -o martine $(SOURCEDIR)/main.go $(SOURCEDIR)/process.go
	zip martine-$(appversion)-arm.zip martine  ./resources/*