/**
 * @Author: hiram
 * @Date: 2020/5/12 11:28
 */
package queue

import (
	"fmt"
	"github.com/leandro-lugaresi/hub"
	"sync"
	"testing"
)

type List struct {
	name   string
	hub    *hub.Hub
	cap    int
	queues sync.Map
}
type EventHandler func(string, hub.Message)

func CreateNew(name string, cap int) *List {
	return &List{name: name, hub: hub.New(), cap: cap}
}

//@key==topic:addr
func (l *List) Publish(key string, data []byte) bool {
	l.hub.Publish(hub.Message{
		Name:   key,
		Fields: hub.Fields{"data": data},
	})
	return true
}

//sub
func (l *List) Subscribe(key string, handler EventHandler) hub.Subscription {
	sub := l.hub.Subscribe(l.cap, key)
	go func(s hub.Subscription) {
		for msg := range s.Receiver {
			handler(key, msg)
			fmt.Printf("receive msg with topic %s and data %d\n", msg.Name, msg.Fields["data"])
		}
	}(sub)
	return sub
}

//sub
func (l *List) Unsubscribe(sub hub.Subscription) {
	l.hub.Unsubscribe(sub)
}

//close
func (l *List) Close() {
	l.hub.Close()
}

func TestList(t *testing.T) {

}
