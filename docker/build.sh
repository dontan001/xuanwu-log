#!/bin/bash

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o bin/xuanwu-log main.go

docker build -t dontan001/xuanwu-log:latest .