FROM golang:1.23-alpine3.21

EXPOSE 8085

RUN apk update && apk add --no-cache mysql-client build-base

RUN mkdir /app
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o /usr/local/bin/main cmd/main.go

CMD ["main"]
