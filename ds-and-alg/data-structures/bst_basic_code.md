## 本届大纲
* [一、二叉树遍历](#1)
* [二、数组转换树](#2)
* [三、二叉查找树](#3)

## <span id="1"> 一、二叉树遍历</span>

~~~go
package main

import (
	"fmt"
	"sync"
)

func main() {

	//二叉树 深度遍历(前、中、后序遍历)及广度优先遍历(层序遍历)

	tree := &TreeNode{
		Val: 1,
		Left: &TreeNode{
			Val: 2,
			Left: &TreeNode{
				Val: 4,
			},
			Right: &TreeNode{
				Val: 5,
			},
		},
		Right: &TreeNode{
			Val: 3,
			Left: &TreeNode{
				Val: 6,
			},
			Right: &TreeNode{
				Val: 7,
			},
		},
	}

	fmt.Println("start")
	//前序 中左右 ret: [1 2 4 5 3 6 7]
	retPreIterative := preOrderIterative(tree)
	fmt.Println("retPreIterative:", retPreIterative)

	retPre := preOrderRecursive(tree)
	fmt.Println("preOrderRecursive", retPre)

	fmt.Println()
	//中序 迭代版 左中右 ret: [4 2 5 1 6 3 7]
	retInIterative := inOrderIterative(tree)
	fmt.Println("retInIterative:", retInIterative)
	//中序 递归版
	retIn := inOrderRecursive(tree)
	fmt.Println("retIn-recursive:", retIn)

	fmt.Println()

	//后序迭代版 左右中 ret: [1 2 4 5 3 6 7]
	retPostIterative := postOrderIterative(tree)
	fmt.Println("retPostIterative:", retPostIterative)
	//后序 递归版
	retPost := postOrderRecursive(tree)
	fmt.Println("retPost-recursive:", retPost)

	//层序(迭代版) ret: [1 2 4 5 3 6 7]
	retLevel := levelIterative(tree)
	fmt.Println("retLevel:", retLevel)

	//层序(递归版) ret: [1 2 3 4 5 6 7]
	retLevel2 := levelRecursive(tree)
	fmt.Println("levelRecursive:", retLevel2)
	fmt.Println("end")

}

var ret [][]int

func levelRecursive(root *TreeNode) [][]int {
	ret = make([][]int, 0)
	if root == nil {
		return ret
	}
	dfsHandler(root, 0)
	return ret
}

func dfsHandler(node *TreeNode, level int) {
	if node == nil {
		return
	}
	if len(ret) < level+1 {
		ret = append(ret, make([]int, 0))
	}
	//ret = append(ret, node.Val)
	ret[level] = append(ret[level], node.Val)

	dfsHandler(node.Left, level+1)
	dfsHandler(node.Right, level+1)
}

// TreeNode
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

type Stack struct {
	lock sync.Mutex
	node []*TreeNode
}

func NewStack() *Stack {

	return &Stack{lock: sync.Mutex{}, node: []*TreeNode{}}
}

func (s *Stack) Push(node *TreeNode) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.node = append(s.node, node)
}

func (s *Stack) Pop() *TreeNode {
	s.lock.Lock()

	defer s.lock.Unlock()
	n := len(s.node)
	if n == 0 {
		return nil
	}
	ret := s.node[n-1]
	s.node = s.node[:(n - 1)]

	return ret
}

type Queue struct {
	queue []*TreeNode
	len   int
	lock  *sync.Mutex
}

func NewQueue() *Queue {
	queue := &Queue{}
	queue.queue = make([]*TreeNode, 0)
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

func (q *Queue) Push(el *TreeNode) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.queue = append(q.queue, el)
	q.len++

}

func (q *Queue) Pop() (el *TreeNode) {

	if q.IsEmpty() {
		return &TreeNode{}
	}
	q.lock.Lock()
	defer q.lock.Unlock()
	el, q.queue = q.queue[0], q.queue[1:]
	q.len--
	return el
}

//前序 迭代版：(中左右)时间复杂度O(n) 空间复杂度(n)
func preOrderIterative(tree *TreeNode) []int {

	s := NewStack()

	//将tree,push 到栈里面
	s.Push(tree)
	var ret []int

	//使用stack 栈，注意其特性，先进后出。

	for len(s.node) != 0 {

		//获取当前节点,首次实际就是根节点值
		curr := s.Pop()
		ret = append(ret, curr.Val)

		//如果当前节点存在右节点，则入栈
		//如果当前节点存在左节点，则入栈
		//先 push right节点，后 push left,是因为栈的特性：先进后出，又因为前序遍历，需要先让左节点出来，后右节点出来。
		if curr.Right != nil {
			s.Push(curr.Right)
		}
		if curr.Left != nil {
			s.Push(curr.Left)
		}

	}
	//fmt.Println("len:", len(s.node)) //0
	return ret
}

//前序 递归版
func preOrderRecursive(root *TreeNode) []int {
	ret := []int{}
	if root == nil {
		return ret
	}
	//fmt.Println(root.Val)
	ret = append(ret, root.Val)
	leftRet := preOrderRecursive(root.Left)
	ret = append(ret, leftRet...)
	ret = append(ret, preOrderRecursive(root.Right)...)

	return ret
}

//  中序 迭代版 ：(左中右)时间复杂度O(n) 空间复杂度(n)
func inOrderIterative(root *TreeNode) []int {
	var ret []int

	if root == nil {
		return ret
	}

	stack := NewStack()
	curr := root

	for {
		//现将root入栈，迭代找到左节点的起点，即左节点为nil的节点。
		for curr != nil {
			stack.Push(curr)
			curr = curr.Left
		}

		//因为上一步左节点最后一个入栈，则第一个出栈，符合中序遍历，先左节点开始，然后中节点，最后是右节点
		node := stack.Pop()
		if node == nil {
			break
		}

		ret = append(ret, node.Val)

		// 再将右节点入栈，入栈是 右节点部分也要按照 左-中-右
		if node.Right != nil {
			curr = node.Right
		}
	}

	return ret
}

//中序 递归版
func inOrderRecursive(root *TreeNode) []int {
	var ret []int

	if root == nil {
		return ret
	}

	ret = append(ret, inOrderRecursive(root.Left)...)
	ret = append(ret, root.Val)
	ret = append(ret, inOrderRecursive(root.Right)...)

	return ret
}

//后序 迭代版 左右中
func postOrderIterative(root *TreeNode) []int {
	var ret []int

	if root == nil {
		return ret
	}
	sTemp := NewStack()
	s := NewStack()

	sTemp.Push(root)

	for {

		node := sTemp.Pop()

		if node == nil {
			break
		}
		s.Push(node)

		if node.Left != nil {

			sTemp.Push(node.Left)

		}
		if node.Right != nil {
			sTemp.Push(node.Right)
		}

	}
	for {
		node := s.Pop()
		if node == nil {
			break
		}
		ret = append(ret, node.Val)
	}

	return ret
}

//后序 递归版
func postOrderRecursive(root *TreeNode) []int {
	var ret []int
	if root == nil {
		return ret
	}

	ret = append(ret, postOrderRecursive(root.Left)...)

	ret = append(ret, postOrderRecursive(root.Right)...)
	ret = append(ret, root.Val)
	//fmt.Println(root.Val)

	return ret
}

//层序遍历 迭代版，实际是 Breadth First Search,广度优先算法，从顶而下，每一层 从左到右。
func levelIterative(root *TreeNode) []int {

	var ret []int

	//s := NewStack()
	s := NewQueue() //使用队列，
	s.Push(root)
	for {

		curr := s.Pop()
		if curr == nil || curr.Val == 0 {
			break
		}
		ret = append(ret, curr.Val)
		//先左后右，使用队列：所以现在push left 再push right

		if curr.Left != nil {
			s.Push(curr.Left)
		}
		if curr.Right != nil {
			s.Push(curr.Right)
		}
	}
	return ret
}

~~~

结果：
~~~

start
retPreIterative: [1 2 4 5 3 6 7]
preOrderRecursive [1 2 4 5 3 6 7]

retInIterative: [4 2 5 1 6 3 7]
retIn-recursive: [4 2 5 1 6 3 7]

retPostIterative: [4 5 2 6 7 3 1]
retPost-recursive: [4 5 2 6 7 3 1]
retLevel: [1 2 3 4 5 6 7]
levelRecursive: [[1] [2 3] [4 5 6 7]]
end
~~~

## <span id="2"> 二.数组转换树</span>
~~~go

func Int2TreeNode(nums []int) *TreeNode {
	n := len(nums)
	if n == 0 {
		return nil
	}

	root := &TreeNode{
		Val: nums[0],
	}

	tempNums := make([]*TreeNode, 1)
	tempNums[0] = root

	i := 1
	for i < n {
		//每次取第一个作为根，然后再去找左右
		node := tempNums[0]
		tempNums = tempNums[1:]

		if i < n && nums[i] > 0 {
			node.Left = &TreeNode{Val: nums[i]}
			tempNums = append(tempNums, node.Left)
		}
		i++

		if i < n && nums[i] > 0 {
			node.Right = &TreeNode{Val: nums[i]}
			tempNums = append(tempNums, node.Right)
		}
		i++
	}

	return root
}
~~~

## <span id="3">三、二叉查找树(新增、修改、删除、遍历、查找)</span>
~~~go
package main

import (
	"fmt"
)

func main() {

	node := NewTreeNode()
	node.Insert(90)
	node.Insert(50)
	node.Insert(150)
	node.Insert(20)
	node.Insert(5)
	node.Insert(25)
	fmt.Println("node:", node)

	isContains := node.Search(63)
	fmt.Println("isContains:", isContains)
	retOrder := node.InOrderTraversal()
	fmt.Println("retOrder:", retOrder)
	retMin := node.FindMin()
	fmt.Println("retMin:", retMin)
	retMax := node.FindMax()
	fmt.Println("retMax:", retMax)
	retDel := node.deteleNode(50) //deleteNode2(node, 50)
	fmt.Println("retDel:", retDel)
	fmt.Println("end")

}

// TreeNode
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func NewTreeNode() *TreeNode {
	return &TreeNode{}
}

//新增
func (node *TreeNode) Insert(v int) {

	if node.Val == 0 {

		node.Val = v
	} else if v < node.Val {
		//左
		if node.Left != nil {
			node.Left.Insert(v)
		} else {
			node.Left = &TreeNode{Val: v}
		}

	} else {
		//右
		if node.Right != nil {
			node.Right.Insert(v)
		} else {
			node.Right = &TreeNode{Val: v}
		}
	}

}

//查找
func (node *TreeNode) Search(v int) bool {

	if node == nil {
		return false
	}
	if node.Val == v {
		return true
	} else if node.Val > v {

		node.Left.Search(v)
	} else {
		node.Right.Search(v)
	}

	return false
}

//遍历 递归版(前、中、后序遍历),本例拿中序(左中右)
func (node *TreeNode) InOrderTraversal() []int {

	var ret []int

	if node == nil {
		return ret
	}
	ret = append(ret, node.Left.InOrderTraversal()...)
	ret = append(ret, node.Val)
	//fmt.Println("node:", node.Val)
	ret = append(ret, node.Right.InOrderTraversal()...)
	return ret
}

// 迭代 如果有左分支，则最小的一定是 左分支最后一个左侧节点
func (node *TreeNode) FindMin() int {

	for {
		if node.Left != nil {
			node = node.Left
		} else {
			return node.Val

		}
	}
	return node.Val

}

//递归 查找最大的，一定在右节点
func (node *TreeNode) FindMax() int {
	var ret int
	if node.Right != nil {
		ret = node.Right.FindMax()
	} else {
		return node.Val
	}
	return ret
}

//删除
func (node *TreeNode) deteleNode(v int) *TreeNode {

	if node == nil {
		return node
	}

	if node.Val == v {
		//找到了
		return mergeNode(node)

	}

	if node.Val > v {
		node.Left = node.Left.deteleNode(v)

	} else {

		node.Right = node.Right.deteleNode(v)
	}

	return node
}

func mergeNode(root *TreeNode) *TreeNode {

	right := root.Right
	left := root.Left

	//1.如果删除的节点没有右节点，则让它左节点代替删除的节点
	if right == nil {
		return left
	}
	node := right
	//2.如果删除的节点有右节点，且右节点下面还有左节点，则让最末层左节点替代删除的节点
	for node.Left != nil {
		node = node.Left
	}
	//3.如果删除的节点有右节点，且右节点下面没有做左节点，则让右节点替代删除节点，即右节点的左节点为删除节点的左节点
	node.Left = left
	return root.Right
}

~~~
