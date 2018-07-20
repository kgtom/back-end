~~~go
//测试实例
fmt.Println("queue....")
	q := queue.NewQueue()
	q.Push(5)
	q.Push(6)
	q.Push(7)
	fmt.Println(q)
	b := q.Pop()
	fmt.Println(b)
	fmt.Println(q)
	one := q.Peek()
~~~

~~~ go
package queue

import "sync"

type Queue struct {
	queue []interface{}
	len   int
	lock  *sync.Mutex
}

func NewQueue() *Queue {
	queue := &Queue{}
	queue.queue = make([]interface{}, 0)
	queue.len = 0
	queue.lock = new(sync.Mutex)
	return queue
}

func (q *Queue) Len() int {

	return q.len
}

func (q *Queue) IsEmpty() bool {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.len == 0
}

func (q *Queue) Push(el interface{}) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.queue = append(q.queue, el)
	q.len++

}

func (q *Queue) Pop() (el interface{}) {

	q.lock.Lock()
	defer q.lock.Unlock()
	el, q.queue = q.queue[0], q.queue[1:]
	q.len--
	return el
}
func (q *Queue) Peek() interface{} {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.queue[0]
}


~~~
