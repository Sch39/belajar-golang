# Stage 1: Build frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /app/web
COPY web/package*.json ./
RUN npm install --legacy-peer-deps
COPY web/ ./
# Build frontend ke folder /public
RUN npm run build

# Stage 2: Build backend
FROM golang:1.24-alpine AS backend-builder
WORKDIR /app

# Install swag CLI
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy go.mod dan go.sum dulu untuk caching
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source code backend
COPY . .

# Generate Swagger docs
RUN swag init -g cmd/server/main.go --output ./docs --parseDependency --parseInternal

# Build Go binary
ARG BINARY_NAME=my-kasir-gw
ARG GO_MAIN=cmd/server/main.go
RUN go build -o ${BINARY_NAME} ${GO_MAIN}

# Stage 3: Final image
FROM alpine:3.18
WORKDIR /app
RUN apk add --no-cache ca-certificates

# Copy binary, Swagger docs, dan frontend build
COPY --from=backend-builder /app/my-kasir-gw .
COPY --from=backend-builder /app/docs ./docs
COPY --from=frontend-builder /app/public ./public

EXPOSE 8080
CMD ["./my-kasir-gw"]
