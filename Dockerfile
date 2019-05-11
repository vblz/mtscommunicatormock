FROM golang:1.12-alpine as build

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN apk add --no-cache --update git && \
    rm -rf /var/cache/apk/*

ADD . /build
WORKDIR /build

RUN go build -mod=vendor -o out/mtsMock .

FROM alpine:3.9

ENV HTTP_PORT=9000
ENV UTC_OFFSET=-4

COPY --from=build /build/out/mtsMock /srv
COPY --from=build /build/static /srv/static

WORKDIR /srv

ENTRYPOINT ["/srv/mtsMock"]