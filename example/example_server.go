package example

import (
	"github.com/8090Lambert/tracer_thrift/example/api/echo"
	"github.com/8090Lambert/tracer_thrift/example/api/trace"
	"github.com/8090Lambert/tracer_thrift/service/rpc"
)

func NewServer() {
	rpc.NewServer(echo.NewEchoProcessor(
		EchoServerInterface{}),
		new(trace.Trace),
	).Run()
}
