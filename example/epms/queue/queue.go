/**
 * @Author: hiram
 * @Date: 2020/5/9 16:26
 */
package queue

/// Queue 队列信息
type Queue struct {
	list *SingleList
	Name string
}

// Init 队列初始化 param@name:addr  @id:sys/biz  (每个连接同时在外部初始化4个队列sys/biz-read/write)
func (q *Queue) Init(name string, id int) {
	q.Name = name
	q.list = new(SingleList)
	q.list.Init(id)
}

// Size 获取队列长度
func (q *Queue) Size() uint {
	return q.list.Size
}

// Enqueue 进入队列
func (q *Queue) Enqueue(data interface{}) bool {
	return q.list.Append(&SingleNode{Data: data})
}

// Dequeue 出列
func (q *Queue) Dequeue() interface{} {
	node := q.list.Get(0)
	if node == nil {
		return nil
	}
	q.list.Delete(0)
	return node.Data
}

// Peek 查看队头信息
func (q *Queue) Peek() interface{} {
	node := q.list.Get(0)
	if node == nil {
		return nil
	}
	return node.Data
}
