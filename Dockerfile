FROM golang:1.23

WORKDIR ./data-enrich

copy go.mod .

COPY . .

RUN go mod tidy -v

run go build -o /data-enrich

EXPOSE 8080

CMD ["/data-enrich"]
