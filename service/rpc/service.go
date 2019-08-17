package rpc

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/8090Lambert/tracer_thrift/service"
	"log"
	"runtime/debug"
	"sync"
	"sync/atomic"
)

func init() {
	SafeStorage = service.NewCoSafeStorage()
}

type TServer struct {
	mu        sync.Mutex
	wg        sync.WaitGroup // 平滑重启
	processor thrift.TProcessor
	bindAddr  string
	timeout   int64
	quit      chan struct{}
	stopped   uint32

	processorFactory       thrift.TProcessorFactory
	serverTransport        *thrift.TServerSocket
	inputTransportFactory  thrift.TTransportFactory
	outputTransportFactory thrift.TTransportFactory
	inputProtocolFactory   thrift.TProtocolFactory
	outputProtocolFactory  thrift.TProtocolFactory
}

const (
	UNSTOP = 0
	STOP   = 1
)

func (s *TServer) Serve() error {
	err := s.Listen()
	if err != nil {
		return err
	}
	s.AcceptLoop()
	return nil
}

func (s *TServer) Listen() error {
	return s.serverTransport.Listen()
}

func (s *TServer) AcceptLoop() error {
	for {
		closed, err := s.innerAccept()
		if err != nil {
			return err
		}

		if closed != UNSTOP {
			return nil
		}
	}
}

func (s *TServer) innerAccept() (int, error) {
	client, err := s.serverTransport.Accept()
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.IsStop() {
		return STOP, nil
	}

	if err != nil {
		return UNSTOP, err
	}

	if client != nil {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			if err := s.processRequests(client); err != nil {
				log.Printf("request error : %s", err)
			}
		}()
	}

	return UNSTOP, nil
}

func (s *TServer) processRequests(client thrift.TTransport) error {
	processor := s.processorFactory.GetProcessor(client)
	inputTransport := s.inputTransportFactory.GetTransport(client)
	outputTransport := s.outputTransportFactory.GetTransport(client)
	inputProtocol := s.inputProtocolFactory.GetProtocol(inputTransport)
	outputProtocol := s.outputProtocolFactory.GetProtocol(outputTransport)
	defer func() {
		if e := recover(); e != nil {
			log.Printf("panic in processor: %v: %s", e, string(debug.Stack()))
		}
	}()
	if inputTransport != nil {
		defer inputTransport.Close()
	}
	if outputTransport != nil {
		defer outputTransport.Close()
	}
	for {
		if s.IsStop() {
			return nil
		}

		ok, err := processor.Process(inputProtocol, outputProtocol)

		if err, ok := err.(thrift.TTransportException); ok && err.TypeId() == thrift.END_OF_FILE {
			return nil
		} else if err != nil {
			log.Printf("error processing request: %s", err)
			return err
		}

		if !ok {
			break
		}
	}
	return nil
}

func (s *TServer) IsStop() bool {
	return atomic.LoadUint32(&s.stopped) != UNSTOP
}

func (s *TServer) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.IsStop() {
		return nil
	}
	atomic.StoreUint32(&s.stopped, STOP)
	s.serverTransport.Interrupt()
	s.wg.Wait() // 等所有running 状态的goroutine结束

	return nil
}

func (s *TServer) ProcessorFactory() thrift.TProcessorFactory {
	return s.processorFactory
}

func (s *TServer) ServerTransport() thrift.TServerTransport {
	return s.serverTransport
}

func (s *TServer) InputTransportFactory() thrift.TTransportFactory {
	return s.inputTransportFactory
}

func (s *TServer) OutputTransportFactory() thrift.TTransportFactory {
	return s.outputTransportFactory
}

func (s *TServer) InputProtocolFactory() thrift.TProtocolFactory {
	return s.inputProtocolFactory
}

func (s *TServer) OutputProtocolFactory() thrift.TProtocolFactory {
	return s.outputProtocolFactory
}
