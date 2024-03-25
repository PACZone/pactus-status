FROM golang:1.22.1-alpine3.18 as builder

WORKDIR /app

COPY . .

RUN go build -o pactus-status .

EXPOSE 9854

CMD ["./pactus-status"]
