CC=go
RM=rm
MV=mv


SOURCEDIR=./cli
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')


VERSION:=$(shell grep -m1 "AppVersion" ./common/app.go | sed 's/[", ]//g' | cut -d= -f2)
suffix=$(shell grep -m1 "version" $(SOURCEDIR)/*.go | sed 's/[", ]//g' | cut -d= -f2 | sed 's/[0-9.]//g')
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
	(make package-linux ARCH=amd64 OS=linux)
else 
ifeq ($(UNAME), Darwin)
	(make package-darwin ARCH=amd64 OS=darwin)
else 
	(make package-windows ARCH=amd64 OS=windows EXT=.exe)
endif 
endif
	
init:
	@echo "Compilation for ${ARCH} ${OS} bits"
	mkdir ${BINARY}/martine-${OS}-${ARCH}

compile:
	GOOS=${OS} GOARCH=${ARCH} go build ${LDFLAGS} -o ${BINARY}/martine-${OS}-${ARCH}/martine${EXT} $(SOURCEDIR)/main.go $(SOURCEDIR)/process.go $(SOURCEDIR)/export_handler.go
	GOOS=${OS} GOARCH=${ARCH} go build ${LDFLAGS} -o ${BINARY}/martine-${OS}-${ARCH}/prepare_delta${EXT} ./resources/formatter/delta/prepare_delta.go
	GOOS=${OS} GOARCH=${ARCH} go build ${LDFLAGS} -o ${BINARY}/martine-${OS}-${ARCH}/format_sprite${EXT} ./resources/formatter/sprites/format_sprite.go
	GOOS=${OS} GOARCH=${ARCH} go build ${LDFLAGS} -o ${BINARY}/martine-${OS}-${ARCH}/format_data${EXT} ./resources/formatter/data/format_data.go

archive:
	zip ${BINARY}/martine-$(appversion)-${OS}-${ARCH}.zip ${BINARY}/martine-${OS}-${ARCH}/* ./resources/*

build-linux:
	@echo "Compilation for linux"
	(make init ARCH=amd64 OS=linux)
	(make compile ARCH=amd64 OS=linux)
	(make archive ARCH=amd64 OS=linux)

build-windows:
	@echo "Compilation for windows"
	(make init ARCH=amd64 OS=windows EXT=.exe)
	(make compile ARCH=amd64 OS=windows EXT=.exe)
	(make archive ARCH=amd64 OS=windows EXT=.exe)

build-darwin:
	@echo "Compilation for macos"
	(make init ARCH=arm64 OS=darwin)
	(make compile ARCH=amd64 OS=darwin)
	(make archive ARCH=amd64 OS=darwin)

build-raspbian:
	@echo "Compilation for raspberry pi Raspbian 64 bits"
	(make init ARCH=arm64 OS=linux)
	(make compile ARCH=arm64 OS=linux)
	(make archive  ARCH=arm64 OS=linux)

build-raspbian-i386:
	@echo "Compilation for raspberry pi Raspbian 32 bits"
	(make init ARCH=arm OS=linux GOARM=5)
	(make compile ARCH=arm OS=linux GOARM=5)
	(make archive ARCH=arm OS=linux GOARM=5)

build-windows-i386:
	@echo "Compilation for windows 32bits"
	(make init ARCH=386 OS=windows  EXT=.exe)
	(make compile ARCH=386 OS=windows  EXT=.exe)
	(make archive ARCH=386 OS=windows  EXT=.exe)

package-darwin:
	(make init)
	@echo "Compilation and packaging for darwin"
	fyne package -os darwin -icon martine-logo.png -sourceDir ${SOURCEDIR} -name martine
	mv martine.app ${BINARY}/martine-${OS}-${ARCH}/
	(make archive)

package-windows:
	(make init)
	@echo "Compilation and packaging for windows"
	fyne package -os windows -icon martine-logo.png -sourceDir ${SOURCEDIR} -name martine
	mv martine.exe ${BINARY}/martine-${OS}-${ARCH}/
	(make archive)

package-linux:
	(make init)
	@echo "Compilation and packaging for linux"
	fyne package -os linux -icon martine-logo.png -sourceDir ${SOURCEDIR} -name martine
	mv martine ${BINARY}/martine-${OS}-${ARCH}/
	(make archive)
		
get-linter:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.50.1

lint:
	@echo "Lint the whole project"
	golangci-lint run --timeout 2m 
