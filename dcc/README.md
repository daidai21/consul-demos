# Distributed Configuration Center

> with consul
> client for pull config for update

## run demo

```shell
#拉取镜像
docker pull consul

#run
docker run -d -p 8500:8500 --restart=always --name=consul consul:latest agent -server -bootstrap -ui -node=1 -client='0.0.0.0'

#浏览器打开
open http://localhost:8500/ui/dc1/services

#`Key/Value`下面新建文件夹`DCC`，配置的key都在这里

#设置`DDC/test`的value为`123`

#using look unittest
```

## UT

`Test_Get_UpdateDelay` run log:

```shell
=== RUN   Test_Get_UpdateDelay
getFromLocalCache  key=DCC/test
getFromConsul  key=DCC/test
{}
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
{}
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
{}
getFromLocalCache  key=DCC/test
watchAndUpdate  更新KV  pair=&{Key:DCC/test Value:[123 34 103 101 34 58 34 115 100 34 125]}
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
{"ge":"sd"}
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
{"ge":"sd"}
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
{"ge":"sd"}
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
getFromConsul  key=DCC/test
{"ge":"sd"}
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
{"ge":"sd"}
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
{"ge":"sd"}
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
getFromConsul  key=DCC/test
{"ge":"sd"}
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
getFromLocalCache  key=DCC/test
--- PASS: Test_Get_UpdateDelay (10.03s)
PASS

Debugger finished with the exit code 0
```

## lines

```shell
➜  consul-demos git:(master) ✗ find dcc | grep "\.go" | grep -v "test" | xargs wc -l
      43 dcc/client/options.go
     165 dcc/client/client.go
      13 dcc/client/tool.go
      89 dcc/client/watcher.go
     310 total

```
