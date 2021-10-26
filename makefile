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
	(make compile ARCH=amd64 OS=linux)
	@echo "Compilation for windows"
	(make compile ARCH=amd64 OS=windows EXT=.exe)
	@echo "Compilation for macos"
	(make compile ARCH=amd64 OS=darwin )
	@echo "Compilation for raspberry pi Raspbian 64 bits"
	(make compile ARCH=arm64 OS=linux)
	@echo "Compilation for windows 32bits"
	(make compile ARCH=386 OS=windows  EXT=.exe)
	@echo "Compilation for raspberry pi Raspbian 32 bits"
	(make compile ARCH=arm OS=linux GOARM=5)


compile:
	@echo "Compilation for ${ARCH} ${OS} bits"
	mkdir ${BINARY}/martine-${OS}-${ARCH}
	GOOS=${OS} GOARCH=${ARCH} go build ${LDFLAGS} -o ${BINARY}/martine-${OS}-${ARCH}/martine${EXT} $(SOURCEDIR)/main.go $(SOURCEDIR)/process.go $(SOURCEDIR)/export_handler.go
	GOOS=${OS} GOARCH=${ARCH} go build ${LDFLAGS} -o ${BINARY}/martine-${OS}-${ARCH}/prepare_delta${EXT} $(SOURCEDIR)/resources/formatter/delta/prepare_delta.go
	GOOS=${OS} GOARCH=${ARCH} go build ${LDFLAGS} -o ${BINARY}/martine-${OS}-${ARCH}/format_sprite${EXT} $(SOURCEDIR)/resources/formatter/sprites/format_sprite.go
	GOOS=${OS} GOARCH=${ARCH} go build ${LDFLAGS} -o ${BINARY}/martine-${OS}-${ARCH}/format_data${EXT} $(SOURCEDIR)/resources/formatter/data/format_data.go
	zip ${BINARY}/martine-$(appversion)-${OS}-${ARCH}.zip ${BINARY}/martine-${OS}-${ARCH}/* ./resources/*