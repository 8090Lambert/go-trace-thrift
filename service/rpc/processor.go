package rpc

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/8090Lambert/tracer_thrift/service"
	"github.com/huandu/go-tls"
)

var (
	SafeStorage *service.CoSafeStorage
)

type Processor struct {
	processor thrift.TProcessor
}

func NewProcessorDefault(processor thrift.TProcessor) *Processor {
	return &Processor{
		processor: processor,
	}
}

func (p *Processor) Process(inProtocol, outProtocol thrift.TProtocol) (bool, thrift.TException) {
	if custom.Reader.Column != nil {
		// Read special column.
		readCallback := custom.Reader.CallBack
		if readCallback == nil {
			readCallback = custom.ReadCallBackStub()
		}
		res, transport, err := readCallback(inProtocol.(*Protocol), custom.Reader.Column)
		inProtocol = NewProtocolDefault(transport)
		if err != nil {
			if exception, ok := err.(thrift.TTransportException); ok && exception.TypeId() == thrift.END_OF_FILE {
				return true, nil
			}
		}
		SafeStorage.Set(tls.ID(), res)
		defer SafeStorage.Delete(tls.ID())
	}

	success, err := p.processor.Process(inProtocol, outProtocol)

	return success, err
}
