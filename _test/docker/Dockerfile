FROM golang:1.19-buster as builder

WORKDIR /bot

COPY go.* ./
RUN go mod download

COPY . ./
RUN go build -v -o test

FROM debian:buster-slim
RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /bot/test /bot/test

CMD ["/bot/test"]
