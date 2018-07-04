FROM golang:1.10 as builder
ADD . /go/src/github.com/hanakoa/alpaca
WORKDIR /go/src/github.com/hanakoa/alpaca/services/password-reset
# won't need to go get vgo in golang 1.11
RUN go get -u -v golang.org/x/vgo && \
    CGO_ENABLED=0 GOOS=linux vgo build -a -o ./bin/alpaca-password-reset .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/hanakoa/alpaca/services/password-reset/bin/alpaca-password-reset .
CMD ["./alpaca-password-reset"]