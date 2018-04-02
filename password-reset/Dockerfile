FROM golang:1.10 as builder
ADD . /go/src/github.com/hanakoa/alpaca
WORKDIR /go/src/github.com/hanakoa/alpaca/password-reset
RUN go get ./... && \
    CGO_ENABLED=0 GOOS=linux go build -a -o ./bin/alpaca-password-reset .

FROM alpine:latest
#ENV ALPACA_USER="alpaca" \
#    ALPACA_UID="8080" \
#    ALPACA_GROUP="alpaca" \
#    ALPACA_GID="8080"
RUN apk --no-cache add ca-certificates
#&& \
#    addgroup -S -g $ALPACA_GID $ALPACA_GROUP && \
#    adduser -S -u $ALPACA_UID -G $ALPACA_GROUP $ALPACA_USER
WORKDIR /root/
COPY --from=builder /go/src/github.com/hanakoa/alpaca/password-reset/bin/alpaca-password-reset .
#USER $ALPACA_USER
CMD ["./alpaca-password-reset"]