export VERSION = v0.0.1


# local GOOS
GOOS_LOCAL ?= $(shell uname -s | tr '[:upper:]' '[:lower:]')
HOSTARCH ?= $(shell uname -m)
KERNEL_NAME ?= $(shell uname -s)
ifeq ($(HOSTARCH),x86_64)
	ARCH ?= amd64
	BUILDPLATFORM ?= linux/amd64
else
 ifeq ($(HOSTARCH),aarch64)
	ARCH := arm64
	BUILDPLATFORM ?= linux/arm64/v8
 endif
endif

include makefiles/const.mk

# 打印版本号
.PHONY: all
.EXPORT_ALL_VARIABLES:
all: clean go-build rpm-build deb-build

.PHONY: help
help:           ## help
	@sed -ne '/@sed/!s/## //p' $(MAKEFILE_LIST)

# 构建rpm 包
.PHONY: go-build
go-build:     ## 构建rpm 包
	go mod tidy \
	&&  $(GOBUILD_ENV) go build -a -ldflags $(LDFLAGS) -o out/infra-pulse

# 构建rpm 包
.PHONY: rpm-build
rpm-build:     ## 构建rpm 包
	go-bin-rpm generate -f service/centos7-${HOSTARCH}/rpm.json -o ./out/infra-pulse-${VERSION}-${KERNEL_NAME}-${HOSTARCH}.rpm


# 构建deb 包
.PHONY: deb-build
deb-build:    ## 构建 deb 包
	go-bin-deb generate -f service/ubuntu18-${HOSTARCH}/deb.json -o ./out/infra-pulse-${VERSION}-${KERNEL_NAME}-${HOSTARCH}.deb

.PHONY: clean
clean:       ## 清理安装包
	rm -rf out/ pkg-build/