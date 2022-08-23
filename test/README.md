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
curl 'http://localhost:8080/log?query=\{job="fluent-bit",app="yinglong"\}&start=now-6h&end=now' -o test/yinglong.zip

```
```browser
http://localhost:8080/log?query={job="fluent-bit",app="yinglong"}&start=now-6h&end=now
http://localhost:8080/log/v2?query={job="fluent-bit",app="yinglong"}&start=now-6h&end=now
http://localhost:8080/log/big?file=test1g.txt

```