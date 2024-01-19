FROM golang:1.21-alpine AS builder
WORKDIR /build
COPY . /build
ENV CGO_ENABLED=0  \
  GOCACHE=/.cache/go-build

RUN --mount=type=cache,target=/.cache/go-build go build -o ./bot -buildvcs=false -trimpath -ldflags "-s -w" cmd/meme_bot/main.go
RUN chmod +x /build/bot
#####################
FROM linuxserver/ffmpeg:latest AS ffmpeg
#####################
FROM scratch
COPY --from=builder /build/bot /usr/bin/bot
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=ffmpeg /usr/local/bin/ffmpeg /usr/local/bin/ffmpeg

CMD ["/usr/bin/bot"]
