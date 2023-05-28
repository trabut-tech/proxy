FROM golang:1.20-alpine

WORKDIR /app/

RUN apk update && apk add git

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o /usr/bin/app .

WORKDIR /app/

ENV LISTEN_ADDR=":33080"

CMD app --listen=$LISTEN_ADDR --verbose