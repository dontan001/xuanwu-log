# Use busybox as the base image
FROM golang:1.17-alpine

WORKDIR /app
COPY bin/xuanwu-api xuanwu-api

ENTRYPOINT [ "/app/xuanwu-api" ]