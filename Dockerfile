FROM golang:1.16.2-alpine AS build-env

WORKDIR /build

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux

RUN apk --no-cache add git=~2

COPY *.go go.mod go.sum /build/

RUN go version
RUN go build

FROM scratch

ENV PATH=/

COPY --from=build-env /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
COPY --from=build-env /etc/ssl/certs/ /etc/ssl/certs/
COPY --from=build-env /build/reviewbot /reviewbot

ENTRYPOINT ["reviewbot"]
