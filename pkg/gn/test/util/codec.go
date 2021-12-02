package util

import (
    "encoding/binary"
    "errors"
    "fmt"
    "golang-im/pkg/gn"
    "net"
)

var (
    HeaderLen = 2
    MaxLen    = 1024
)

// Codec 编解码器，用来处理tcp的拆包粘包
type Codec struct {
    Conn    net.Conn
    ReadBuf *gn.Buffer // 读缓冲
}

// NewCodec 创建一个编解码器
func NewCodec(conn net.Conn) *Codec {
    return &Codec{
        Conn:    conn,
        ReadBuf: gn.NewBuffer(make([]byte, MaxLen)),
    }
}

// Read 从conn里面读取数据，当conn发生阻塞，这个方法也会阻塞
func (c *Codec) Read() (int, error) {
    return c.ReadBuf.ReadFromReader(c.Conn)
}

// Decode 解码数据
// Package 代表一个解码包
// bool 标识是否还有可读数据
func (c *Codec) Decode() ([]byte, bool, error) {
    var err error
    // 读取数据长度
    lenBuf, err := c.ReadBuf.Seek(HeaderLen)
    if err != nil {
        return nil, false, nil
    }

    // 读取数据内容
    valueLen := int(binary.BigEndian.Uint16(lenBuf))

    // 数据的字节数组长度大于buffer的长度，返回错误
    if valueLen > MaxLen {
        fmt.Println("out of max len")
        return nil, false, errors.New("out of max len")
    }

    valueBuf, err := c.ReadBuf.Read(HeaderLen, valueLen)
    if err != nil {
        return nil, false, nil
    }
    return valueBuf, true, nil
}

// Encode 编码数据
func Encode(bytes []byte) []byte {
    l := len(bytes)
    buffer := make([]byte, l+2)
    // 将消息长度写入buffer
    binary.BigEndian.PutUint16(buffer[0:2], uint16(l))
    // 将消息内容内容写入buffer
    copy(buffer[2:], bytes)
    return buffer[0 : 2+l]
}
