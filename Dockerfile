FROM rust:1.76 as builder

WORKDIR /app

COPY . .

RUN cargo build --release

FROM ubuntu:latest as runner

LABEL authors="lemonneko"

RUN apt update && apt install libpq-dev ca-certificates -y && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/target/release/reminder_bot /app/reminder_bot

ENTRYPOINT ["/app/reminder_bot"]
