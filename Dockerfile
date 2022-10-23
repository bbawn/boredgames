# syntax=docker/dockerfile:1

FROM golang

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /boredgames cmd/api

EXPOSE 8080

CMD [ "/boredgames" ]
