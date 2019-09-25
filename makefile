CC=go
RM=rm
MV=mv


SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

VERSION:=$(shell grep -m1 "version" *.go | sed 's/[", ]//g' | cut -d= -f2)

.DEFAULT_GOAL:=build

SNAPSHOT=`date +%FT%T%z`
STATUS=`echo ${VERSION} | grep ".rc"`


build: 
	@echo "Compilation for linux"
	@if [ "${STATUS}" = "${VERSION}" ]; then echo "Version :${VERSION} is snapshot" && suffix=${SNAPSHOT} ; else echo "Version :${VERSION} is release" && suffix=""; fi
	@echo $$suffix
	@GOOS=linux go build ${LDFLAGS} -o martine $(SOURCEDIR)/main.go $(SOURCEDIR)/process.go
	zip martine-${VERSION}$$suffix-linux.zip martine ./resources/*
	@echo "Compilation for windows"
	@GOOS=windows go build ${LDFLAGS} -o martine.exe $(SOURCEDIR)/main.go $(SOURCEDIR)/process.go
	zip martine-${VERSION}$$suffix-windows.zip martine.exe ./resources/*
	@echo "Compilation for macos"
	@GOOS=darwin go build ${LDFLAGS} -o martine $(SOURCEDIR)/main.go $(SOURCEDIR)/process.go
	zip martine-${VERSION}$$suffix-macos.zip martine  ./resources/*
