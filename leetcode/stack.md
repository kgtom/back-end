~~~ go

package stack

import "sync"

type Stack struct {
	stack []interface{}
	len   int
	lock  sync.Mutex
}

func NewStack() *Stack {

	return &Stack{
		stack: make([]interface{}, 0),
		len:   0,
	}
}

func (s *Stack) Len() int {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.len
}
func (s *Stack) IsEmpty() bool {

	s.lock.Lock()
	defer s.lock.Unlock()
	return s.len == 0
}

func (s *Stack) Pop() (el interface{}) {

	s.lock.Lock()
	defer s.lock.Unlock()
	el, s.stack = s.stack[0], s.stack[1:]
	s.len--
	return el
}

func (s *Stack) Push(el interface{}) {

	s.lock.Lock()
	defer s.lock.Unlock()
	new := make([]interface{}, 1)
	new[0] = el
	s.stack = append(new, s.stack...) //参数为两个slice时，不要忘记...

	s.len++
}
func (s *Stack) Peek() interface{} {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.stack[0]
}
~~~
