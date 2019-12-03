FROM golang:buster AS builder
RUN update-ca-certificates
WORKDIR /
COPY . .
RUN go build -o /wallet ./cmd/wallet/main.go

# Image to use
FROM debian:buster-slim
COPY --from=builder /wallet /wallet
ADD ./client.toml /client.toml
 
CMD /wallet
