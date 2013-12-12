#CGO_ENABLED=0
#GOOS=windows|linux
#GOARCH=386|amd64




TOP_PKG=aupd

CC=llvm-gcc
#GOPATH=$(shell pwd)
PWD=$(shell pwd)
GOCMD=GOPATH=$(GOPATH):$(PWD) go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOINSTALL=$(GOCMD) install
GOTEST=$(GOCMD) test
GODEP=$(GOTEST) -i
GOFMT=gofmt -w
GOGET=$(GOCMD) get
#GOPATH=$(shell pwd):$(shell pwd)/lib



all: install
	@echo ${GOPATH}
	@echo ${TOP_PKG}
	
build:
	$(GOBUILD) -o $(PWD)/bin/$(TOP_PKG)_$(shell date +%Y%m%d%H%M%S) $(TOP_PKG)
	

clean:
	$(GOCLEAN) $(TOP_PKG)
	rm -rf pkg
	
install:
	$(GOINSTALL) $(TOP_PKG)
	mv $(PWD)/bin/$(TOP_PKG) $(PWD)/bin/$(TOP_PKG)_$(shell date +%Y%m%d%H%M%S)
	
test:
	$(GOTEST) $(TOP_PKG)
	
fmt:
	$(GOFMT) src

get:
	$(GOGET) 
