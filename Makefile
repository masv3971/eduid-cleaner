.PHONY: update clean build build-all run package deploy test authors dist

NAME 					:= eduid-cleaner
VERSION                 := $(shell cat VERSION)
LDFLAGS                 := -ldflags "-w -s --extldflags '-static'"

default: linux

linux:
		@echo building-static
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./bin/${NAME} ${LDFLAGS} ./cmd/main.go
		@echo Done
