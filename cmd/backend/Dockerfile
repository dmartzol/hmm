# Start from golang base image
FROM --platform=$BUILDPLATFORM golang:1.19-alpine AS builder
# Set the current working directory inside the container
WORKDIR /build
ENV GOOS linux
ENV CGO_ENABLED 0
ENV GOPATH /go
ENV GOCACHE /go-build
# Copy go.mod, go.sum files and download deps
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod/cache go mod download
COPY . .
# Build the Go app
RUN --mount=type=cache,target=/go/pkg/mod/cache \
    --mount=type=cache,target=/go-build \
    go build -o backend ./cmd/backend

# Start a new stage from alpine
FROM alpine:3.17.2 as prod
RUN apk --no-cache add tzdata
WORKDIR /dist
# Copy the build artifacts from the previous stage
COPY --from=builder /build/backend .
CMD ["/dist/backend"]

