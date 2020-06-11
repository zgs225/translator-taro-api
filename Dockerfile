FROM golang:1.14

WORKDIR /go/src/app
COPY . .

RUN go install -v ./...

CMD ["translator-api"]