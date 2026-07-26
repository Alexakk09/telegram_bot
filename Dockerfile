FROM golang:1.26.5

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o weatherbot .

CMD ["./weatherbot"]# ---------- Stage 1: Build ----------
FROM golang:1.26.5 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o weatherbot .

# ---------- Stage 2: Runtime ----------
FROM debian:bookworm-slim

RUN apt-get update && \
    apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/weatherbot .



CMD ["./weatherbot"]