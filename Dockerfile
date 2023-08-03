FROM golang:1.20-alpine

RUN mkdir -p /app

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build

CMD ["./go-eats-server"]

EXPOSE 8080