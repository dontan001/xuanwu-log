# Testing Guide

## How to generate big files
```shell 
// mac os
dd if=/dev/zero of=test/test1g.txt  bs=1000000 count=1000
dd if=/dev/zero of=test/test10g.txt bs=1000000 count=10000
dd if=/dev/zero of=test/test50g.txt bs=1000000 count=50000
```

## How to make a http request
```shell
curl 'http://localhost:8080/log/download?query=\{job="fluent-bit",app="yinglong"\}&start=now-6h&end=now' -o test/yinglong.zip
// query encoded
curl 'http://localhost:8080/log/download?query=%7Bjob%3D%22fluent-bit%22%2Capp%3D%22yinglong%22%7D&start=now-3h&end=now' -o test/yinglong.zip

```
```browser
http://localhost:8080/log/big?file=test1g.txt
http://localhost:8080/log/download?query={job="fluent-bit",app="yinglong"}&start=now-6h&end=now

```

## How to install the app via helm
```shell
// install with override
helm install test helm/xuanwu-log -f test/xuanwu-log.override --namespace xuanwu-log \
 --create-namespace --description "install v1" --debug --dry-run

// uninstall
helm uninstall test -n xuanwu-log --debug --dry-run 
```