FROM golang:1.26.2-alpine AS builder

WORKDIR /app

COPY backend/go.work ./backend/
COPY backend/tasks/go.mod backend/tasks/go.sum ./backend/tasks/
COPY backend/pkg/go.mod backend/pkg/go.sum ./backend/pkg/

COPY . .
WORKDIR /app/backend/tasks
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/tasks ./cmd/app

FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache ca-certificates

COPY --from=builder /bin/tasks /app/tasks

CMD ["./tasks"]
