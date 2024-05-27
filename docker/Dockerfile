FROM alpine:latest

ENV TZ=Asia/Shanghai

COPY docker/install.sh /tmp

RUN apk add --no-cache \
    binutils wget curl && \
    ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && \
    echo $TZ > /etc/timezone && \
    chmod +x /tmp/install.sh && \
    /tmp/install.sh

ENTRYPOINT ["/opt/miaospeed"]