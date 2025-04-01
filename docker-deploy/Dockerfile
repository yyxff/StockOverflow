FROM golang:1.24.1

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy
RUN go mod download

COPY . .

RUN go build -o main ./cmd/exchange/main.go

CMD ["./main"]