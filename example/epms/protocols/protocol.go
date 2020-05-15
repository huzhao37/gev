package protocols

import (
	"encoding/binary"
	"github.com/Allenxuxu/ringbuffer"
	"github.com/gobwas/pool/pbytes"
	"github.com/huzhao37/gev/connection"
	nc "github.com/huzhao37/gev/example/epms/protocols/gen-go/ncbaseheader"
	"github.com/huzhao37/gev/log"
	"reflect"
	"unsafe"
)

const (
	magicNumber    int = 0x33533
	version        int = 0x2
	checkSum       int = 0x0
	epmsHeaderLen      = 16
	SizeOfEpmsBody     = int(unsafe.Sizeof(nc.NcEPMSMsgHeader{}))
)

//var sizeOfEpmsHeader = int(unsafe.Sizeof(EpmsHeader{}))

type EpmsHeader struct {
	magicNum int //魔法数(0x33533)
	version  int //版本(2)
	BodyLen  int //body长度
	checkSum int //crc32算法(0表示跳过检测)
}

//type EpmsBody struct {
//	MsgType   nc.NcEPMSMsgType // 消息类型
//	MsgName   string           // 消息名称
//	SourceId  int64            // 回复消息和发送结果回复，需要带之前的sourceId
//	ProtoName string           // 消息类型名，用于消息类型校验
//	BufLength int32            // 缓冲块长度 - 【消息内容长度】
//	Buffer    []byte           // 缓存 buf   - 【消息内容二进制块】
//	Option    int32            // 消息选项（0）
//}

type EpmsProtocol struct{}

//just deal header protocol in here,body-dealing in biz,because of it refers to biz msg type
func (d *EpmsProtocol) UnPacket(c *connection.Connection, buffer *ringbuffer.RingBuffer) (interface{}, []byte) {
	if buffer.VirtualLength() > epmsHeaderLen {
		header := pbytes.GetLen(epmsHeaderLen)
		defer pbytes.Put(header)
		_, err := buffer.VirtualRead(header)
		if err != nil {
			log.Error("[epms-unpack]:header error%s", err)
			buffer.VirtualFlush()
			return nil, nil
		}
		epmsHEAD := BytesToEpmsHeader(header)
		if epmsHEAD.magicNum != magicNumber {
			log.Error("[epms-unpack]:magic error%d", epmsHEAD.magicNum)
			buffer.VirtualFlush()
			return nil, nil
		}
		if epmsHEAD.version != version {
			log.Error("[epms-unpack]:version error%d", epmsHEAD.version)
			buffer.VirtualFlush()
			return nil, nil
		}
		//crc32校验
		if epmsHEAD.checkSum != checkSum {
			log.Error("[epms-unpack]:checkSum error%d", epmsHEAD.checkSum)
			buffer.VirtualFlush()
			return nil, nil
		}
		//解析body报文
		if buffer.VirtualLength() == epmsHEAD.BodyLen {
			//返回body数据
			body := make([]byte, epmsHEAD.BodyLen)
			_, err := buffer.VirtualRead(body)
			if err != nil {
				log.Error("[epms-unpack]:header error%s", err)
				buffer.VirtualFlush()
				return nil, nil
			}
			buffer.VirtualFlush()
			return nil, body
		} else {
			//数据长度不对应
			log.Error("[epms-unpack]:bodyLen error%d", epmsHEAD.BodyLen)
			buffer.VirtualRevert()
			return nil, nil
		}
	}
	return nil, nil
}

//param:data is epms-body,the func add epms-header
func (d *EpmsProtocol) Packet(c *connection.Connection, data []byte) []byte {
	dataLen := len(data)
	header := &EpmsHeader{magicNum: 0x33533, version: 2, checkSum: 0, BodyLen: dataLen}
	ret := make([]byte, epmsHeaderLen+dataLen)
	binary.LittleEndian.PutUint32(ret, uint32(dataLen))
	copy(ret[:epmsHeaderLen], epmsHeaderToBytes(header))
	copy(ret[epmsHeaderLen:], data)
	return ret
}

//convert//
//model to []byte
func epmsHeaderToBytes(s *EpmsHeader) []byte {
	data := pbytes.GetLen(epmsHeaderLen)
	defer pbytes.Put(data)
	//data := bufferPool.Get(16)
	binary.LittleEndian.PutUint32((data)[0:4], uint32(s.magicNum))
	binary.LittleEndian.PutUint32((data)[4:8], uint32(s.version))
	binary.LittleEndian.PutUint32((data)[8:12], uint32(s.BodyLen))
	binary.LittleEndian.PutUint32((data)[12:16], uint32(s.checkSum))
	return data
	//var x reflect.SliceHeader
	//x.Len = epmsHeaderLen
	//x.Cap = epmsHeaderLen
	//x.Data = uintptr(unsafe.Pointer(s))
	//return *(*[]byte)(unsafe.Pointer(&x))
}
func EpmsBodyToBytes(s *nc.NcEPMSMsgHeader) []byte {
	var x reflect.SliceHeader
	x.Len = SizeOfEpmsBody
	x.Cap = SizeOfEpmsBody
	x.Data = uintptr(unsafe.Pointer(s))
	return *(*[]byte)(unsafe.Pointer(&x))
}

//[]byte to model
func BytesToEpmsHeader(b []byte) *EpmsHeader {
	return &EpmsHeader{
		int(binary.LittleEndian.Uint32(b[0:4])),
		int(binary.LittleEndian.Uint32(b[4:8])),
		int(binary.LittleEndian.Uint32(b[8:12])),
		int(binary.LittleEndian.Uint32(b[12:16])),
	}
	//return (*EpmsHeader)(unsafe.Pointer(
	//	(*reflect.SliceHeader)(unsafe.Pointer(&b)).Data,
	//))
}

func BytesToEpmsBody(b []byte) *nc.NcEPMSMsgHeader {
	return (*nc.NcEPMSMsgHeader)(unsafe.Pointer(
		(*reflect.SliceHeader)(unsafe.Pointer(&b)).Data,
	))
}
