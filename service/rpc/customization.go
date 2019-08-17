package rpc

import (
	"errors"
	"git.apache.org/thrift.git/lib/go/thrift"
	"reflect"
)

var (
	custom      *Customization
	ReadMethod  = `Read`
	WriteMethod = `Write`
)

type Customization struct {
	Reader *ReadExecutor
	Writer *WriteExecutor
}

type ColumnType interface{}

type ReadClosure func(cp *Protocol, column ColumnType) (interface{}, thrift.TTransport, error)

type WriteClosure func(cp *Protocol, column ColumnType) (interface{}, error)

type ReadExecutor struct {
	Column      ColumnType
	HasCallBack bool
	CallBack    ReadClosure
}

type WriteExecutor struct {
	Column      ColumnType
	HasCallBack bool
	CallBack    WriteClosure
}

func init() {
	custom = new(Customization)
	custom.Reader = new(ReadExecutor)
	custom.Writer = new(WriteExecutor)
}

// Must ptr
func RegisterReadCustom(column ColumnType) error {
	return custom.SetReadColumn(column)
}

func (c *Customization) SetReadColumn(column ColumnType) error {
	reflectField := reflect.ValueOf(column)
	if reflectField.Kind().String() != `ptr` {
		return errors.New("`column` must ptr")
	}
	c.Reader.Column = column
	return nil
}

func (c *Customization) SetWriteColumn(column ColumnType) error {
	reflectField := reflect.ValueOf(column)
	if reflectField.Kind().String() != `ptr` {
		return errors.New("`column` must ptr")
	}
	c.Writer.Column = column
	return nil
}

// Only one column
func (c *Customization) ReadCallBackStub() ReadClosure {
	return func(cp *Protocol, column ColumnType) (interface{}, thrift.TTransport, error) {
		// when read stream,
		// and write to another way
		// because Transport.tran is private,
		// process need a complete transport to parse.
		originalTransport := cp.Transport()
		wrapperTransport := NewStreamReWinder(originalTransport)
		pc := NewProtocolDefault(wrapperTransport)

		_, _, _, err := pc.ReadMessageBegin()
		defer pc.ReadMessageEnd()

		if err != nil {
			return nil, originalTransport, err
		}

		for {
			_, fieldType, fieldSeq, err := pc.ReadFieldBegin()
			if err != nil || fieldType == thrift.STOP {
				pc.ReadFieldEnd()
				break
			}

			if fieldSeq == 0 && fieldType == thrift.STRUCT {
				reflectField := reflect.ValueOf(column)
				args := []reflect.Value{reflect.ValueOf(pc.protocol)}
				reflectField.MethodByName(ReadMethod).Call(args)
			} else {
				pc.protocol.Skip(fieldType)
			}
		}

		return column, wrapperTransport.Output, nil
	}
}

func (c *Customization) WriteCallBackStub() WriteClosure {
	return func(cp *Protocol, column ColumnType) (interface{}, error) {
		cp.protocol.WriteFieldBegin("custom", thrift.STRUCT, 0)
		reflectColumn := reflect.ValueOf(column)
		result := reflectColumn.MethodByName(WriteMethod).Call([]reflect.Value{reflect.ValueOf(cp.protocol)})
		cp.protocol.WriteFieldEnd()

		var err error
		if _, ok := result[0].Interface().(error); ok {
			err = result[0].Interface().(error)
		}

		return nil, err
	}
}
