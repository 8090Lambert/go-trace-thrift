package rpc

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/8090Lambert/tracer_thrift/helper"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

func NewServer(processor thrift.TProcessor, specialColumn interface{}) *TServer {
	helper.CheckAddrAlreadyUse(vars.port)
	helper.SetPid2File(vars.pidFile)
	if specialColumn != nil {
		custom.SetReadColumn(specialColumn)
	}
	return initServer(processor, vars.port, vars.timeout)
}

func (s *TServer) Run() {
	exit := make(chan struct{}, 1)
	defer close(exit)
	go func() {
		log.Printf("thrift server will run on addr=%s", vars.port)
		if err := s.Serve(); err != nil {
			exit <- struct{}{}
			log.Printf("fail to start thrift server. err:%v addr:s%s", err, vars.port)
		}
		exit <- struct{}{}
	}()

	c := make(chan os.Signal)
	signal.Notify(c, vars.captureSign...)
	defer signal.Stop(c)

	select {
	case sign := <-c:
		if err := s.Stop(); err != nil {
			log.Printf("server stop failed, err: %s", err)
		}
		log.Printf("server graceful stop, signal=%v", sign)
	}
}

func initServer(processor thrift.TProcessor, addr string, timeout int64) *TServer {
	server := new(TServer)
	server.bindAddr = net.JoinHostPort("", addr)
	server.timeout = timeout
	server.quit = make(chan struct{}, 1)
	server.initAttribute(processor)
	return server
}

func (s *TServer) initAttribute(processor thrift.TProcessor) {
	s.initProcessor(processor)
	s.initTransport()
	s.initProtocol()
}

func (s *TServer) initProcessor(processor thrift.TProcessor) {
	s.processor = NewProcessorDefault(processor)
	s.processorFactory = thrift.NewTProcessorFactory(s.processor)
}

func (s *TServer) initTransport() {
	serverTransport, err := thrift.NewTServerSocketTimeout(s.bindAddr, time.Duration(s.timeout*int64(time.Millisecond)))
	if err != nil {
		log.Printf("new server socket on addr %s failed, error:%s", s.bindAddr, err)
		os.Exit(1)
	}
	s.serverTransport = serverTransport
	s.serverTransport.BufferSize = 0
	transportFactory := NewTransportFactory(thrift.NewTTransportFactory())
	s.inputTransportFactory = transportFactory
	s.outputTransportFactory = transportFactory
}

func (s *TServer) initProtocol() {
	protocolFactory := NewProtocolFactoryDefault()
	s.inputProtocolFactory = protocolFactory
	s.outputProtocolFactory = protocolFactory
}
