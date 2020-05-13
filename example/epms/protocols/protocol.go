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

const epmsHeaderLen = 16
const SizeOfEpmsBody = int(unsafe.Sizeof(nc.NcEPMSMsgHeader{}))

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
		buf := pbytes.GetLen(epmsHeaderLen)
		defer pbytes.Put(buf)
		_, _ = buffer.VirtualRead(buf)
		dataLen := binary.LittleEndian.Uint32(buf) //小端字节流(需要重写ringBuffer包)
		//解析header报文
		if buffer.VirtualLength() >= int(dataLen) {
			//todo
			epmsHEAD := &EpmsHeader{}
			header := make([]byte, 16)
			_, err := buffer.VirtualRead(header)
			if err != nil {
				log.Error("[epms-unpack]:%s", err)
			}
			epmsHEAD = BytesToEpmsHeader(header)
			if epmsHEAD.magicNum != 0x33533 {
				log.Error("[epms-unpack]:magic error%d", epmsHEAD.magicNum)
				buffer.VirtualFlush()
				return nil, nil
			}
			if epmsHEAD.version != 2 {
				log.Error("[epms-unpack]:version error%d", epmsHEAD.version)
				buffer.VirtualFlush()
				return nil, nil
			}
			//crc32校验
			if epmsHEAD.checkSum != 0 {
				//todo
			}
			//数据长度不对应
			if int(dataLen) != epmsHEAD.BodyLen+16 {
				log.Error("[epms-unpack]:bodyLen error%d", epmsHEAD.BodyLen)
				buffer.VirtualFlush()
				return nil, nil
			}
			//返回body数据
			body := make([]byte, dataLen-16)
			_, _ = buffer.VirtualRead(body)

			buffer.VirtualFlush()
			return nil, body
		} else {
			buffer.VirtualRevert()
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
	var x reflect.SliceHeader
	x.Len = epmsHeaderLen
	x.Cap = epmsHeaderLen
	x.Data = uintptr(unsafe.Pointer(s))
	return *(*[]byte)(unsafe.Pointer(&x))
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
	return (*EpmsHeader)(unsafe.Pointer(
		(*reflect.SliceHeader)(unsafe.Pointer(&b)).Data,
	))
}

func BytesToEpmsBody(b []byte) *nc.NcEPMSMsgHeader {
	return (*nc.NcEPMSMsgHeader)(unsafe.Pointer(
		(*reflect.SliceHeader)(unsafe.Pointer(&b)).Data,
	))
}
