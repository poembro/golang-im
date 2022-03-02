package main

// Start Commond eg: ./client 1 1000 localhost:3101
// first parameter：beginning userId
// second parameter: amount of clients
// third parameter: comet server ip

import (
	"bufio"
	"io"

	"encoding/json"
	"flag"
	"fmt"
	"golang-im/pkg/protocol"
	"math/rand"
	"net"
	"os"
	"runtime"
	"strconv"
	"sync/atomic"
	"time"

	log "github.com/golang/glog"
)

const (
	opHeartbeat      = int32(2)
	opHeartbeatReply = int32(3)
	opAuth           = int32(7)
	opAuthReply      = int32(8)
)

const (
	rawHeaderLen = uint16(16)
	heart        = 15 * time.Second
)

type Int64 int64

// AuthToken auth token.
type AuthToken struct {
	UserId   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
	Face     string `json:"face"`

	RoomID string `json:"room_id"`

	ShopId   string `json:"shop_id"`
	ShopName string `json:"shop_name"`
	ShopFace string `json:"shop_face"`

	Platform string `json:"platform"`
}

var (
	countDown  int64
	aliveCount int64
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	fmt.Println(os.Args)
	begin, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}
	num, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}
	go result()
	for i := begin; i < begin+num; i++ {
		go func(mid int64) {
			for {
				startClient(mid)
				time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
			}
		}(int64(i))
	}
	// signal
	var exit chan bool
	<-exit
}

func result() {
	var (
		lastTimes int64
		interval  = int64(5)
	)
	for {
		nowCount := atomic.LoadInt64(&countDown)
		nowAlive := atomic.LoadInt64(&aliveCount)
		diff := nowCount - lastTimes
		lastTimes = nowCount
		fmt.Println(fmt.Sprintf("%s 活跃连接:%d down:%d down/s:%d", time.Now().Format("2006-01-02 15:04:05"), nowAlive, nowCount, diff/interval))
		time.Sleep(time.Second * time.Duration(interval))
	}
}

func startClient(key int64) {
	time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	atomic.AddInt64(&aliveCount, 1)
	quit := make(chan bool, 1)
	defer func() {
		close(quit)
		atomic.AddInt64(&aliveCount, -1)
	}()
	// connnect to server
	conn, err := net.Dial("tcp", os.Args[3])
	if err != nil {
		log.Errorf("net.Dial(%s) error(%v)", os.Args[3], err)
		return
	}
	seq := int32(0)
	wr := bufio.NewWriter(conn)
	rd := bufio.NewReader(conn)

	authToken := &AuthToken{
		key,
		"xxx",
		"https://img.wxcha.com/m00/86/59/7c6242363084072b82b6957cacc335c7.jpg",
		"1646188353",
		"1646188353",
		"yyyyyshop",
		"http://img.touxiangwu.com/2020/3/uq6Bja.jpg",
		"web",
	}
	proto := new(protocol.Proto)
	proto.Ver = 1
	proto.Op = 7
	proto.Seq = seq
	proto.Body, _ = json.Marshal(authToken)
	if err = tcpWriteProto(wr, proto); err != nil {
		log.Errorf("tcpWriteProto() error(%v)", err)
		return
	}
	if err = tcpReadProto(rd, proto); err != nil {
		log.Errorf("tcpReadProto() error(%v)", err)
		return
	}
	log.Infof("key:%d auth ok, proto: %v", key, proto)
	seq++

	// writer
	go func() {
		for {
			proto.Op = 2
			if err = tcpWriteProto(wr, proto); err != nil {
				log.Errorf("key:%d tcpWriteProto() error(%v)", key, err)
				return
			}
			log.Infof("key:%d Write heartbeat", key)
			time.Sleep(heart)
			seq++
			select {
			case <-quit:
				return
			default:
			}
		}
	}()
	// reader
	for {
		if err = tcpReadProto(rd, proto); err != nil {
			log.Errorf("key:%d tcpReadProto() error(%v)", key, err)
			quit <- true
			return
		}
		if proto.Op == opAuthReply {
			log.Infof("key:%d auth success", key)
		} else if proto.Op == opHeartbeatReply {
			log.Infof("key:%d receive heartbeat", key)
			// 设置读取超时
			//golang的标准网络库是最后期限方式  (平常linux 是空闲超时)
			if err = conn.SetReadDeadline(time.Now().Add(heart + 60*time.Second)); err != nil {
				log.Errorf("conn.SetReadDeadline() error(%v)", err)
				quit <- true
				return
			}
		} else {
			log.Infof("key:%d op:%d msg: %s", key, proto.Op, string(proto.Body))
			atomic.AddInt64(&countDown, 1)
		}
	}
}

func tcpWriteProto(wr *bufio.Writer, proto *protocol.Proto) (err error) {
	// write
	p, err := proto.Encode()

	wr.Write(p)

	fmt.Printf("发送协议包: %#v \r\n", proto.Op)
	//fmt.Println("缓冲中已使用的字节数。:", n, "----", wr.Buffered()) //269
	//fmt.Println("缓冲中还有多少字节未使用。:", wr.Available())         //3827

	err = wr.Flush()
	return
}

func tcpReadProto(rd *bufio.Reader, proto *protocol.Proto) (err error) {
	var p []byte
	_, err = rd.Read(p[0:])
	if err != nil && err != io.EOF {
		return err
	}
	proto.Decode(p)
	return nil
}
