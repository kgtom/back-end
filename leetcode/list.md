### 链表

~~~ go
/*

双向链表
 * 单向链表只可向一个方向遍历，一般查找一个结点的时候需要从第一个结点开始每次访问下一个结点，一直访问到需要的位置。
 * 双向链表的每个结点都有指向的前一个结点和后一个节点，既有链表表头又有表尾，即可从链头向链尾遍历，又可从链尾向链头遍历。
*/

package list

type List struct {
	Length int
	Head   *Node
	Tail   *Node
}

func NewList() *List {

	l := new(List)
	l.Length = 0
	return l
}

type Node struct {
	Value interface{}
	Prev  *Node
	Next  *Node
}

func (l *List) Len() int {
	return l.Length
}

func (l *List) IsEmpty() bool {
	return l.Length == 0
}
func (l *List) Prepend(value interface{}) {
	node := NewNode(value)
	if l.Len() == 0 {
		l.Head = node
		l.Tail = l.Head
	} else {
		formerHead := l.Head
		formerHead.Prev = node

		node.Next = formerHead
		l.Head = node
	}

	l.Length++
}

func (l *List) Append(value interface{}) {
	node := NewNode(value)

	if l.Len() == 0 {
		l.Head = node
		l.Tail = l.Head
	} else {
		formerTail := l.Tail
		formerTail.Next = node

		node.Prev = formerTail
		l.Tail = node
	}

	l.Length++
}
func NewNode(v interface{}) *Node {
	return &Node{Value: v}
}


~~~
