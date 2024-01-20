FROM golang:1.21.6-alpine3.19 AS builder
WORKDIR /build
COPY . /build
ENV CGO_ENABLED=0  \
  GOCACHE=/.cache/go-build

RUN --mount=type=cache,target=/.cache/go-build go build -o ./service -buildvcs=false -trimpath -ldflags "-s -w" cmd/meme_bot/main.go
RUN chmod +x /build/service

#####################
FROM alpine:3.19
RUN apk update
RUN apk upgrade
RUN apk add --no-cache ffmpeg font-terminus font-inconsolata font-dejavu font-noto font-noto-cjk font-awesome font-noto-extra

COPY --from=builder /build/service /usr/bin/service

CMD ["/usr/bin/service"]
