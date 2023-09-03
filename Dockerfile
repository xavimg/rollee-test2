FROM golang:1.20.7-alpine3.17

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CONFIG_FILE=config.yml

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd

EXPOSE 3001

CMD ["./app"]
