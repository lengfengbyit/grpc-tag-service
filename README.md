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

