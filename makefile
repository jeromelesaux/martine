CC=go
RM=rm
MV=mv


SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

VERSION:=$(shell grep -m1 "appVersion" *.go | sed 's/[", ]//g' | cut -d= -f2)
suffix=$(shell grep -m1 "version" *.go | sed 's/[", ]//g' | cut -d= -f2 | sed 's/[0-9.]//g')
snapshot=$(shell date +%FT%T)
UNAME := $(shell uname)
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
ifeq ($(UNAME),Linux)
	(make build-linux)
else 
ifeq ($(UNAME), Darwin)
	(make build-darwin)
else 
	(make build-windows)
endif 
endif

package:
	rm -fr ${BINARY}
	mkdir ${BINARY}
ifeq ($(UNAME),Linux)
	(make package-linux)
else 
ifeq ($(UNAME), Darwin)
	(make package-darwin)
else 
	(make package-windows)
endif 
endif
	
compile:
	@echo "Compilation for ${ARCH} ${OS} bits"
	mkdir ${BINARY}/martine-${OS}-${ARCH}
	GOOS=${OS} GOARCH=${ARCH} go build ${LDFLAGS} -o ${BINARY}/martine-${OS}-${ARCH}/martine${EXT} $(SOURCEDIR)/main.go $(SOURCEDIR)/process.go $(SOURCEDIR)/export_handler.go
	GOOS=${OS} GOARCH=${ARCH} go build ${LDFLAGS} -o ${BINARY}/martine-${OS}-${ARCH}/prepare_delta${EXT} $(SOURCEDIR)/resources/formatter/delta/prepare_delta.go
	GOOS=${OS} GOARCH=${ARCH} go build ${LDFLAGS} -o ${BINARY}/martine-${OS}-${ARCH}/format_sprite${EXT} $(SOURCEDIR)/resources/formatter/sprites/format_sprite.go
	GOOS=${OS} GOARCH=${ARCH} go build ${LDFLAGS} -o ${BINARY}/martine-${OS}-${ARCH}/format_data${EXT} $(SOURCEDIR)/resources/formatter/data/format_data.go
	zip ${BINARY}/martine-$(appversion)-${OS}-${ARCH}.zip ${BINARY}/martine-${OS}-${ARCH}/* ./resources/*

build-linux:
	@echo "Compilation for linux"
	mkdir ${BINARY}/martine-${OS}-${ARCH}
	(make compile ARCH=amd64 OS=linux)

build-windows:
	@echo "Compilation for windows"
	mkdir ${BINARY}/martine-${OS}-${ARCH}
	(make compile ARCH=amd64 OS=windows EXT=.exe)

build-darwin:
	@echo "Compilation for macos"
	mkdir ${BINARY}/martine-${OS}-${ARCH}
	(make compile ARCH=amd64 OS=darwin)

build-raspbian:
	@echo "Compilation for raspberry pi Raspbian 64 bits"
	mkdir ${BINARY}/martine-${OS}-${ARCH}
	(make compile ARCH=arm64 OS=linux)

build-raspbian-i386:
	@echo "Compilation for raspberry pi Raspbian 32 bits"
	mkdir ${BINARY}/martine-${OS}-${ARCH}
	(make compile ARCH=arm OS=linux GOARM=5)

build-windows-i386:
	@echo "Compilation for windows 32bits"
	mkdir ${BINARY}/martine-${OS}-${ARCH}
	(make compile ARCH=386 OS=windows  EXT=.exe)

package-darwin:
	@echo "Compilation and packaging for darwin"
	fyne package -os darwin -icon martine-logo.png -sourceDir ./

package-windows:
	@echo "Compilation and packaging for darwin"
	fyne package -os windows -icon martine-logo.png -sourceDir ./

package-linux:
	@echo "Compilation and packaging for darwin"
	fyne package -os linux -icon martine-logo.png -sourceDir ./
		