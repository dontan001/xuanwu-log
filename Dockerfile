# Use busybox as the base image
FROM golang:1.17-alpine

WORKDIR /app
COPY bin/xuanwu-log xuanwu-log

CMD [ "/app/xuanwu-log" ]