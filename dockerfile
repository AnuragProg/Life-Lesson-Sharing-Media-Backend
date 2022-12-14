FROM golang:latest

COPY . /app

WORKDIR /app

RUN go build -o app

CMD ["/app/app"]