FROM golang:1.26.2 AS builder

WORKDIR /app

COPY backend/go.work ./backend/
COPY backend/user/go.mod backend/user/go.sum ./backend/user/
COPY backend/pkg/go.mod backend/pkg/go.sum ./backend/pkg/

COPY . .

WORKDIR /app/backend/user
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/user ./cmd/app

FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /bin/user /app/user

CMD ["./user"]
