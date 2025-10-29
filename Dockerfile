FROM golang:1.24.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o olmv0-mcp-server ./cmd/olmv0-mcp-server

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/olmv0-mcp-server .

EXPOSE 8080

CMD ["./olmv0-mcp-server"]