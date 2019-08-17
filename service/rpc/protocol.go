package rpc

import (
	"errors"
	"git.apache.org/thrift.git/lib/go/thrift"
	"log"
	"reflect"
)

type ProtocolFactory struct {
	protocolFactory *thrift.TBinaryProtocolFactory
	writeColumn     interface{}
}

func NewProtocolFactoryDefault() *ProtocolFactory {
	return &ProtocolFactory{
		protocolFactory: thrift.NewTBinaryProtocolFactoryDefault(),
		writeColumn:     nil,
	}
}

func NewProtocolFactory(strictRead, strictWrite bool) *ProtocolFactory {
	return &ProtocolFactory{
		protocolFactory: thrift.NewTBinaryProtocolFactory(strictRead, strictWrite),
		writeColumn:     nil,
	}
}

func (cpf *ProtocolFactory) GetProtocol(t thrift.TTransport) thrift.TProtocol {
	protocol := NewProtocolDefault(t)
	protocol.SetWriteColumn(cpf.writeColumn)
	return protocol
}

func (cpf *ProtocolFactory) SetWriteColumn(column interface{}) error {
	reflectColumn := reflect.ValueOf(column)
	if reflectColumn.Kind().String() != `ptr` {
		return errors.New("`column` must ptr")
	}
	cpf.writeColumn = column
	return nil
}

type Protocol struct {
	protocol    *thrift.TBinaryProtocol
	seqId       int32
	writeColumn bool
}

func NewProtocolDefault(t thrift.TTransport) *Protocol {
	return &Protocol{
		protocol:    thrift.NewTBinaryProtocol(t, false, true),
		seqId:       0,
		writeColumn: false,
	}
}

func NewProtocol(t thrift.TTransport, strictRead, strictWrite bool) *Protocol {
	return &Protocol{
		protocol:    thrift.NewTBinaryProtocol(t, strictRead, strictWrite),
		seqId:       0,
		writeColumn: false,
	}
}

func (cp *Protocol) SetWriteColumn(column interface{}) error {
	if err := custom.SetWriteColumn(column); err != nil {
		return err
	}
	cp.writeColumn = true
	return nil
}

func (cp *Protocol) Rewind() {
	cp.seqId = 0
}

func (cp *Protocol) WriteMessageBegin(name string, typeId thrift.TMessageType, seqid int32) error {
	cp.Rewind()
	return cp.protocol.WriteMessageBegin(name, typeId, seqid)
}

func (cp *Protocol) WriteMessageEnd() error {
	return cp.protocol.WriteMessageEnd()
}

func (cp *Protocol) WriteStructBegin(name string) error {
	if cp.writeColumn == false {
		return cp.protocol.WriteStructBegin(name)
	} else {
		if cp.seqId != 0 {
			return cp.protocol.WriteStructBegin(name)
		}
		cp.seqId++

		// Write special column
		writeCallback := custom.Writer.CallBack
		if custom.Writer.HasCallBack == false {
			writeCallback = custom.WriteCallBackStub()
		}

		if _, err := writeCallback(cp, custom.Writer.Column); err != nil {
			log.Fatal(err)
		}

		return cp.protocol.WriteStructBegin(name)
	}
}

func (cp *Protocol) WriteStructEnd() error {
	return cp.protocol.WriteStructEnd()
}

func (cp *Protocol) WriteFieldBegin(name string, typeId thrift.TType, id int16) error {
	return cp.protocol.WriteFieldBegin(name, typeId, id)
}

func (cp *Protocol) WriteFieldEnd() error {
	return cp.protocol.WriteFieldEnd()
}

func (cp *Protocol) WriteFieldStop() error {
	return cp.protocol.WriteFieldStop()
}

func (cp *Protocol) WriteMapBegin(keyType thrift.TType, valueType thrift.TType, size int) error {
	return cp.protocol.WriteMapBegin(keyType, valueType, size)
}

func (cp *Protocol) WriteMapEnd() error {
	return cp.protocol.WriteMapEnd()
}

func (cp *Protocol) WriteListBegin(elemType thrift.TType, size int) error {
	return cp.protocol.WriteListBegin(elemType, size)
}

func (cp *Protocol) WriteListEnd() error {
	return cp.protocol.WriteListEnd()
}

func (cp *Protocol) WriteSetBegin(elemType thrift.TType, size int) error {
	return cp.protocol.WriteSetBegin(elemType, size)
}

func (cp *Protocol) WriteSetEnd() error {
	return cp.protocol.WriteSetEnd()
}

func (cp *Protocol) WriteBool(value bool) error {
	return cp.protocol.WriteBool(value)
}

func (cp *Protocol) WriteByte(value byte) error {
	return cp.protocol.WriteByte(value)
}

func (cp *Protocol) WriteI16(value int16) error {
	return cp.protocol.WriteI16(value)
}

func (cp *Protocol) WriteI32(value int32) error {
	return cp.protocol.WriteI32(value)
}

func (cp *Protocol) WriteI64(value int64) error {
	return cp.protocol.WriteI64(value)
}

func (cp *Protocol) WriteDouble(value float64) error {
	return cp.protocol.WriteDouble(value)
}

func (cp *Protocol) WriteString(value string) error {
	return cp.protocol.WriteString(value)
}

func (cp *Protocol) WriteBinary(value []byte) error {
	return cp.protocol.WriteBinary(value)
}

func (cp *Protocol) ReadMessageBegin() (name string, typeId thrift.TMessageType, seqid int32, err error) {
	return cp.protocol.ReadMessageBegin()
}

func (cp *Protocol) ReadMessageEnd() error {
	return cp.protocol.ReadMessageEnd()
}

func (cp *Protocol) ReadStructBegin() (name string, err error) {
	return cp.protocol.ReadStructBegin()
}

func (cp *Protocol) ReadStructEnd() error {
	return cp.protocol.ReadStructEnd()
}

func (cp *Protocol) ReadFieldBegin() (name string, typeId thrift.TType, id int16, err error) {
	return cp.protocol.ReadFieldBegin()
}

func (cp *Protocol) ReadFieldEnd() error {
	return cp.protocol.ReadFieldEnd()
}

func (cp *Protocol) ReadMapBegin() (keyType thrift.TType, valueType thrift.TType, size int, err error) {
	return cp.protocol.ReadMapBegin()
}

func (cp *Protocol) ReadMapEnd() error {
	return cp.protocol.ReadMapEnd()
}

func (cp *Protocol) ReadListBegin() (elemType thrift.TType, size int, err error) {
	return cp.protocol.ReadListBegin()
}

func (cp *Protocol) ReadListEnd() error {
	return cp.protocol.ReadListEnd()
}

func (cp *Protocol) ReadSetBegin() (elemType thrift.TType, size int, err error) {
	return cp.protocol.ReadSetBegin()
}

func (cp *Protocol) ReadSetEnd() error {
	return cp.protocol.ReadSetEnd()
}

func (cp *Protocol) ReadBool() (value bool, err error) {
	return cp.protocol.ReadBool()
}

func (cp *Protocol) ReadByte() (value byte, err error) {
	return cp.protocol.ReadByte()
}

func (cp *Protocol) ReadI16() (value int16, err error) {
	return cp.protocol.ReadI16()
}

func (cp *Protocol) ReadI32() (value int32, err error) {
	return cp.protocol.ReadI32()
}

func (cp *Protocol) ReadI64() (value int64, err error) {
	return cp.protocol.ReadI64()
}

func (cp *Protocol) ReadDouble() (value float64, err error) {
	return cp.protocol.ReadDouble()
}

func (cp *Protocol) ReadString() (value string, err error) {
	return cp.protocol.ReadString()
}

func (cp *Protocol) ReadBinary() (value []byte, err error) {
	return cp.protocol.ReadBinary()
}

func (cp *Protocol) Skip(fieldType thrift.TType) (err error) {
	return cp.protocol.Skip(fieldType)
}

func (cp *Protocol) Flush() (err error) {
	return cp.protocol.Flush()
}

func (cp *Protocol) Transport() thrift.TTransport {
	return cp.protocol.Transport()
}
