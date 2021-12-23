# go-rpc

go-rpc 是在阅读 [gee-orm](https://github.com/geektutu/7days-golang/tree/master/gee-orm) 源码时的一些注解。

orm 框架由如下几部分构成，分别为 :
- 消息编解码
- 客户端
- 服务端
- 服务注册
- 服务发现
- 超时控制
- 负载均衡

接下来通过图片的方式， 给出各个组件的模块和构成。

# 整体流程


