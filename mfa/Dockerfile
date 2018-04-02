FROM golang:1.9.2 as builder
RUN go get -d -v \
  github.com/gorilla/mux \
  github.com/lib/pq \
  github.com/bwmarrin/snowflake \
  gopkg.in/guregu/null.v3 \
  github.com/kevinmichaelchen/my-go-utils
ADD . /go/src/github.com/hanakoa/alpaca/mfa
WORKDIR /go/src/github.com/hanakoa/alpaca/mfa
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

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
COPY --from=builder /go/src/github.com/hanakoa/alpaca/mfa/app .
#USER $ALPACA_USER
CMD ["./app"]