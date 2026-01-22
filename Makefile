BINARY_NAME=my-kasir-gw
GO_MAIN=cmd/server/main.go
FRONTEND_DIR=web
BUILD_DIR=public

# Deteksi OS
ifeq ($(OS),Windows_NT)
    BINARY_FINAL=$(BINARY_NAME).exe
else
    BINARY_FINAL=$(BINARY_NAME)
endif

.PHONY: all build build-fe build-be clean help

all: build

## build: Melakukan build FE dan BE secara berurutan
build: build-fe build-be build-swagger

## build-fe: Build frontend
build-fe:
	@echo "Building Frontend..."
	cd $(FRONTEND_DIR) && npm ci && npm run build

## build-be: Build backend (menggunakan BINARY_FINAL)
build-be:
	@echo "Building Backend..."
	go build -o $(BINARY_FINAL) $(GO_MAIN)

build-swagger:
	@echo "Generating Swagger Documentation..."
	swag init -g $(GO_MAIN) --output ./docs

## clean: Membersihkan hasil build
clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME) $(BINARY_NAME).exe
	rm -rf $(BUILD_DIR)

help:
	@echo "Daftar perintah:"
	@grep -E '^##' $(MAKEFILE_LIST) | sed -e 's/## //g' | column -t -s ':'