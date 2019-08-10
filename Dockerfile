FROM golang:1.12 as builder

WORKDIR /go/src/conchatenate

COPY . .

RUN cd server

RUN go get -v -d ./...

RUN go install -v ./...

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/server ./server

RUN useradd -M appuser

FROM scratch as deploy

COPY --from=builder /go/bin/server /go/bin/server

COPY --from=builder /go/src/conchatenate/static/testpage.html ./static/testpage.html

COPY --from=builder /etc/passwd /etc/passwd

EXPOSE 8080

USER appuser

CMD ["/go/bin/server"]
