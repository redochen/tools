package concurrent

import (
	queue "gopkg.in/eapache/queue.v1"
	"sync"
)

//ConcurrentQueue 线程安全的队列类
type ConcurrentQueue struct {
	lock   *sync.Mutex
	cond   *sync.Cond
	queue  *queue.Queue
	closed bool
}

//NewConcurrentQueue 获取新的ConcurrentQueue实例
func NewConcurrentQueue() *ConcurrentQueue {
	q := &ConcurrentQueue{
		lock:   &sync.Mutex{},
		queue:  queue.New(),
		closed: false,
	}

	q.cond = sync.NewCond(q.lock)
	return q
}

//Push 将一项Push到队列中，总是直接返回
func (q *ConcurrentQueue) Push(v interface{}) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if !q.closed {
		q.queue.Add(v)
		q.cond.Signal()
	}
}

//Pop 从队列中Pop一项，如果队列为空则会阻塞
func (q *ConcurrentQueue) Pop() interface{} {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.closed {
		return nil
	}

	if q.queue.Length() <= 0 {
		q.cond.Wait()
	}

	return q.queue.Remove()
}

//TryPop 尝试从队列中Pop一项，如果队列为空则直接返回
func (q *ConcurrentQueue) TryPop() (interface{}, bool) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.queue.Length() > 0 {
		return q.queue.Remove(), true
	}

	return nil, q.closed
}

//Length 获取队列的长度
func (q *ConcurrentQueue) Length() int {
	q.lock.Lock()
	defer q.lock.Unlock()

	return q.queue.Length()
}

//Close 关闭队列
func (q *ConcurrentQueue) Close() {
	q.lock.Lock()
	defer q.lock.Unlock()

	if !q.closed {
		q.closed = true
		q.cond.Signal()
	}
}
