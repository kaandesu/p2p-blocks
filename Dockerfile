FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o p2p-blocks .

EXPOSE 3000 

CMD ["./p2p-blocks", "--apex"]
