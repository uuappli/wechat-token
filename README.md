# Wechat Token
微信Access Token中控服务器，用来统一管理各个公众号的access_token，提供统一的接口进行获取和自动刷新Access Token。


## 项目特点

* REST风格的Web服务方式提供一致的获取接口
* 使用Echo框架编写的REST API，性能优秀
* 使用BuntDB作为access_token的内存缓存数据库，并且同时支持数据持久化
* 支持Basic Auth的认证方式，需要通过HTTP Basic认证才能访问，增强了安全性


## 安装

```bash
git clone https://github.com/gnuos/wechat-token.git
cd wechat-token
go get -v .
go build
```


## 快速开始

1. 在安装之后，会生成一个 wechat-token 程序，修改account.json文件，把其中的appid和secret替换成你自己的值；

2. 在项目目录中运行 ./wechat-token 就启动了一个Echo服务。如果操作系统开启了防火墙，需要防火墙开放8080端口的访问；

3. 如果有多个微信公众号的access_token需要管理，只需要在 account.json 文件中按格式把你的AppID和AppSecret添加到数组中就可以了。

4. 默认情况下，没有给框架配置日志输出，如果需要定制日志输出方式，请参考Echo框架的文档。


## License

Apache License, Version 2.0

