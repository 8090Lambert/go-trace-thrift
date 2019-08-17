# Customized Thrift

> 自定义 thrift 服务

## Feature
针对对请求链路有追踪需求的应用，不必显式定义在每个 IDL 中 trace 参数，或是兼容 thrift 的 Generation 自动生成。
思路是：依赖于请求头，启动服务时，提前注册 struct, 解析请求参数时反解，重置 stream, struct 只做单次请求生命周期
内的存储 storage (协程安全的)。详见 example/

## Install
```
$ go get github.com/8090Lambert/CustomizedThrift
```

## Usage
```
import (
	"github.com/8090Lambert/tracer_thrift/example/api/echo"
	"github.com/8090Lambert/tracer_thrift/example/api/trace"
	"github.com/8090Lambert/tracer_thrift/service/rpc"
)

func NewServer() {
	rpc.NewServer(echo.NewEchoProcessor(
		EchoServerInterface{}),     // process handler
		new(trace.Trace),           // 约定好的 trace struct
	).Run()
}
```



