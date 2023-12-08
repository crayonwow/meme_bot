FROM golang:1.21-alpine AS builder
WORKDIR /build
COPY . /build
ENV CGO_ENABLED=0  \
  GOCACHE=/.cache/go-build

RUN --mount=type=cache,target=/.cache/go-build go build -o ./bot_bin -buildvcs=false -trimpath -ldflags "-s -w" cmd/meme_bot/main.go
RUN chmod +x /build/bot_bin

FROM scratch

COPY --from=builder /build/bot_bin /usr/bin/bot_bin

ARG TOKEN
ARG PORT
ARG CHAT_ID

ENV TOKEN=$TOKEN \
  PORT=$PORT \
  CHAT_ID=$CHAT_ID 

CMD ["/usr/bin/bot_bin"]

EXPOSE 8080
