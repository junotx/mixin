# partially refer to https://github.com/jmb12686/docker-filebeat

FROM golang:1.17.10-alpine3.15 as builder

RUN apk update && apk add --no-cache git && \
    apk add --no-cache make && \
    apk add --no-cache bash && \
    apk add --no-cache gcc && \
    apk add --no-cache libc-dev && \
    apk add --no-cache binutils-gold

RUN mkdir -p $GOPATH/src/github.com/elastic/beats && \
    git clone --branch v6.7.0 --depth 1 https://github.com/elastic/beats.git $GOPATH/src/github.com/elastic/beats
WORKDIR $GOPATH/src/github.com/elastic/beats/filebeat
ENV GO111MODULE=off
RUN make

## Run stage - Install dependencies and copy filebeat from builder
FROM alpine:3.14

RUN mkdir -p /usr/share/filebeat
RUN apk update && \
    apk add --no-cache bash

COPY --from=builder /go/src/github.com/elastic/beats/filebeat/filebeat /usr/share/filebeat/filebeat
RUN ["chmod", "+x", "/usr/share/filebeat/filebeat"]
COPY docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
RUN ["chmod", "+x", "/usr/local/bin/docker-entrypoint.sh"]
ENTRYPOINT [ "/usr/local/bin/docker-entrypoint.sh" ]
WORKDIR /usr/share/filebeat/
CMD [ "/usr/share/filebeat/filebeat" , "--strict.perms=false"]