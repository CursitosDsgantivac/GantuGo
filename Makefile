USE_ENV_FILE = true
ENV ?= test
ARG ?= chiServer


ifeq ($(USE_ENV_FILE), true)
    ifeq (,$(wildcard .env.$(ENV)))
        $(error .env.$(ENV) file not found!)
    endif

    include .env.$(ENV)
    export
endif

.DEFAULT_GOAL := clean

.PHONY:fmt vet build run clean

fmt:
	go fmt ./...

vet: fmt 
	go vet ./...

build: vet
	go build -o gantuProgram

run: build
	./gantuProgram ${ARG}

clean: run
	rm -f gantuProgram

