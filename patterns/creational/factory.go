package main

import "fmt"

//接口，使用一个共同的接口创建新的对象
type OrderCreater interface {
	Create()
}

const (
	Inland int = 1 << iota
	International
)

//接口具体实现类

//国内订单
type InlandOrder struct {
}

func (i *InlandOrder) Create() {
	fmt.Println("inland order create")
}

//国际订单
type InternationalOrder struct {
}

func (i *InternationalOrder) Create() {
	fmt.Println("internation order create")
}
//工厂方法
func NewCreateOrderFactory(orderType int) OrderCreater {

	switch orderType {

	case Inland:
		return new(InlandOrder)
	case International:
		return new(InternationalOrder)
	default:
		fmt.Println("unkonw type")
		return nil
	}
}
func main() {

	//根据条件，调用工厂方法
	order := NewCreateOrderFactory(Inland)
	order.Create()

	order2 := NewCreateOrderFactory(International)
	order2.Create()

}
