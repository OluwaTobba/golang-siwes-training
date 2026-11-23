FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY app ./app
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server ./app

FROM scratch
COPY --from=builder /app/server /server
USER 1001:1001
EXPOSE 8080
ENTRYPOINT ["/server"]