# starGo
starGo是一款高性能、分布式、轻量级、微服务的游戏服务器框架。框架采用了go语言开发，得益于go本身对高并发的强大支持，框架足够简洁，效率足够高效。服务器框架二次开发简单易上手，实现了高性能的异步网络库，分布式节点间的通信采用了高性能通信中间件nats，能够实现每秒20万的Qps,日志管理，常规的关系型数据库（mysql）和非关系型数据库redis的支持，goroutine的安全定时器工具等。

服务器框架可用于包括但不限于游戏服务器等的应用，可以在框架开发阶段上节省大量时间。同时可以通过nats的发布订阅通道进行服务器的热更新等操作.
#### 优势特点
    1) 开发效率高
    2) 支持自定义的通信协议
    3) 采用nats高性能通信中间，实现了异步发布订阅，请求响应等模式的通信请求
    4) 分布式、微服务的架构，方便横向拓展
    5) 协程安全的定时器实现
    6) 对协程安全封装，优雅的实现开启退出
    7) 内置redis、mysql数据库支持
#### 安装教程
由于使用了nats做为通信中间件，所以需要安装nats服务器。使用docker安装非常的方便，简洁.

    docker pull nats:latest
    docker run -p 4222:4222 -p 8222:8222 -p 6222:6222 -ti nats:latest
    
然后安装依赖项，golang 1.13版本提供非常方便的依赖项管理工具go mod，你只需要输入go mod tidy，便可非常方便快捷的倒入依赖项。

当然你也可通过go get来手动导入依赖包

    go get github.com/go-redis/redis
    go get github.com/jinzhu/gorm
    go get github.com/nats-io/nats.go
    go get github.com/satori/go.uuid
    go get github.com/zhnxin/csvreader
#### 使用说明
    
#### 参与贡献

1. Fork 本仓库
2. 新建 Feat_xxx 分支
3. 提交代码
4. 新建 Pull Request
