# syntax=docker/dockerfile:1

FROM golang:1.17

WORKDIR ./meli

copy go.mod .

COPY . .

RUN go mod tidy -v

run go build -o /meli

EXPOSE 8080

CMD ["/meli"]