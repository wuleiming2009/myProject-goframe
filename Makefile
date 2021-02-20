SHELL := /bin/bash
GO     ?= go

define check_command
	command -v $(1) 1>/dev/null || $(2)
endef

# build项目
.PHONY: build
build:
	cd src;go build -v -o ../bin/myProject

.PHONY: all
all: build swag

.PHONY: clean
clean:
	rm -rf bin/*

# 生成swagger文档
.PHONY: swag
swag:
	@@$(call check_command,swag,$(GO) get -u github.com/swaggo/swag/cmd/swag)
	cd src;swag init

check:
	$(call check_comand,asdf,ddd)