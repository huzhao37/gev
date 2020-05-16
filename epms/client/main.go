package main

import (
	"encoding/binary"
	"github.com/huzhao37/gev/epms/DivModService/protocols/gen-go/cluster"
	"github.com/huzhao37/gev/epms/protocols"
	nc "github.com/huzhao37/gev/epms/protocols/gen-go/ncbaseheader"
	t "github.com/huzhao37/gev/epms/thrift"
	"github.com/huzhao37/gev/log"
	"net"
	"time"
)

func Packet(data []byte) []byte {
	buffer := make([]byte, 16+len(data))
	// 将buffer前面16个字节设置为包长度，大端序
	binary.LittleEndian.PutUint32(buffer, uint32(len(data)))
	copy(buffer[16:], data)
	return buffer
}

func main() {
	prot := &protocols.EpmsProtocol{}
	th := &t.Thrift{}
	conn, e := net.Dial("tcp", ":1833")
	if e != nil {
		log.Fatal(e)
	}
	//defer conn.Close()

	for {
		//reader := bufio.NewReader(os.Stdin)
		//fmt.Print("Text to send: ")
		//text, _ := reader.ReadString('\n')
		st := cluster.NewDivModDoDivModArgs()
		st.Arg1 = 19800
		st.Arg2 = 100
		err, argsBuffer := th.GetStructBuffer(st)
		epmsBody := nc.NcEPMSMsgHeader{MsgType: 2, MsgName: "msg://epms/cluster/DoDivMod",
			SourceId: 1, ProtoName: "", Buffer: argsBuffer, BufLength: int32(len(argsBuffer))}

		err, data := th.GetStructBuffer(&epmsBody)

		buffer := prot.Packet(nil, data)

		_, err = conn.Write(buffer)
		if err != nil {
			log.Error("[Client-UnPacket]%v", err)
		}
		time.Sleep(1 * time.Second)
		// listen for reply
		msg, err := protocols.ClientUnPacket(conn)
		if err != nil {
			log.Error("[Client-UnPacket]%v", err)
		}
		reply := &nc.NcEPMSMsgHeader{}
		th.GetStructValue(reply, msg)
		res := th.GetStruct("DivModDoDivModResult")
		th.GetStructValue(res, reply.Buffer)
		log.Info("Message from server (len %d) : %s", len(msg), res)
		time.Sleep(5 * time.Second)
	}
}
