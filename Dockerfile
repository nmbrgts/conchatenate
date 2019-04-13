FROM golang:1.12

WORKDIR /go/src/conchatenate

COPY . .

RUN cd server

RUN go get -v -d ./...

RUN go install -v ./...

EXPOSE 8080

CMD ["server"]
