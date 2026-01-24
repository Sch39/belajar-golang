# Stage 1: Build frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /app/web
COPY web/package*.json ./
RUN npm install --legacy-peer-deps
COPY web/ ./
RUN npm run build

# Stage 2: Build backend
FROM golang:1.21-alpine AS backend-builder
WORKDIR /app
# Install swag CLI
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy backend source
COPY cmd ./cmd
COPY go.mod go.sum ./
RUN go mod download

# Copy frontend build
COPY --from=frontend-builder /app/web/dist ./web/dist
COPY --from=frontend-builder /app/web/public ./web/public

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
COPY --from=backend-builder /app/my-kasir-gw .
COPY --from=backend-builder /app/docs ./docs
COPY --from=backend-builder /app/web/dist ./web/dist
COPY --from=backend-builder /app/web/public ./web/public

EXPOSE 8080
CMD ["./my-kasir-gw"]
