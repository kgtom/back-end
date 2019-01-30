/*
例子：订单待支付、支付两种状态，状态改变时，执行相应的操作。

*/

package main

import "fmt"

//接口
type OrderStater interface {
	Handle()string
}


//接口具体的实现类

//待支付订单
type UnPayOrder struct {
}

func (i *UnPayOrder) Handle() string {
	return "unPayOrder handle"
}

// 已支付订单
type PayOrder struct {
}

func (i *PayOrder) Handle() string {
	return "PayOrder  handle"
}

// 创建Context类
type OrderContext struct {

	OrderStater

}

func (c *OrderContext)Handle()string  {
	return c.OrderStater.Handle()
}
func (c *OrderContext)SetState(state OrderStater)  {
	c.OrderStater=state
}

func main() {

	//调用：使用contex中改变state,执行不同行为(方法)
	ctx:=OrderContext{}
	ctx.SetState(new(UnPayOrder))
	ret:=ctx.Handle()
	fmt.Println(ret)

	ctx.SetState(new(PayOrder))
    ret=ctx.Handle()
    fmt.Println(ret)
}
