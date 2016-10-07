FROM gliderlabs/alpine
MAINTAINER Thibault Deutsch <thibault.deutsch@gmail.com>

EXPOSE 80

VOLUME /static
WORKDIR /app

RUN apk-install ca-certificates

COPY statico /app/
ENTRYPOINT ["./statico"]
