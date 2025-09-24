FROM golang:1.25.1

WORKDIR /subscription

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o subscription ./cmd/server

EXPOSE 8080

CMD ["./subscription"]