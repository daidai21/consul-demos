# Distributed Configuration Center

> with consul

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
```


