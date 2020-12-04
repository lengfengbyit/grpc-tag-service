# grpc-tag-service
> 使用grpc调用go-blog的tag服务

[grpc中文文档](http://doc.oschina.net/grpc)

### 调试grpc接口
> grpc 是基于HTTP/2协议的，不能像HTTP/1.1接口那样可以直接通过postman或者普通的curl进行调试， 这里可以使用grpcurl来请求grpc接口

1. 安装 grpcurl 
```shell script
go get github.com/fullstorydev/grpcurl
go install github.com/fullstorydev/grpcurl/cmd/grpcurl
```

2. 注册服务
> 使用 grpcurl 的前提是已经注册了反射服务
```go
import "google.golang.org/grpc/reflection"

func main() {
    s := grpc.NewServer()

    // 注册服务，方便使用grpcurl调试
    reflection.Register(s)
}
```
-  查看服务列表
```shell script
grpcurl -plaintext localhost:8080 list
```
-  查看服务提供的接口 
```shell script
grpcurl -plaintext localhost:8080 list proto.TagServer
```
1. plaintext: grpcurl工具默认使用TLS认证(可通过-cert和-key参数设置公钥和秘钥)， 这个参数是用来忽略TLS认证的
2. localhost:8080 指定运行的服务的地址和端口
3. list: 不跟参数列出注册的服务列表， 跟上服务参数，可以查看服务所提供的接口列表

- 调用接口
```shell script
grpcurl -plaintext -d '{"name":"Go"}' 127.0.0.1:8080 proto.TagService.GetTagList
```

### grpc状态码
![grpc状态码1](https://gitee.com/fym321/picgo/raw/master/imgs/20201201221159.png)
![grpc状态码2](https://gitee.com/fym321/picgo/raw/master/imgs/20201201221306.png)

### 在同端口监听HTTP和GRPC
> 通过开源库实现，cmux根据有效负载(payload)对连接进行多路复用(只匹配连接的前面的字节来区分当前连接的类型), 可以在同一 TCP Listener上提供gRPC,SSH,HTTPS,HTTP和Go RPC等几乎所有其他协议的服务，是一个相对通用的方案

- 下载 cmux
```shell script
go get -u github.com/soheilhy/cmux
```

### 同端口同方法提供双流量支持 grpc-gateway
> grpc-gateway是protoc的一个插件，它能够读取Protobuf的服务定义，生成一个反向代理服务器，将RESTful Json API转换为gRPC。 它主要是根据Protobuf的服务定义中的google.api.http来生成的。

![grpc-gateway架构图](https://gitee.com/fym321/picgo/raw/master/imgs/20201202143717.png)

> 简单的来说， grpc-gateway能够将RESTful转换为gRPC请求，实现用同一个RPC方法提供gRPC协议和HTTP/1.1的双流量支持。

- 安装grpc-gateay
```shell script
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@v1.16.0
```

- 重新编译proto文件
```shell script
// -I 参数执行 proto 文件中 import 的搜索目录， 不指定则为当前目录
// 该命令会再proto目录下生产一个 .pb.gw.go的文件， 如果配置好了的话
protoc -I/Users/fym/Documents/code/go/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis -I. -I$GOPATH/src --grpc-gateway_out=logtostderr=true:. ./proto/*proto
```

### 接口文档
> 使用 protoc-gen-swagger 来根据 protoc 文件自动生成 swagger 定义

- 安装插件
```shell script
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
```

- 使用 bindata 将资源转为 go 文件
```shell script
# 下载
go get -u github.com/go-bindata/go-bindata/...
# 将资源转为go文件
go-bindata --nocompress -pkg swagger -o pkg/swagger/data.go third_party/swagger-ui/...
# 让静态资源代码能够被外部访问，该库结合net/http标准库和go-bindata库生成的
# swagger ui的go代码供外部访问
go get -u github.com/elazarl/go-bindata-assetfs/...

# 生成swagger.json文件
protoc -I$GOPATH/src -I. -I$GOPATH/src \
-I$GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis \
--swagger_out=logtostderr=true:. ./proto/*.proto

```

### gRPC拦截器
- 一元拦截器(Unary Interceptor): 拦截和处理一元RPC调用
- 流拦截器(Stream Interceptor): 拦截和处理流式RPC调用

> 由于客户端和服务端有各自的一元拦截器和流拦截器，因此，在gRPC中， 也可以说共偶四种类型的拦截器

***1. 客户端***
- 一元拦截器， 类型为 UnaryClientInterceptor, 原型如下
```go
type UnaryClientInterceptor func(
    ctx context.Context,            // RPC上下文
    method string,                  // 所调用的方法名称
    req, reply interface{},         // RPC方法的请求参数和响应结果
    cc *ClientConn,                 // 客户端连接句柄
    invoker UnaryInvoker,           // 所调用的RPC方法
    opts ...CallOption              // 调用的配置
) error
```

> 一元拦截器的实现通常分为三部分： 预处理、调用RPC方法和后处理

- 流拦截器， 类型为 StreamClientInterceptor, 原型如下
```go
type StreamClientInterceptor func(
    ctx context.Context,            // RPC上下文
    desc *StreamDesc,               // 流描述
    cc *ClientConn,                 // 客户端连接句柄
    method string,                  // 所调用的方法名称
    streamer Streamer,              // 所调用RPC方法的流对象
    opts ...CallOption              // 所调用RPC方法的配置
) (ClientStream, error)
```

> 流拦截器的实现包括预处理和流操作拦截两部分，不能再时候进行RPC方法调用和后处理，只能拦截用户对流的操作

***2. 服务端***
- 一元拦截器， 类型为 UnaryServerInterceptor， 原型如下
```go
type UnaryServerInterceptor func(
    ctx context.Context,            // RPC上下文
    req interface{},                // RPC方法的请求参数
    info *UnaryServerInfo,          // RPC方法的所有信息
    handler UnaryHandler            // RPC方法本身
) (resp interface{}, err error)
```

- 流拦截器，类型为 StreamServerInterceptor 原型如下
```go
type StreamServerInterceptor func(
    srv interface{},                
    ss ServerStream, 
    info *StreamServerInfo, 
    handler StreamHandler
) error
```

***3.使用多个拦截器***
> grpc-go官方只提供一个拦截器的钩子，以便开发人员在其上构建各种复杂的拦截器模式，而不会遇到多个拦截器的执行顺序问题，同时还能保持grpc-go自身的间接性，尽可能的最小化grpc-go的公共API

> 可以使用go-grpc-middleware提供的grpc.UnaryInterceptor和grpc.StreamInterceptro，以链式方式达到多个拦截器的目的

- 1. 安装 go-grpc-middleware
```shell script
go get -u github.com/grpc-ecosystem/go-grpc-middleware@v1.2.2
```

- 2. 使用
```go
opts := []grpc.ServerOption{
    // 使用 grpc_middleware.ChainUnaryServer 链式调用多个拦截器
    grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
        HelloInterceptor,
        WorldInterceptor,
    )),
}
s := grpc.NewServer(opts...)
```

### metadata
> 在http/1.1中， 我们常常通过Header来传递数据，而在gRPC中，因为他是基于HTTP/2的，所以本质上也可以通过Header来传递数据，但不是直接去操纵它，
> 而是通过gRPC中的metadata来进行传递和操作的， 使用metadata需要我们使用的库支持。

1. metadata 数据结构
```go
type MD map[string][]string
```

2. 创建 metadata
- google.golang.org/grpc/metadata 提供两个方法来创建metadata,第一种是metadata.New方法，如下
```go
metadata.New(map[string]string{"go":"Hello", "php": "World"})
```
使用New方法创建的metadata会直接被转换MD结构，如下
```go
map[
    go:[]string{"Hello"}
    php:[]string{"world"}
]
```

- 第二种是 metadata.Pairs 方法， 如下
```go
metadata.Pairs(
    "go", "Hello",
    "php", "World",
)
```

> 通过抓包分析， 实际上 metadata 也是通过 header 来传递数据的。

### 对RPC方法做自定义认证

如果需要对某些模块的RPC方法做特殊认证或校验，可以使用 grpc 提供的 Token 接口，代码如下：

```go
type PerRPCCredentials interface {
    // 获取当前请求认证所需的元数据
	GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error)
    // 是否需要基于TLS认证进行安全传输
	RequireTransportSecurity() bool
}
```

### 链路追踪 grpc + jaeger
1. 注入追踪信息
> 做链路追踪的基本条件就是要注入追踪信息，最简单的方式莫过于使用拦截器实现

- 服务端拦截器： 从metadata中提取链路信息，并将其设置追踪到服务端的上下文中。
- 客户端拦截器: 从调用的上下文中提取链路信息，并作为metadata追缴到RPC的调用中

2. 下载 jaeger
```shell script
go get -u github.com/opentracing/opentracing-go@v1.2.0
go get -u github.com/uber/jaeger-client-go@v1.6.0
```

3. metadata的读取和配置
在OPenTracing中，可以指定SpanContexts的三种传输表达方式，代码如下:
```go
type BuiltinFormat byte
const (
    // 不透明的二进制数据
    Binary BuiltinFormat = iota
    // 键值字符串对
    TextMap
    // HTTP header 格式的字符串
    HTTPHeaders
)

type TextMapWriter interface {
    Set(key, val string)
}

type TextMapReader interface {
    ForeachKey(handler func(key, val string) error) error
}
```

### 服务注册与发现

1. 所涉及到的角色

- 注册中心：承担对服务信息进行注册、协调、管理等工作
- 服务提供者：暴露特定端口，并提供一个到多个服务允许外部访问
- 服务消费者：服务提供者的调用方

2. 负载均衡策略
- 客户端负载： 客户端到注册中心查询服务提供者清单，并使用特定的负载均衡策略在服务清单中选择一个或多个服务进行调用。
    - 优点： 去中心化、不需要借助独立的外部负载均衡组件
    - 缺点: 实现成本较高，不同语言的客户端都需要实现对应负载均衡策略

![客户端负载架构图](https://gitee.com/fym321/picgo/raw/master/imgs/20201204075110.png)

- 服务端负载：有被称为代理模式，在服务端搭建独立的负载均衡服务器，客户端统一发送请求到代理服务器，由代理服务器实现负载均衡策略并反向代理到某个具体的业务服务器。
    - 优点: 简单、透明，客户端不需要知道背后的逻辑，直接调用即可
    - 缺点: 可能会成为性能瓶颈，与客户端负载想比，可能会出现更高的网络延迟

![服务端负载架构图](https://gitee.com/fym321/picgo/raw/master/imgs/20201204075029.png)

- gRPC官方设计思路

> 官方并没有直接该出gRPC服务发现和负载均衡相关的功能代码，而是在官方文档load-balancing中进行了详细的介绍，说明了实现思路，并在gRPC API中提供了各类应用的接口，以便外部扩展

[gRPC load-balancing 官方文档](https://github.com/grpc/grpc/blob/master/doc/load-balancing.md)

![gRPC负载架构](https://gitee.com/fym321/picgo/raw/master/imgs/20201204075543.png)

- 在进行gRPC调用时，gRPC客户端会想名称解析器(Name Resolver)发出服务端名称的解析请求。名称解析器会将服务名解析成一个或多个IP地址，每个IP地址都会标识它是服务端地址，还是负载均衡地址，以及客户端使用的负载均衡策略(如round_robin或grpclb)

- 客户端实例化负载均衡策略，如果gRPC客户端获取的地址是负载均衡地址，那么客户端将使用grpclb策略，否则使用服务配置请求的负载均衡策略。如果服务端未配置负载均衡策略，则客户端端默认选择第一个可用的服务端地址。

- 负载均衡策略会为每一个服务端地址创建一个子通道

- 对于每一个请求，都有负载均衡策略决定将其发送到那个子通道(即grpc服务端)

3. 服务注册于发现















