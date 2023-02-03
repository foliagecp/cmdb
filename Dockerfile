FROM golang:1.19 as builder

LABEL maintainer="NJWS, Inc."

WORKDIR /src/

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=0.1.1 -X main.release=$(git rev-parse --short HEAD)" -o /build/cmdb ./cmd/cmdb

FROM ubuntu:18.04

LABEL maintainer="NJWS, Inc."

RUN apt update && \
    apt install ca-certificates netcat -y && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /build/cmdb /usr/bin/

RUN chmod +x /usr/bin/cmdb

ENTRYPOINT ["/usr/bin/cmdb"]
