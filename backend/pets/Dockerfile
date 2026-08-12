FROM golang:1.26.2 AS builder

WORKDIR /app

COPY backend/go.work ./backend/
COPY backend/pets/go.mod backend/pets/go.sum ./backend/pets/
COPY backend/pkg/go.mod backend/pkg/go.sum ./backend/pkg/

COPY . .

WORKDIR /app/backend/pets
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/pets ./cmd

FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata curl

RUN curl -k \
    https://gu-st.ru/content/lending/russian_trusted_root_ca_pem.crt \
    -o /usr/local/share/ca-certificates/russian_trusted_root_ca.crt \
    && update-ca-certificates

COPY --from=builder /bin/pets /app/pets

CMD ["./pets"]