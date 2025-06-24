SHELL := /bin/bash

GOBUILD_ENV = GO111MODULE=on CGO_ENABLED=0 GOARCH=$(ARCH) GOOS=$(GOOS_LOCAL)


TARGETS     := darwin/amd64 linux/amd64 windows/amd64 linux/arm64
DIST_DIRS   := find * -type d -exec

TIME_LONG	= `date +%Y-%m-%d' '%H:%M:%S`
BUILD_TIME	= `date +%Y-%m-%d-%H:%M:%S`
TIME_SHORT	= `date +%H:%M:%S`
BUILD_TIME1 = echo ${BUILD_TIME}

NEXUS = 'devops_test:Cloud_dev123'

#NEXUS_PASSWORD = echo ${NEXUS_PASSWORD}

BLUE         := $(shell printf "\033[34m")
YELLOW       := $(shell printf "\033[33m")
RED          := $(shell printf "\033[31m")
GREEN        := $(shell printf "\033[32m")
CNone        := $(shell printf "\033[0m")

INFO	= echo ${TIME} ${BLUE}[ .. ]${CNone}
WARN	= echo ${TIME} ${YELLOW}[WARN]${CNone}
ERR		= echo ${TIME} ${RED}[FAIL]${CNone}
OK		= echo ${TIME} ${GREEN}[ OK ]${CNone}
FAIL	= (echo ${TIME} ${RED}[FAIL]${CNone} && false)


# Repo info
GIT_COMMIT          ?= git-$(shell git rev-parse --short HEAD)

GIT_BRANCH ?=$(shell git rev-parse --abbrev-ref HEAD)

ifeq ("master",$(VERSION))
else
VERSION ?=$(GIT_BRANCH)-git-$(GIT_COMMIT)
endif


ifeq (,$(shell git config --get user.name))
ACCESS_TOKEN_USR=$(GIT_USR)
else
ACCESS_TOKEN_USR=$(shell git config --get user.name)
endif

ifeq (,$(shell git config --get user.password))
ACCESS_TOKEN_PWD=$(GIT_PASSWORD)
else
ACCESS_TOKEN_PWD=$(shell git config --get user.password)
endif

#	buildVersion     = "unknown"
#	buildGitRevision = "unknown"
#	buildStatus      = "unknown"
#	buildTag         = "unknown"
#	buildHub         = "unknown"

GIT_COMMIT_LONG     ?= $(shell git rev-parse HEAD)
GIT_BRANCH ?= $(shell git branch --show-current)
GO_MODULE_NAME ?= $(shell go list)
BUILD_VERSION_KEY ?= "$(GO_MODULE_NAME)/src/version.buildVersion"
BUILD_GIT_REVISION_KEY ?= "$(GO_MODULE_NAME)/src/version.buildGitRevision"
BUILD_STATUS_KEY ?= "$(GO_MODULE_NAME)/src/version.buildStatus"
BUILD_TAG_KEY ?= "$(GO_MODULE_NAME)/src/version.buildTag"
BUILD_HUB_KEY ?= "$(GO_MODULE_NAME)/src/version.buildHub"
BUILD_GOLANG_VERSION ?= "$(GO_MODULE_NAME)/src/version.golangVersion"
BUILD_DATE_KEY ?= "$(GO_MODULE_NAME)/src/version.buildDate"
PLATFORM_KEY ?= "$(GO_MODULE_NAME)/src/version.platform"

BUILD_VERSION ?= $(VERSION)
BUILD_GIT_REVISION ?= $(shell git rev-parse HEAD)
# GOLANG_VERSION ?= $(shell go version|awk '{print $3}')
GOLANG_VERSION ?= go1.23.4
BUILD_STATUS ?=  $(shell git rev-parse --abbrev-ref HEAD)
BUILD_TAG ?= $(shell git describe --tags $(git rev-list --tags --max-count=1))
BUILD_HUB ?= $(REGISTRY)
BUILD_TIME_VALUE ?= $(shell $(BUILD_TIME1))


LDFLAGS ?= "-s -w -X $(BUILD_VERSION_KEY)=$(BUILD_VERSION) -X $(BUILD_DATE_KEY)=$(BUILD_TIME_VALUE) -X $(PLATFORM_KEY)=$(BUILDPLATFORM) -X $(BUILD_GOLANG_VERSION)=$(GOLANG_VERSION) -X $(BUILD_GIT_REVISION_KEY)=$(BUILD_GIT_REVISION) -X $(BUILD_STATUS_KEY)=$(BUILD_STATUS) -X $(BUILD_TAG_KEY)=$(BUILD_TAG) -X $(BUILD_HUB_KEY)=$(BUILD_HUB)"

# GOPROXY := $(shell go env GORPOXY)
GOPROXY := https://goproxy.cn,direct

ifeq (,$(GOPROXY))
GOPROXY = https://goproxy.cn
endif

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif
