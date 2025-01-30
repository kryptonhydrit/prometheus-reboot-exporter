BIN_DIR := $(shell pwd)/bin
BIN_NAME := reboot_exporter

APP_VERSION := $(shell cat VERSION)
APP_BRANCH := $(shell git describe --all --contains --dirty HEAD)
APP_REV := $(shell git rev-parse HEAD)

build:
	@echo "Build ${BIN_NAME} (Version: ${APP_VERSION})"
	GOARCH=amd64 GOOS=linux go build \
		-ldflags "-X github.com/prometheus/common/version.Version=${APP_VERSION} \
		-X github.com/prometheus/common/version.Revision=$(APP_REV) \
		-X github.com/prometheus/common/version.Branch=$(APP_BRANCH)" \
		-o ${BIN_DIR}/${BIN_NAME} .

run: build
	@echo "Starting ${BIN_NAME} (Version: ${APP_VERSION})"
	${BIN_DIR}/./${BIN_NAME}

clean:
	go clean
	rm ${BIN_DIR}/${BIN_NAME}