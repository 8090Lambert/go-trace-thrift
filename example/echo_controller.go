package example

import (
	"github.com/8090Lambert/tracer_thrift/example/api/trace"
	"github.com/8090Lambert/tracer_thrift/service/rpc"
	"github.com/huandu/go-tls"
	"log"
)

type EchoServerInterface struct {
}

func (e EchoServerInterface) Emit(message string) (string, error) {
	tracer, _ := rpc.SafeStorage.Get(tls.ID())
	if _, ok := tracer.(*trace.Trace); !ok {
		log.Fatal("parse failed")
		res := message + "_service"
		return res, nil
	}

	traceConcrete := tracer.(*trace.Trace)
	res := "success, trace=" + traceConcrete.Id + ", pid=" + traceConcrete.Pid + ", cid=" + traceConcrete.Cid

	return res, nil
}
