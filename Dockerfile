FROM golang:1.13.8-alpine3.11 as builder
COPY . /go/src/github.com/kmaris/docker-volume-iso
WORKDIR /go/src/github.com/kmaris/docker-volume-iso
RUN set -ex \
  && apk add --no-cache --virtual .builder gcc libc-dev \
  && go install --ldflags '-extldflags "-static"' \
  && apk del .builder

FROM alpine
RUN mkdir -p /run/docker/plugins /mnt/volumes
COPY --from=builder /go/bin/docker-volume-iso /bin
CMD ["/bin/docker-volume-iso"]
