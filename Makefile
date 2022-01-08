
# バージョン
VERSION:=$(shell cat VERSION)

# リビジョン
REVISION:=$(shell git rev-parse --short HEAD 2> /dev/null || cat REVISION)

## コマンド定義
GO:=go
GO_BUILD:=$(GO) build
BINDIR:=bin
GOOS:=windows
GOARCH:=amd64
ENV:=GOOS=$(GOOS) GOARCH=$(GOARCH)

ROOT_PACKAGE:=$(shell go list .)

COMMAND_PACKAGES:=$(shell go list ./cmd/...)

BINARIES:=$(COMMAND_PACKAGES:$(ROOT_PACKAGE)/cmd/%=$(BINDIR)/%)

GO_FILES:=$(shell find . -type f -name '*.go' -print)

# ldflag
GO_LDFLAGS_VERSION:=-X '${ROOT_PACKAGE}.VERSION=${VERSION}' -X '${ROOT_PACKAGE}.REVISION=${REVISION}'
GO_LDFLAGS:=$(GO_LDFLAGS_VERSION)

# go build
GO_BUILD_OPTION:=-ldflags "$(GO_LDFLAGS)"

# ビルドタスク
.PHONY: build
build: $(BINARIES)

# 実ビルドタスク
$(BINARIES): $(GO_FILES) VERSION .git/HEAD
	$(ENV) $(GO_BUILD) -o $@ $(GO_BUILD_OPTION) $(@:$(BINDIR)/%=$(ROOT_PACKAGE)/cmd/%)
	cp $@ $@.exe


