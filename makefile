CC=go
RM=rm
MV=mv


SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

VERSION:=$(shell grep -m1 "version" *.go | sed 's/[", ]//g' | cut -d= -f2)

.DEFAULT_GOAL:=build

build: 
	@echo "Compilation for linux"
	@GOOS=linux go build ${LDFLAGS} -o martine $(SOURCEDIR)/main.go
	zip martine-${VERSION}-linux.zip martine
	@echo "Compilation for windows"
	@GOOS=windows go build ${LDFLAGS} -o martine.exe $(SOURCEDIR)/main.go
	zip martine-${VERSION}-windows.zip martine.exe
	@echo "Compilation for macos"
	@GOOS=darwin go build ${LDFLAGS} -o martine $(SOURCEDIR)/main.go
	zip martine-${VERSION}-macos.zip martine    