FROM golang:1.24-bookworm AS builder
WORKDIR /app
ENV GOPROXY=https://goproxy.cn,direct

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \ 
    go build -o tool-server cmd/tool-server/main.go

FROM debian:bookworm-slim AS runtime

WORKDIR /app
COPY --from=builder /app/tool-server /app/

EXPOSE 8080

ENTRYPOINT ["/app/tool-server"]
