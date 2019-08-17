package rpc

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	"net"
	"reflect"
	"time"
)

// rpc call client,
// use this to write special struct with message header
type client struct {
	specialColumn    interface{}
	socket           *thrift.TSocket
	handler          func(ts thrift.TTransport, pf thrift.TProtocolFactory) interface{}
	protocolFactory  thrift.TProtocolFactory
	transportFactory thrift.TTransportFactory
}

func NewClient(ip, port string, timeout int64, callback func(ts thrift.TTransport, pf thrift.TProtocolFactory) interface{}, specialColumn interface{}) (*client, error) {
	socket, err := thrift.NewTSocketTimeout(net.JoinHostPort(ip, port), time.Duration(timeout*int64(time.Microsecond)))
	if err != nil {
		return nil, err
	}
	protocolFactory := NewProtocolFactoryDefault()
	transportFactory := NewTransportFactory(thrift.NewTTransportFactory())
	if specialColumn != nil {
		protocolFactory.SetWriteColumn(specialColumn)
	}

	return &client{
		specialColumn:    specialColumn,
		socket:           socket,
		handler:          callback,
		protocolFactory:  protocolFactory,
		transportFactory: transportFactory,
	}, nil
}

func (c *client) Call(method string, parameters ...interface{}) ([]reflect.Value, error) {
	transport := c.transportFactory.GetTransport(c.socket)
	if transport.IsOpen() == false {
		if err := transport.Open(); err != nil {
			return nil, err
		}
	}
	defer transport.Close()

	arguments := make([]reflect.Value, 0, len(parameters))
	for _, value := range parameters {
		arguments = append(arguments, reflect.ValueOf(value))
	}
	res := reflect.ValueOf(c.handler(transport, c.protocolFactory)).MethodByName(method).Call(arguments)

	var err error
	for _, value := range res {
		if _, ok := value.Interface().(error); ok {
			err = value.Interface().(error)
			break
		}
	}

	return res, err
}
