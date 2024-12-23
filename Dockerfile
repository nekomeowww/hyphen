# syntax=docker/dockerfile:1

# --- builder ---
FROM golang:1.23 as builder

RUN GO111MODULE=on

RUN mkdir /app

COPY . /app/hyphen

WORKDIR /app/hyphen

RUN go env
RUN go env -w CGO_ENABLED=0
RUN go mod download
RUN go build -a -o "release/hyphen" "github.com/nekomeowww/hyphen/cmd/hyphen"

# --- runner ---
FROM debian as runner

RUN apt update && apt upgrade -y && apt install -y ca-certificates curl && update-ca-certificates

COPY --from=builder /app/hyphen/release/hyphen /usr/local/bin/

VOLUME [ "/app/hyphen/logs" ]

EXPOSE 9010

ENTRYPOINT [ "/usr/local/bin/hyphen", "run", "--data", "/var/lib/hyphen/bbolt_data.db", "-l", "0.0.0.0:9010" ]
