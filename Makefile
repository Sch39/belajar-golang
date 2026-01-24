BINARY_NAME=my-kasir-gw
GO_MAIN=cmd/server/main.go
FRONTEND_DIR=web
BUILD_DIR=public

# Deteksi OS
ifeq ($(OS),Windows_NT)
    BINARY_FINAL=$(BINARY_NAME).exe
    RM=del /Q
    MKDIR=mkdir
else
    BINARY_FINAL=$(BINARY_NAME)
    RM=rm -f
    MKDIR=mkdir -p
endif

.PHONY: all build build-fe build-be build-swagger run run-dev clean help install-swag

all: build

## build: Build frontend + swagger + backend
build: build-fe build-swagger build-be

## build-fe: Build frontend
build-fe:
	@echo "Building Frontend..."
	cd $(FRONTEND_DIR) && npm install --legacy-peer-deps && npm run build

## install-swag: Install swag CLI jika belum ada
install-swag:
	@which swag >/dev/null || go install github.com/swaggo/swag/cmd/swag@latest

## build-be: Build backend (OS-aware)
build-be:
	@echo "Building Backend..."
	go build -o $(BINARY_FINAL) $(GO_MAIN)

## build-swagger: Generate Swagger docs
build-swagger: install-swag
	@echo "Generating Swagger Documentation..."
	swag init -g $(GO_MAIN) --output ./docs --parseDependency --parseInternal

## run: Jalankan binary
run:
	@echo "Running the application..."
	./$(BINARY_FINAL)

## run-dev: Jalankan frontend + backend tanpa build binary
run-dev:
	@echo "Running Backend in dev mode..."
	go run $(GO_MAIN)

## clean: Bersihkan hasil build
clean:
	@echo "Cleaning up..."
	$(RM) $(BINARY_FINAL)
	-@rm -rf $(BUILD_DIR) docs 2>nul || true

help:
	@echo "Daftar perintah:"
	@grep -E '^##' $(MAKEFILE_LIST) | sed -e 's/## //g' | column -t -s ':'
