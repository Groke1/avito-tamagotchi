FROM golang:1.26 AS builder

WORKDIR /app

COPY backend/pets/go.mod backend/pets/go.sum ./
COPY backend/pkg/go.mod backend/pkg/go.sum ./backend/pkg/

RUN go mod download

COPY backend/pets/ .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/app ./cmd

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /bin/app .

CMD ["./app"]