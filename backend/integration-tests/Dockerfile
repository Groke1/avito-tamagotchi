FROM golang:1.26.2

WORKDIR /tests

COPY backend/integration-tests/go.mod backend/integration-tests/go.sum ./
RUN go mod download

COPY backend/integration-tests/. .
RUN go mod tidy

CMD ["go", "test", "-v", "./..."]
