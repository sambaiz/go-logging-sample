FROM golang:1.10 AS builder
ADD . /go/src/github.com/sambaiz/go-logging-sample/
WORKDIR /go/src/github.com/sambaiz/go-logging-sample
RUN go get -u github.com/golang/dep/cmd/dep \
    && dep ensure \
    && go build -o go-logging-sample main.go

FROM alpine
WORKDIR /root/
COPY --from=builder /go/src/github.com/sambaiz/go-logging-sample/go-logging-sample .
CMD ["./go-logging-sample"]
