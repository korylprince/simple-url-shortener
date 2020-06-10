FROM golang:1.14-alpine as builder

ARG VERSION

RUN apk add --no-cache git ca-certificates build-base

RUN git clone --branch "v1.1" --single-branch --depth 1 \
    https://github.com/korylprince/fileenv.git /go/src/github.com/korylprince/fileenv

RUN git clone --branch "$VERSION" --single-branch --depth 1 \
    https://github.com/korylprince/simple-url-shortener.git  /go/src/github.com/korylprince/simple-url-shortener

RUN go install github.com/korylprince/fileenv
RUN go install github.com/korylprince/simple-url-shortener

FROM alpine:3.12

RUN apk add --no-cache ca-certificates

COPY --from=builder /go/bin/fileenv /
COPY --from=builder /go/bin/simple-url-shortener /
COPY setenv.sh /

CMD ["/fileenv", "sh", "/setenv.sh", "/simple-url-shortener"]
