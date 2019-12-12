package main

import (
	"fmt"
)

type userNode struct {
	id   int
	name string
	next *userNode
}

//法1：在链表尾部插入一个新节点
func (u *userNode) insertUserNodeOfLast(head *userNode, newUserNode *userNode) {
	//思路：
	//1.创建一个临时node
	//2.遍历先找该链表的最后这个节点
	//3.追加到最后一个节点

	temp := head
	for {
		if temp.next == nil { //表示找到最后节点返回
			break
		}
		temp = temp.next //一直遍历到最后节点
	}
	//3.将newUserNode加入到链表的最后
	temp.next = newUserNode
}

//法2：根据id的编号大小顺序插入，比较id大小，注意不用重复添加
func (u *userNode) insertUserNodeByID(head *userNode, newUserNode *userNode) {

	temp := head
	exists := false
	for {
		if temp.next == nil {
			break
		} else if temp.next.id > newUserNode.id {
			//说明插入到当前temp的后面
			break
		} else if temp.next.id == newUserNode.id {
			exists = true
			break
		}
		temp = temp.next
	}
	if exists {
		fmt.Println("exists id", newUserNode)
	} else {
		//下面两行代码，顺序不能变，否则变成了无穷尽的list
		newUserNode.next = temp.next
		temp.next = newUserNode

	}
}

//根据id 删除一个节点
func (u *userNode) delUserNode(head *userNode, id int) {
	temp := head
	exists := false
	for {
		if temp.next == nil {
			break
		} else if temp.next.id == id {
			exists = true
			break
		}
		temp = temp.next
	}
	if exists {
		temp.next = temp.next.next //跳过这个节点，
	} else {
		fmt.Println("no exists id")
	}
}

//根据id 修改节点
func (u *userNode) modifyUserNode(head *userNode, id int, newUserNode *userNode) {
	temp := head
	exists := false

	for {
		if temp.next == nil {
			break
		} else if temp.next.id == id {
			exists = true
			break
		}
		temp = temp.next
	}
	if exists {
		temp.next = newUserNode
	} else {
		fmt.Println("not exists id")
	}
}

//查询所有节点信息
func (u *userNode) listUserNode(head *userNode) {
	temp := head

	if temp.next == nil {
		fmt.Println("empty list")
		return
	}
	for {
		fmt.Printf("%d,%s-->", temp.next.id, temp.next.name)
		temp = temp.next //判断是否是链表的末尾
		if temp.next == nil {
			break
		}
	}

}

func main() {
	//1.先创建一个头结点
	head := &userNode{}

	//2.创建一批新的userNode
	stu1 := &userNode{
		id:   1,
		name: "tom1",
	}
	stu2 := &userNode{
		id:   2,
		name: "tom2",
	}
	stu3 := &userNode{
		id:   3,
		name: "tom3",
	}
	stu4 := &userNode{
		id:   4,
		name: "tom4",
	}

	stu5 := &userNode{
		id:   5,
		name: "tom5",
	}

	//3.加入节点
	head.insertUserNodeOfLast(head, stu1)
	head.insertUserNodeOfLast(head, stu2)

	//查询
	head.listUserNode(head)
	fmt.Println()
	//4.加入结点（第二种方法）
	head.insertUserNodeByID(head, stu4) //id是4
	head.insertUserNodeByID(head, stu3) //id是3
	head.listUserNode(head)
	fmt.Println()
	//5.删除结点
	head.delUserNode(head, 1)
	//显示链表
	head.listUserNode(head)
	fmt.Println()
	//6.修改结点
	head.modifyUserNode(head, 3, stu5)
	head.listUserNode(head)
	fmt.Println()
}
