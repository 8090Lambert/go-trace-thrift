package rpc

import "git.apache.org/thrift.git/lib/go/thrift"

type TransportFactory struct {
	transportFactory thrift.TTransportFactory
}

type Transport struct {
	transport *thrift.TFramedTransport
}

func NewTransportFactory(factory thrift.TTransportFactory) *TransportFactory {
	return &TransportFactory{thrift.NewTFramedTransportFactory(factory)}
}

func (factory *TransportFactory) GetTransport(base thrift.TTransport) thrift.TTransport {
	return NewTransport(base)
}

func NewTransport(transport thrift.TTransport) *Transport {
	return &Transport{
		thrift.NewTFramedTransport(transport),
	}
}

func (t *Transport) Open() error {
	return t.transport.Open()
}

func (t *Transport) IsOpen() bool {
	return t.transport.IsOpen()
}

func (t *Transport) Close() error {
	return t.transport.Close()
}

func (t *Transport) Read(buf []byte) (int, error) {
	return t.transport.Read(buf)
}

func (t *Transport) ReadByte() (c byte, err error) {
	return t.transport.ReadByte()
}

func (t *Transport) Write(p []byte) (int, error) {
	return t.transport.Write(p)
}

func (t *Transport) WriteByte(c byte) error {
	return t.transport.WriteByte(c)
}

func (t *Transport) WriteString(s string) (int, error) {
	return t.transport.WriteString(s)
}

func (t *Transport) Flush() error {
	return t.transport.Flush()
}
