/*
例子：第三方通过代理类 OrderProxy 访问订单的信息，而不是直接访问 OrderInfo

代理模式：访问控制行为，针对单个类的代理，而外观模式则提供了一组更高层次接口，封装了若干个类。
 */

package main

import "fmt"

//1.创建一个接口
type OrderInfoer interface {
	GetOrder(int) interface{}
}

//2.创建接口的实现类

type OrderInfo struct {
	ID   int64
	Name string
}

func (order OrderInfo) GetOrder(id int) interface{} {
	//todo form db getOrderById
	return order
}

//3.创建实现类的代理类

type OrderProxy struct {
	Order OrderInfoer
}

func (proxy OrderProxy) GetOrder(id int) interface{} {

	return proxy.Order.GetOrder(id)
}
func NewOrderProxy(order OrderInfoer) OrderInfoer {
	return &OrderProxy{Order: order}
}

func main() {

	//4.请求时，使用代理类获取接口实现类的数据

	proxy := NewOrderProxy(&OrderInfo{ID: 1, Name: "d001"})
	ret := proxy.GetOrder(1)
	fmt.Println("ret:", ret)

}
