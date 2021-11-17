package gn

import (
	"encoding/binary"
	"io"
	"sync"
	"syscall"
)

// Decoder 解码器
type Decoder interface {
	Decode(c *Conn) error
}

// Encoder 编码器
type Encoder interface {
	EncodeToFD(fd int32, bytes []byte) error
}

type headerLenDecoder struct {
	headerLen int // TCP包的头部长度，用来描述这个包的字节长度
}

// NewHeaderLenDecoder 创建基于头部长度的解码器
// headerLen TCP包的头部内容，用来描述这个包的字节长度
// readMaxLen 所读取的客户端包的最大长度，客户端发送的包不能超过这个长度
func NewHeaderLenDecoder(headerLen int) Decoder {
	if headerLen <= 0 {
		panic("headerLen or readMaxLen must must greater than 0")
	}

	return &headerLenDecoder{
		headerLen: headerLen,
	}
}

func Int32(b []byte) int32 {
	return int32(b[3]) | int32(b[2])<<8 | int32(b[1])<<16 | int32(b[0])<<24
}

// Decode 解码 (golang-im自定义一下)
func (d *headerLenDecoder) Decode(c *Conn) error {
	for {
		header, err := c.buffer.Seek(4)
		if err == ErrNotEnough {
			return nil
		}
		valueLen := Int32(header)
		value, err := c.buffer.Read(0, int(valueLen))
		if err == ErrNotEnough {
			return nil
		}

		c.server.handler.OnMessage(c, value)
	}
}

// Decode 解码
func (d *headerLenDecoder) Decode2(c *Conn) error {
	for {
		header, err := c.buffer.Seek(d.headerLen)
		if err == ErrNotEnough {
			return nil
		}
		valueLen := int(binary.BigEndian.Uint16(header))
		value, err := c.buffer.Read(d.headerLen, valueLen)
		if err == ErrNotEnough {
			return nil
		}

		c.server.handler.OnMessage(c, value)
	}
}

type headerLenEncoder struct {
	headerLen       int        // TCP包的头部长度，用来描述这个包的字节长度
	writeBufferLen  int        // 服务器发送给客户端包的建议长度，当发送的包小于这个值时，会利用到内存池优化
	writeBufferPool *sync.Pool // 写缓存区内存池
}

// NewHeaderLenEncoder 创建基于头部长度的编码器
// headerLen TCP包的头部内容，用来描述这个包的字节长度
// writeBufferLen 服务器发送给客户端包的建议长度，当发送的包小于这个值时，会利用到内存池优化
func NewHeaderLenEncoder(headerLen, writeBufferLen int) *headerLenEncoder {
	if headerLen <= 0 || writeBufferLen <= 0 {
		panic("headerLen or writeBufferLen must must greater than 0")
	}

	return &headerLenEncoder{
		headerLen:      headerLen,
		writeBufferLen: writeBufferLen,
		writeBufferPool: &sync.Pool{
			New: func() interface{} {
				b := make([]byte, writeBufferLen)
				return b
			},
		},
	}
}

// EncodeToFD 编码数据,并且写入文件描述符
func (e headerLenEncoder) EncodeToFD(fd int32, bytes []byte) error {
	l := len(bytes)
	var buffer []byte
	if l <= e.writeBufferLen-e.headerLen {
		obj := e.writeBufferPool.Get()
		defer e.writeBufferPool.Put(obj)
		buffer = obj.([]byte)[0 : l+e.headerLen]
	} else {
		buffer = make([]byte, l+e.headerLen)
	}

	// 将消息长度写入buffer
	binary.BigEndian.PutUint16(buffer[0:2], uint16(l))
	// 将消息内容内容写入buffer
	copy(buffer[e.headerLen:], bytes)

	_, err := syscall.Write(int(fd), buffer)
	return err
}

// EncodeToWriter 编码数据,并且写入Writer
func (e headerLenEncoder) EncodeToWriter2(w io.Writer, bytes []byte) error {
	l := len(bytes)
	var buffer []byte
	if l <= e.writeBufferLen-e.headerLen {
		obj := e.writeBufferPool.Get()
		defer e.writeBufferPool.Put(obj)
		buffer = obj.([]byte)[0 : l+e.headerLen]
	} else {
		buffer = make([]byte, l+e.headerLen)
	}

	// 将消息长度写入buffer
	binary.BigEndian.PutUint16(buffer[0:2], uint16(l))
	// 将消息内容内容写入buffer
	copy(buffer[e.headerLen:], bytes)

	_, err := w.Write(buffer)
	return err
}

// EncodeToWriter 编码数据,并且写入Writer (golang-im自定义一下)
func (e headerLenEncoder) EncodeToWriter(w io.Writer, bytes []byte) error {
	l := len(bytes)
	var buffer []byte
	if l <= e.writeBufferLen-e.headerLen {
		obj := e.writeBufferPool.Get()
		defer e.writeBufferPool.Put(obj)
		buffer = obj.([]byte)[0:l]
	} else {
		buffer = make([]byte, l)
	}

	// 将消息内容内容写入buffer
	copy(buffer[0:], bytes)

	_, err := w.Write(buffer)
	return err
}
