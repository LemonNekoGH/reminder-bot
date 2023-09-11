FROM golang:1.20-alpine as builder

WORKDIR /app

ARG version

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY . .

RUN go build -ldflags "-X github.com/lemonnekogh/reminderbot/config.Version=${version}" -o reminder-bot cmd/reminderbot/main.go

FROM alpine:latest as runner

WORKDIR /app

COPY --from=builder /app/reminder-bot /app/reminder-bot

ENTRYPOINT [ "/app/reminder-bot" ]
