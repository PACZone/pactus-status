FROM golang:1.21.1

WORKDIR /app

COPY . .

RUN go build -o pactus-status .

EXPOSE 9854

CMD ["./pactus-status"]
