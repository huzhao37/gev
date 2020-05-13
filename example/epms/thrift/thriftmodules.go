/**
 * @Author: hiram
 * @Date: 2020/5/13 10:44
 */
package thrift

import (
	"bytes"
	"context"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/huzhao37/gev/example/epms/DivModService/protocols/gen-go/cluster"
	nc "github.com/huzhao37/gev/example/epms/protocols/gen-go/ncbaseheader"
	"github.com/huzhao37/gev/log"
	"reflect"
)

//!#obselete
type Thrift struct {
	processor       *thrift.TMultiplexedProcessor
	protocolFactory *thrift.TCompactProtocolFactory
}

//!#obselete
type ThriftHandlers struct {
	ServiceName string
	Processor   thrift.TProcessor
}

func (t *Thrift) setArgsStructValue(st thrift.TStruct, buffer []byte) error {
	iprot := t.protocolFactory.GetProtocol(&thrift.TMemoryBuffer{
		Buffer: bytes.NewBuffer(buffer),
	})

	if err := st.Read(iprot); err != nil {
		iprot.ReadMessageEnd()
		log.Error("【GetArgsStructValue】%s", err)
		return err
	}
	iprot.ReadMessageEnd()
	return nil
}

func (t *Thrift) GetArgsStructBuffer(st thrift.TStruct) (error, []byte) {
	bufferProt := &thrift.TMemoryBuffer{}
	iprot := t.protocolFactory.GetProtocol(bufferProt)

	if err := st.Write(iprot); err != nil {
		return err, nil
	}
	if err := iprot.WriteMessageEnd(); err != nil {
		return err, nil
	}
	return nil, bufferProt.Buffer.Bytes()
}

func (t *Thrift) GetResultStructBuffer(st thrift.TStruct) (error, []byte) {
	bufferProt := &thrift.TMemoryBuffer{}
	oprot := t.protocolFactory.GetProtocol(bufferProt)

	var err error
	if err = st.Write(oprot); err != nil {
		return err, nil
	}
	if err = oprot.WriteMessageEnd(); err != nil {
		return err, nil
	}
	return nil, bufferProt.Buffer.Bytes()
}

//get args struct by name
func (t *Thrift) getArgsStruct(argsName string) thrift.TStruct {
	switch argsName {
	//todo
	}
	//test
	return cluster.NewDivModDoDivModArgs()
}

//get result struct by name
func (t *Thrift) getResultStruct(resultName string) thrift.TStruct {
	switch resultName {
	//todo
	}
	//test
	return cluster.NewDivModDoDivModResult()
}

//get args &r reply
func (t *Thrift) GetArgsAndReply(epmsBody *nc.NcEPMSMsgHeader) (thrift.TStruct, thrift.TStruct) {
	//todo
	//根据msgtype,msgname ,protoName 获取相应的入参和出参
	args := t.getArgsStruct("")
	err := t.setArgsStructValue(args, epmsBody.Buffer)
	if err != nil {
		log.Error("【GetArgsAndReply】%s", err)
	}
	return args, t.getResultStruct("")
}

//!#obselete
// Call calls a service
func (t *Thrift) Call(servicePath, serviceMethod string, args interface{}, reply interface{}) (err error) {
	defer func() {
		if e := recover(); e != nil {
			var ok bool
			if err, ok = e.(error); ok {
				err = fmt.Errorf("failed to call %s.%s because of %v", servicePath, serviceMethod, err)
			}
		}
	}()

	var mv = &reflect.Value{}
	if mv == nil {
		//v := reflect.ValueOf(service)
		//t := v.MethodByName(serviceMethod)
		//if t == (reflect.Value{}) {
		//return fmt.Errorf("method %s.%s not found", servicePath, serviceMethod)
		//}
		//mv = &t
	}

	argv := reflect.ValueOf(args)
	replyv := reflect.ValueOf(reply)

	err = nil
	returnValues := mv.Call([]reflect.Value{argv, replyv})
	errInter := returnValues[0].Interface()
	if errInter != nil {
		err = errInter.(error)
	}

	return err
}

//register thrift service
func (t *Thrift) RegisterThriftProcessor(handlers []ThriftHandlers) {

	t.processor = thrift.NewTMultiplexedProcessor()
	// 给每个service起一个名字
	if len(handlers) > 0 {
		for _, handler := range handlers {
			t.processor.RegisterProcessor(handler.ServiceName, handler.Processor)
		}
	}
	t.protocolFactory = thrift.NewTCompactProtocolFactory()
}

func (t *Thrift) Process(ctx context.Context, data []byte) {
	transport := &thrift.TMemoryBuffer{
		Buffer: bytes.NewBuffer(data),
	}
	ipro := t.protocolFactory.GetProtocol(transport)
	t.processor.Process(ctx, ipro, ipro)
}
