FROM golang:1.19

WORKDIR /p2p-orders-parser

COPY . ./

RUN go mod download

RUN go build -o parser .

EXPOSE 8080

CMD ["./parser"]