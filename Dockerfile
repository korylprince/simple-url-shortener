FROM golang:1-alpine as builder

ARG VERSION

RUN go install github.com/korylprince/fileenv@v1.1.0
RUN go install "github.com/korylprince/simple-url-shortener@$VERSION"


FROM alpine:3.15

RUN apk add --no-cache ca-certificates

COPY --from=builder /go/bin/fileenv /
COPY --from=builder /go/bin/simple-url-shortener /
COPY setenv.sh /

CMD ["/fileenv", "sh", "/setenv.sh", "/simple-url-shortener"]
