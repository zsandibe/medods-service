FROM golang:1.22.6

WORKDIR /api-server

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

EXPOSE 7777