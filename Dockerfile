FROM golang:1.19.6-alpine3.17 as builder
WORKDIR /app
COPY . .
RUN \
    go mod download && \
    CGO_ENABLED=0 GOOS=linux go build ./cmd/server 

FROM alpine:3.17 as dev
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["/app/server"]