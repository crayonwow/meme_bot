FROM golang:1.21-alpine AS builder
WORKDIR /build
COPY . /build
ENV CGO_ENABLED=0  \
  GOCACHE=/.cache/go-build

RUN --mount=type=cache,target=/.cache/go-build go build -o ./service -buildvcs=false -trimpath -ldflags "-s -w" cmd/meme_bot/main.go
RUN chmod +x /build/service
#####################
#####################
FROM alpine:3.19
RUN apk update
RUN apk upgrade
RUN apk add --no-cache ffmpeg

COPY --from=builder /build/service /usr/bin/service

CMD ["/usr/bin/service"]
