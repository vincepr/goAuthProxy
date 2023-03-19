# step 1 build the executable binary
FROM golang:alpine As builder

RUN apk update && apk add --no-cache git

WORKDIR $GOPATH/src/go_auth_proxy
COPY src/. .

RUN go get -d -v
RUN go build -o /go/bin/goAuthProxy

# step 2 build a minimal image from scratch (just binary + html/js)
FROM scratch
COPY --from=builder /go/bin/goAuthProxy /go/bin/goAuthProxy
COPY src/public/. /var/www/goAuthProxy/

EXPOSE 5555

ENTRYPOINT [ "/go/bin/goAuthProxy", "--port", "5555" ]