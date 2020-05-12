package main

import (
	"container/list"
	"context"
	"flag"
	"github.com/Allenxuxu/toolkit/sync/atomic"
	"github.com/huzhao37/gev/example/epms/protocols"
	"github.com/huzhao37/gev/example/epms/queue"
	"github.com/leandro-lugaresi/hub"
	"log"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/huzhao37/gev"
	"github.com/huzhao37/gev/connection"
)

const (
	SystemRead  = "SystemRead"
	SystemWrite = "SystemWrite"
	BizRead     = "BizRead"
	BizWrite    = "BizWrite"
)

type EpmsServer struct {
	mu                 sync.RWMutex
	conn               *list.List
	clientNum          atomic.Int64
	maxConnection      int64
	server             *gev.Server
	ReadQueue          *queue.List
	SystemWriteQueue   *queue.List
	BusinessWriteQueue *queue.List
}

// New Epms Server
func NewEpmsServer(ip string, port int, maxConnection int64, loops int) (*EpmsServer, error) {
	var err error
	s := new(EpmsServer)
	s.maxConnection = maxConnection
	s.server, err = gev.NewServer(s,
		gev.Network("tcp"),
		gev.Address(":"+strconv.Itoa(port)),
		gev.NumLoops(runtime.NumCPU()*3/5), //loops
		gev.Protocol(&protocols.EpmsProtocol{}),
		gev.IdleTime(5*time.Second),
		gev.ReusePort(true))
	if err != nil {
		return nil, err
	}

	return s, nil
}

// Start server
func (s *EpmsServer) Start() {
	s.ReadQueue = queue.CreateNew(SystemRead, 10_000)
	s.SystemWriteQueue = queue.CreateNew(SystemWrite, 10_000)
	s.BusinessWriteQueue = queue.CreateNew(BizWrite, 10_000)

	s.server.RunEvery(5*time.Second, s.RunPush) //定时发送给hello消息
	s.server.Start()
}

// Stop server
func (s *EpmsServer) Stop() {
	s.server.Stop()
}

func (s *EpmsServer) OnConnect(c *connection.Connection) {
	s.clientNum.Add(1)
	log.Println(" OnConnect ： ", c.PeerAddr())

	if s.clientNum.Get() > s.maxConnection {
		_ = c.ShutdownWrite()
		log.Println("Refused connection")
		return
	}
	s.mu.Lock()
	e := s.conn.PushBack(c)
	s.mu.Unlock()
	c.SetContext(e)

	//订阅
	s.ReadQueue.Subscribe(c.PeerAddr(), s.DistributingMsg)
	s.SystemWriteQueue.Subscribe(c.PeerAddr(), s.SystemHandlerWrite)
	s.BusinessWriteQueue.Subscribe(c.PeerAddr(), s.HandlerWrite)
}

func (s *EpmsServer) OnMessage(c *connection.Connection, ctx interface{}, data []byte) (out []byte) {
	log.Println("OnMessage：", data)
	//发布
	s.ReadQueue.Publish(c.PeerAddr(), data)
	out = data
	return
}

func (s *EpmsServer) OnClose(c *connection.Connection) {
	s.clientNum.Add(-1)
	log.Println("OnClose")
	e := c.Context().(*list.Element)
	s.mu.Lock()
	s.conn.Remove(e)
	s.mu.Unlock()
}

// RunPush push message
func (s *EpmsServer) RunPush() {
	var next *list.Element

	s.mu.RLock()
	defer s.mu.RUnlock()

	for e := s.conn.Front(); e != nil; e = next {
		next = e.Next()

		c := e.Value.(*connection.Connection)
		_ = c.Send([]byte("hello\n"))
	}
}

func main() {
	var port int
	var loops int
	var maxClients int64

	flag.IntVar(&port, "port", 1833, "server port")
	flag.IntVar(&loops, "loops", -1, "num loops")
	flag.Int64Var(&maxClients, "maxClients", 1000, "max clients")
	flag.Parse()

	s, err := NewEpmsServer("127.0.0.1", port, maxClients, loops)
	if err != nil {
		panic(err)
	}

	log.Println("server start")
	s.Start()
}

func (s *EpmsServer) RegisterService(serviceName string, svc interface{}, metaData string) error {
	//todo
	return InprocessClient.Register(serviceName, svc, metaData)
}

func (s *EpmsServer) UnRegisterService(serviceName string) error {
	//todo
	return InprocessClient.Unregister(serviceName)
}

func (s *EpmsServer) SystemProcessor(ctx context.Context, servicePath, serviceMethod string, args interface{}, reply interface{}) *Call {
	return InprocessClient.Go(ctx, servicePath, serviceMethod, args, reply, nil)
}

func (s *EpmsServer) BusinessProcessor(ctx context.Context, servicePath, serviceMethod string, args interface{}, reply interface{}) *Call {
	return InprocessClient.Go(ctx, servicePath, serviceMethod, args, reply, nil)
}

//***add async queue for solve msg order***
//according to msg type ,write data to bussiness—read-queue  or system-read-queue
func (s *EpmsServer) DistributingMsg(addr string, msg hub.Message) {
	//todo
	data := msg.Fields["data"].([]byte)
	//1.unpack(协议中处理后的bytes，再次转换成thrift struct)
	epmsBody := protocols.BytesToEpmsBody(data)
	//2.get msg type
	//3.wirte msg data to queue
	//系统消息写入系统队列，业务消息写入业务队列
	if epmsBody.MsgType == protocols.NC_EPMS_HEARTBEAT {
		s.SystemWriteQueue.Publish(addr, data)
	} else {
		s.BusinessWriteQueue.Publish(addr, data)
	}
}

func (s *EpmsServer) HandlerWrite(addr string, msg hub.Message) {
	//todo
	//1.msg unpack
	//2.valid
	//3.call BusinessProcessor,return result
	//4.send client
}

func (s *EpmsServer) SystemHandlerWrite(addr string, msg hub.Message) {
	//todo
	//1.msg unpack
	//2.valid
	//3.call SystemProcessor,return result
	//4.send client
}
