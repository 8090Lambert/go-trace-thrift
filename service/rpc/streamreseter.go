package rpc

import "git.apache.org/thrift.git/lib/go/thrift"

type StreamReWinder struct {
	thrift.TTransport
	Input  thrift.TTransport
	Output thrift.TTransport
	retry  int // retry count
}

func NewStreamReWinder(original thrift.TTransport) *StreamReWinder {
	pipe := new(StreamReWinder)
	pipe.Input = original
	pipe.Output = thrift.NewTMemoryBuffer()
	pipe.retry = 3
	return pipe
}

func (pipe *StreamReWinder) IsOpen() bool {
	return pipe.Input.IsOpen()
}

func (pipe *StreamReWinder) Open() error {
	return pipe.Input.Open()
}

func (pipe *StreamReWinder) Close() error {
	return pipe.Input.Close()
}

func (pipe *StreamReWinder) Read(b []byte) (n int, err error) {
	n, err = pipe.Input.Read(b)
	if err == nil {
		if _, err = pipe.Output.Write(b); err != nil {
			for i := 0; i < pipe.retry; i++ {
				if _, err = pipe.Output.Write(b); err == nil {
					break
				}
			}
		}
	}
	return
}

func (pipe *StreamReWinder) Write(b []byte) (int, error) {
	return pipe.Input.Write(b)
}

func (pipe *StreamReWinder) Flush() error {
	return pipe.Output.Flush()
}
