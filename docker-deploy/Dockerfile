FROM golang:1.20.5

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN go build -o main .

CMD ["./main"]