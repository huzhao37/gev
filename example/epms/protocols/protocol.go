package protocols

import (
	"encoding/binary"
	"github.com/Allenxuxu/ringbuffer"
	"github.com/gobwas/pool/pbytes"
	"github.com/huzhao37/gev/connection"
	"github.com/huzhao37/gev/log"
	"reflect"
	"unsafe"
)

const epmsHeaderLen = 16
const SizeOfEpmsBody = int(unsafe.Sizeof(EpmsBody{}))

//var sizeOfEpmsHeader = int(unsafe.Sizeof(EpmsHeader{}))

/**
消息类型
**/
type ncEPMSMsgType int32

const (
	NC_EPMS_SEND_MSG      ncEPMSMsgType = 0 // 发送消息
	NC_EPMS_REPLY_MSG     ncEPMSMsgType = 1 // 回复消息
	NC_EPMS_SEND_SUCCESS  ncEPMSMsgType = 2 // 发送结果 - 成功
	NC_EPMS_SEND_FAILED   ncEPMSMsgType = 3 // 发送结果 - 失败
	NC_EPMS_NO_SUBSCRIBER ncEPMSMsgType = 4 // 没有订阅对象
	NC_EPMS_CONNECT       ncEPMSMsgType = 5 // 连接消息
	NC_EPMS_DISCONNECT    ncEPMSMsgType = 6 // 断开消息
	NC_EPMS_HEARTBEAT     ncEPMSMsgType = 7 // 心跳消息
)

func (p ncEPMSMsgType) String() string {
	switch p {
	case NC_EPMS_SEND_MSG:
		return "NC_EPMS_SEND_MSG"
	case NC_EPMS_REPLY_MSG:
		return "NC_EPMS_REPLY_MSG"
	case NC_EPMS_SEND_SUCCESS:
		return "NC_EPMS_SEND_SUCCESS"
	case NC_EPMS_SEND_FAILED:
		return "NC_EPMS_SEND_FAILED"
	case NC_EPMS_NO_SUBSCRIBER:
		return "NC_EPMS_NO_SUBSCRIBER"
	case NC_EPMS_CONNECT:
		return "NC_EPMS_CONNECT"
	case NC_EPMS_DISCONNECT:
		return "NC_EPMS_DISCONNECT"
	case NC_EPMS_HEARTBEAT:
		return "NC_EPMS_HEARTBEAT"
	default:
		return "UNKNOWN"
	}
}

type EpmsHeader struct {
	magicNum int //魔法数(0x33533)
	version  int //版本(2)
	bodyLen  int //body长度
	checkSum int //crc32算法(0表示跳过检测)
}

type EpmsBody struct {
	MsgType   ncEPMSMsgType // 消息类型
	MsgName   string        // 消息名称
	SourceId  int64         // 回复消息和发送结果回复，需要带之前的sourceId
	ProtoName string        // 消息类型名，用于消息类型校验
	BufLength int32         // 缓冲块长度 - 【消息内容长度】
	Buffer    []byte        // 缓存 buf   - 【消息内容二进制块】
	Option    int32         // 消息选项（0）
}

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
			epmsHEAD = bytesToEpmsHeader(header)
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
			if int(dataLen) != epmsHEAD.bodyLen+16 {
				log.Error("[epms-unpack]:bodyLen error%d", epmsHEAD.bodyLen)
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
	header := &EpmsHeader{magicNum: 0x33533, version: 2, checkSum: 0, bodyLen: dataLen}
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
func EpmsBodyToBytes(s *EpmsBody) []byte {
	var x reflect.SliceHeader
	x.Len = SizeOfEpmsBody
	x.Cap = SizeOfEpmsBody
	x.Data = uintptr(unsafe.Pointer(s))
	return *(*[]byte)(unsafe.Pointer(&x))
}

//[]byte to model
func bytesToEpmsHeader(b []byte) *EpmsHeader {
	return (*EpmsHeader)(unsafe.Pointer(
		(*reflect.SliceHeader)(unsafe.Pointer(&b)).Data,
	))
}

func BytesToEpmsBody(b []byte) *EpmsBody {
	return (*EpmsBody)(unsafe.Pointer(
		(*reflect.SliceHeader)(unsafe.Pointer(&b)).Data,
	))
}
