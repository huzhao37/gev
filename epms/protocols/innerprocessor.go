/**
 * @Author: hiram
 * @Date: 2020/5/9 14:44
 */
package protocols

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"sync"
)

const HANDLER_SEPARATOR = ":"

// ReqMetaDataKey is used to set metatdata in context of requests.
var ReqMetaDataKey = ContextKey("__req_metadata")

// Inprocessor is a in-process client for test.
var Inprocessor = &inprocessClient{
	services: make(map[string]interface{}),
	methods:  make(map[string]*reflect.Value),
}

// ContextKey defines key type in context.
type ContextKey string

// ServiceError is an error from server.
type ServiceError string

func (e ServiceError) Error() string {
	return string(e)
}

// Call represents an active RPC.
type Call struct {
	ServicePath   string            // The name of the service and method to call.
	ServiceMethod string            // The name of the service and method to call.
	Metadata      map[string]string //metadata
	ResMetadata   map[string]string
	Args          interface{} // The argument to the function (*struct).
	Reply         interface{} // The reply from the function (*struct).
	Error         error       // After completion, the error status.
	Done          chan *Call  // Strobes when call is complete.
	Raw           bool        // raw message or not
}

// inprocessClient is a in-process client that call services in process not via TCP/UDP.
// Notice this client is only used for test.
type inprocessClient struct {
	services map[string]interface{}
	sync.RWMutex
	methods map[string]*reflect.Value
	mmu     sync.RWMutex
	//ServerMessageChan chan<- *Message
}

func (call *Call) done() {
	select {
	case call.Done <- call:
		// ok
	default:
		log.Printf("rpc: discarding Call reply due to insufficient Done chan capacity")

	}
}

func (client *inprocessClient) Register(name string, rcvr interface{}, metadata string) (err error) {
	client.Lock()
	client.services[name] = rcvr
	client.Unlock()
	return
}

func (client *inprocessClient) Unregister(name string) error {
	client.Lock()
	delete(client.services, name)
	client.Unlock()
	return nil
}

// Go calls is not async. It still use sync to call.
func (client *inprocessClient) Go(ctx context.Context, servicePath, serviceMethod string, args interface{}, reply interface{}, done chan *Call) *Call {
	call := new(Call)
	call.ServicePath = servicePath
	call.ServiceMethod = serviceMethod
	meta := ctx.Value(ReqMetaDataKey)
	if meta != nil { //copy meta in context to meta in requests
		call.Metadata = meta.(map[string]string)
	}
	call.Args = args
	call.Reply = reply
	if done == nil {
		done = make(chan *Call, 10) // buffered.
	} else {
		// If caller passes done != nil, it must arrange that
		// done has enough buffer for the number of simultaneous
		// RPCs that will be using that channel. If the channel
		// is totally unbuffered, it's best not to run at all.
		if cap(done) == 0 {
			log.Panic("rpc: done channel is unbuffered")
		}
	}
	call.Done = done
	err := client.Call(ctx, servicePath, serviceMethod, args, reply)
	if err != nil {
		call.Error = ServiceError(err.Error())

	}
	call.done()
	return call
}

// Call calls a service synchronously.
func (client *inprocessClient) Call(ctx context.Context, servicePath, serviceMethod string, args interface{}, reply interface{}) (err error) {
	defer func() {
		if e := recover(); e != nil {
			var ok bool
			if err, ok = e.(error); ok {
				err = fmt.Errorf("failed to call %s.%s because of %v", servicePath, serviceMethod, err)
			}
		}
	}()

	client.RLock()
	service := client.services[servicePath]
	client.RUnlock()
	if service == nil {
		return fmt.Errorf("service %s not found", servicePath)
	}

	key := servicePath + "." + serviceMethod
	client.mmu.RLock()
	var mv = &reflect.Value{}
	mv = client.methods[key]
	client.mmu.RUnlock()

	if mv == nil {
		client.mmu.Lock()
		mv = client.methods[key]
		if mv == nil {
			v := reflect.ValueOf(service)
			t := v.MethodByName(serviceMethod)
			if t == (reflect.Value{}) {
				client.mmu.Unlock()
				return fmt.Errorf("method %s.%s not found", servicePath, serviceMethod)
			}
			mv = &t
		}
		client.mmu.Unlock()
	}

	argv := reflect.ValueOf(args)
	replyv := reflect.ValueOf(reply)

	err = nil
	returnValues := mv.Call([]reflect.Value{reflect.ValueOf(ctx), argv, replyv})
	errInter := returnValues[0].Interface()
	if errInter != nil {
		err = errInter.(error)
	}

	return err
}
