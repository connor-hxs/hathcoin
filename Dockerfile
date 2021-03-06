FROM golang:1.9-alpine
MAINTAINER i@boris.tech

COPY . /go/src/github.com/borisding1994/hathcoin
WORKDIR /go/src/github.com/borisding1994/hathcoin
RUN apk --no-cache add make openssl git

# Add USTC repo
RUN apk add --update-cache tzdata \
    && setup-timezone -z Asia/Shanghai \
    && apk del tzdata \
# Sepcify ntp server
    && sed -i 's/pool.ntp.org/cn.pool.ntp.org/' /etc/init.d/ntpd

# build binary
RUN make dep-install && \
    make build && \
    mkdir -p /app

# remove sources
RUN cp -r /go/src/github.com/borisding1994/hathcoin/dist /app && \
    rm -rf /go/src/*

WORKDIR /app
#VOLUME /app/logs
EXPOSE 8081 8081
ENTRYPOINT ["/app/hathcoin", "server"]