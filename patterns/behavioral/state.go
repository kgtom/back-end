/*
例子：订单待审核、待支付、支付三种状态，flag改变时，执行相应类的操作。

*/

package main

import "fmt"

//接口
type OrderStater interface {
	Handle(c *OrderContext)

}


//接口具体的实现类

//待审核的订单
type UnAuditOrder struct {

}
func (i *UnAuditOrder)Handle(c *OrderContext)  {

	if c.flag==1{
		fmt.Println("UnAuditOrder handle")
	}else {

		c.SetState(new(UnPayOrder))
		c.Do()
	}

}

//待支付的订单
type UnPayOrder struct {
}


func (i *UnPayOrder) Handle(c *OrderContext)  {

	if c.flag==2{
		fmt.Println( "UnPayOrder handle")
	}else {

		c.SetState(new(PayOrder))
		c.Do()
	}
}


// 已支付的订单
type PayOrder struct {
}

func (i *PayOrder) Handle(c *OrderContext)  {

	if c.flag==3{
		fmt.Println("PayOrder  handle")
	}else {
		fmt.Println("PayOrder end")
	}
}


// 创建Context类
type OrderContext struct {
    flag int
	OrderStater

}

func NewOrderContex() OrderContext  {
	return OrderContext{flag:1,OrderStater:new(UnAuditOrder)}
}

func (c *OrderContext)Do()  {
	 c.OrderStater.Handle(c)
}
func (c *OrderContext)SetState(state OrderStater)  {

	c.OrderStater=state

}
func (c *OrderContext)Setflag(flag int )  {

	c.flag=flag
}



func main() {

	//调用：使用contex中改变 flag,执行不同行为(方法)
	ctx:=NewOrderContex()

	ctx.Do()
	ctx.Setflag(2)
	ctx.Do()
	ctx.Setflag(3)
	ctx.Do()
	
}

