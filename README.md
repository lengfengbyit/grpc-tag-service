# grpc-tag-service
> 使用grpc调用go-blog的tag服务

### 调试grpc接口
> grpc 是基于HTTP/2协议的，不能像HTTP/1.1接口那样可以直接通过postman或者普通的curl进行调试， 这里可以使用grpcurl来请求grpc接口

1. 安装 grpcurl 
```bash
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
```bash
grpcurl -plaintext localhost:8080 list
```
-  查看服务提供的接口 
```bash
grpcurl -plaintext localhost:8080 list proto.TagServer
```
1. plaintext: grpcurl工具默认使用TLS认证(可通过-cert和-key参数设置公钥和秘钥)， 这个参数是用来忽略TLS认证的
2. localhost:8080 指定运行的服务的地址和端口
3. list: 不跟参数列出注册的服务列表， 跟上服务参数，可以查看服务所提供的接口列表

- 调用接口
```bash
grpcurl -plaintext -d '{"name":"Go"}' 127.0.0.1:8080 proto.TagService.GetTagList
```

3. grpc状态码
![grpc状态码1](https://gitee.com/fym321/picgo/raw/master/imgs/20201201221159.png)
![grpc状态码2](https://gitee.com/fym321/picgo/raw/master/imgs/20201201221306.png)

4. 在同端口监听HTTP和GRPC
> 通过开源库实现，cmux根据有效负载(payload)对连接进行多路复用(只匹配连接的前面的字节来区分当前连接的类型), 可以在同一 TCP Listener上提供gRPC,SSH,HTTPS,HTTP和Go RPC等几乎所有其他协议的服务，是一个相对通用的方案

- 下载 cmux
```bash
go get -u github.com/soheilhy/cmux
```

5. 同端口同方法提供双流量支持 grpc-gateway
> grpc-gateway是protoc的一个插件，它能够读取Protobuf的服务定义，生成一个反向代理服务器，将RESTful Json API转换为gRPC。 它主要是根据Protobuf的服务定义中的google.api.http来生成的。

![grpc-gateway架构图](https://gitee.com/fym321/picgo/raw/master/imgs/20201202143717.png)

> 简单的来说， grpc-gateway能够将RESTful转换为gRPC请求，实现用同一个RPC方法提供gRPC协议和HTTP/1.1的双流量支持。

- 安装grpc-gateay
```bash
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@v1.16.0
```

- 重新编译proto文件
```bash
// -I 参数执行 proto 文件中 import 的搜索目录， 不指定则为当前目录
// 该命令会再proto目录下生产一个 .pb.gw.go的文件， 如果配置好了的话
protoc -I/Users/fym/Documents/code/go/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis -I. -I$GOPATH/src --grpc-gateway_out=logtostderr=true:. ./proto/*proto
```

6. 接口文档
> 使用 protoc-gen-swagger 来根据 protoc 文件自动生成 swagger 定义

- 安装插件
```bash
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
```

- 使用 bindata 将资源转为 go 文件
```bash
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