FROM golang:1.26.2-alpine AS builder

WORKDIR /src/backend/tasks

RUN apk add --no-cache ca-certificates

COPY backend/tasks/go.mod backend/tasks/go.sum ./
RUN go mod download

COPY backend/tasks/ .
RUN find /src/backend/tasks -maxdepth 4 -type f -print
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/tasks ./cmd/app

FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /out/tasks /app/tasks

EXPOSE 8080

ENV HTTP_ADDR=:8080

CMD ["/app/tasks"]