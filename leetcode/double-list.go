package main

import (
	"errors"
	"fmt"
	"sync"
)

//单向链表：包括数据+指向下一个节点指针，只能单向读取，
//双向链表：包括 头部指针和尾部指针，
//单双比较：单：在存储空间上只保存下一个节点，所以比双向链表节省空间；双：在处理时间比单向链表快，因为单需要遍历O(n),双只要找到head->tail

//链表优缺点：添加和删除元素快O(1)，内存利用率高，缺点：查找元素时需要遍历复杂度O(N)
//链表适用于：频繁添加和删除的操作,比如redis中list的push、pop，hash的hget,zset的zscore；经常访问数组比链表更换

//链表与数组比较：数组固定长度，分配栈上，查找下标O(1)，新增、删除需要遍历O(N);链表动态长度，分配堆上，查找O(N),新增、删除O(1)

//链表，类似于火车若干车厢，收尾相连。车厢就是每一个节点，
type List struct {
	Length int
	Head   *Node
	Tail   *Node
	lock   *sync.RWMutex
}

func NewList() *List {
	l := new(List)
	l.lock = new(sync.RWMutex)
	l.Length = 0
	return l
}

type Node struct {
	Value interface{}
	Prev  *Node
	Next  *Node
}

func NewNode(value interface{}) *Node {
	return &Node{Value: value}
}

func (l *List) Len() int {
	return l.Length
}

func (l *List) IsEmpty() bool {
	return l.Length == 0
}

//头插法
func (l *List) InsertAtHead(value interface{}) {

	if value == nil {
		return
	}
	l.lock.Lock()
	defer l.lock.Unlock()
	node := NewNode(value)
	if l.Len() == 0 {
		l.Head = node
		l.Tail = l.Head
	} else {
		tempHead := l.Head   //上次head
		tempHead.Prev = node //上次头的头部插入node

		node.Next = tempHead //新插入的头的next应该是之前的头
		l.Head = node
	}

	l.Length++
}

//尾插法
func (l *List) AppendAtTail(value interface{}) {

	if value == nil {
		return
	}
	l.lock.Lock()
	defer l.lock.Unlock()

	node := NewNode(value)

	if l.Len() == 0 {
		l.Head = node
		l.Tail = l.Head
	} else {
		tempTail := l.Tail
		tempTail.Next = node

		node.Prev = tempTail
		l.Tail = node
	}

	l.Length++
}

func (l *List) Remove(node interface{}) error {
	if l.Len() == 0 {
		return errors.New("Empty list")
	}

	l.lock.Lock()
	defer l.lock.Unlock()
	if l.Head.Value.(*Node).Value == node {
		l.Head = l.Head.Next
		l.Length--
		return nil
	}

	found := 0
	for n := l.Head; n != nil; n = n.Next {

		if n.Value.(*Node).Value == node && found == 0 {
			n.Next, n.Prev.Next = n.Prev, n.Next
			l.Length--
			found++
		}
	}

	if found == 0 {
		return errors.New("node not found")
	}

	return nil
}

func (l *List) GetByIdx(idx int) (*Node, error) {
	if idx >= l.Len() {
		return nil, errors.New("idx out of range")
	}

	l.lock.RLock()
	defer l.lock.RUnlock()
	node := l.Head
	for i := 0; i < idx; i++ {
		node = node.Next
	}

	return node, nil
}

func (l *List) Clear() {
	l.Length = 0
	l.Head = nil
	l.Tail = nil
}
func (l *List) ShowList() {
	if l == nil || l.Length == 0 {
		fmt.Println("list is nil or empty")
		return
	}
	l.lock.RLock()
	defer l.lock.RUnlock()
	fmt.Printf("list lenght is %d \n", l.Length)
	currNode := l.Head
	for currNode != nil {
		fmt.Printf("curr data value is %v", currNode.Value)
		if currNode.Prev != nil {
			fmt.Printf("prev value is %v", currNode.Prev.Value)

		}
		if currNode.Next != nil {
			fmt.Printf("next value is %v\n", currNode.Next.Value)

		}
		fmt.Printf("\n")
		currNode = currNode.Next
	}
}

func main() {
	fmt.Println(" start list ....")
	l := NewList()

	l.InsertAtHead(NewNode(1))
	l.InsertAtHead(NewNode(2))
	l.InsertAtHead(NewNode(3))

	l.ShowList()

	fmt.Println()

	// Test GetByIdx
	v, err := l.GetByIdx(3)
	if err != nil {

		fmt.Println("GetByIdx err:", err)
		return
	}

	fmt.Println("index 3 data value is:", (v.Value))

	fmt.Println()
	return
	// Test AppendAtTail

	l.AppendAtTail(NewNode(4))
	l.AppendAtTail(NewNode(5))
	l.AppendAtTail(NewNode(6))

	l.ShowList()

	//Test clear
	l.Remove(3)

	l.Remove(5)
	l.ShowList()

	//Test clear
	//l.Clear()
	//l.ShowList()
	fmt.Println(" end list ....")
}
