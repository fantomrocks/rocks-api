# --------------------------------------------------------------------------
# Makefile for the Fantom Rocks API
#
# v0.1 (2020/01/29)  - Initial version, base API daemon build.
# (c) GoFantom, 2020
# --------------------------------------------------------------------------

# project related vars
PROJECT := $(shell basename "$(PWD)")

# go related vars
GOBASE := $(shell pwd)
GOBIN=$(CURDIR)/bin
GOFILES := $(wildcard cmd/*.go)

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

## server: Make the API server as bin/frd
server:
	go build -o $(GOBIN)/frd cmd/fantomrocksd.go

.PHONY: help
all: help
help: Makefile
	@echo
	@echo "Choose a make command in "$(PROJECT)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
