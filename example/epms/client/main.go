package main

import (
	"encoding/binary"
	"fmt"
	"github.com/huzhao37/gev/example/epms/DivModService/protocols/gen-go/cluster"
	"github.com/huzhao37/gev/example/epms/protocols"
	nc "github.com/huzhao37/gev/example/epms/protocols/gen-go/ncbaseheader"
	t "github.com/huzhao37/gev/example/epms/thrift"
	"io"
	"log"
	"net"
)

func Packet(data []byte) []byte {
	buffer := make([]byte, 16+len(data))
	// 将buffer前面16个字节设置为包长度，大端序
	binary.LittleEndian.PutUint32(buffer, uint32(len(data)))
	copy(buffer[16:], data)
	return buffer
}

func UnPacket(c net.Conn) ([]byte, error) {
	var header = make([]byte, 16)

	_, err := io.ReadFull(c, header)
	if err != nil {
		return nil, err
	}
	_ = binary.LittleEndian.Uint32(header)
	epmsHeader := protocols.BytesToEpmsHeader(header)
	fmt.Printf("%v", epmsHeader)

	bodyByte := make([]byte, epmsHeader.BodyLen)
	_, e := io.ReadFull(c, bodyByte) //读取内容
	if e != nil {
		return nil, e
	}

	return bodyByte, nil
}

func main() {
	prot := &protocols.EpmsProtocol{}
	th := &t.Thrift{}
	conn, e := net.Dial("tcp", ":1833")
	if e != nil {
		log.Fatal(e)
	}
	defer conn.Close()

	for {
		//reader := bufio.NewReader(os.Stdin)
		//fmt.Print("Text to send: ")
		//text, _ := reader.ReadString('\n')
		st := cluster.NewDivModDoDivModArgs()
		st.Arg1 = 19800
		st.Arg2 = 100
		err, argsBuffer := th.GetArgsStructBuffer(st)

		buffer := prot.Packet(nil, protocols.EpmsBodyToBytes(&nc.NcEPMSMsgHeader{MsgType: 2, MsgName: "msg://epms/cluster/DoDivMod",
			SourceId: 1, ProtoName: "", Buffer: argsBuffer, BufLength: int32(len(argsBuffer))}))

		_, err = conn.Write(buffer)
		if err != nil {
			panic(err)
		}

		// listen for reply
		msg, err := UnPacket(conn)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Message from server (len %d) : %s", len(msg), protocols.BytesToEpmsBody(msg))
	}
}
