# syntax = docker/dockerfile:1-experimental
FROM golang:1.18-alpine3.15 AS base

RUN apk --no-cache --update-cache --upgrade --latest add build-base git gcc bash

WORKDIR /go/src/github.com/kjuulh/scel/server

ADD go.mod go.mod
ADD go.sum go.sum

ENV GO111MODULE on
ENV CGO_ENABLED 1

RUN go mod download

ADD . .

ARG VERSION
ARG COMMIT
ARG BUILD_DATE

RUN --mount=type=cache,target=/root/.cache/go-build go build \
  -o /usr/bin/scel

FROM alpine:3.15.4

RUN addgroup -S scel; \
  adduser -S scel -G scel -D -u 10000 -h /home/scel -s /bin/nologin; \
  chown -R scel:scel /home/scel

COPY --from=base /usr/bin/scel /usr/bin/scel

# Exposing the ory home directory to simplify passing in Kratos configuration (e.g. if the file $HOME/.kratos.yaml
# exists, it will be automatically used as the configuration file).
VOLUME /home/scel

# Declare the standard ports used by Kratos (4433 for public service endpoint, 4434 for admin service endpoint)
EXPOSE 15000 15000

USER 10000

ENTRYPOINT ["scel"]
CMD ["serve"]
