FROM golang:1.20rc1-alpine3.17 AS dependencies

RUN apk update && apk add --no-cache \
  ca-certificates \
  musl-dev \
  gcc \
  build-base \
  git

ENV GO111MODULE=on \
  CGO_ENABLED=1 \
  GOOS=linux \
  GOARCH=arm64

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

FROM dependencies AS builder

COPY . .

RUN ./build.sh

WORKDIR /dist

RUN cp /build/main .

FROM scratch as dist

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /dist/main /

ENTRYPOINT ["/main"]
