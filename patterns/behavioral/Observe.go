/*

例子：订单创建成功后，短信通知供应商、客户。
实际业务中是这样使用的：EventSubject被观察者有短信和邮件两种。观察者：短信只发送给供应商、客户，邮件发送给操作、销售。
  该例子只演示 短信发送给供应商、和客户。
 */

package main

import "fmt"

// 被观察者接口
type EventSubject interface {
	Notify() // 当被观察者数据改变时，通知已注册的观察者
	SetData(string)
	AddObserve(ob EventObserver)
	RemoveObserve(ob EventObserver)
}

// 观察者接口
type EventObserver interface {
	Receive(string)
}

////////////////////////////// 被观察者对象 SMSSubject ///////////////////

// 实例化 被观察者对象 发送短信
type SMSSubject struct {
	data      string          // 观察者监听的数据
	observers []EventObserver // 存储已注册的观察者对象的容器
}

// 更改被观察者对象的数据
func (s *SMSSubject) SetData(data string) {
	s.data = data
}

// 注册观察者
func (s *SMSSubject) AddObserve(ob EventObserver) {
	s.observers = append(s.observers, ob)
}

// 删除观察者
func (s *SMSSubject) RemoveObserve(ob EventObserver) {
	for k, v := range s.observers {
		if v == ob {
			s.observers = append(s.observers[:k], s.observers[k+1:]...)
		}
	}
}

// 通知所有已注册的观察者，变更自身的状态
func (s *SMSSubject) Notify() {
	for _, v := range s.observers {
		v.Receive(s.data)
	}
}

// 创建被观察者对象
func NewSMSSubject(data string) *SMSSubject {
	return &SMSSubject{data: data, observers: []EventObserver{}}
}

////////////////////////////// 观察者对象 SupplierObserve ///////////////////

// 一个观察者对象
type SupplierObserve struct {
	data string
}

// 接收信息
func (s *SupplierObserve) Receive(data string) {
	s.data = data
	fmt.Println("supplier receive data:", data)

}

// 创建实例
func NewSupplierObserve(data string) *SupplierObserve {
	return &SupplierObserve{data: data}
}

////////////////////////////// 观察者对象 CustomerObserve ///////////////////

// 一个观察者对象
type CustomerObserve struct {
	data string
}

// 接收信息
func (c *CustomerObserve) Receive(data string) {
	c.data = data
	fmt.Println("customer receive data:", data)
}

// 创建实例
func NewCustomerObserve(data string) *CustomerObserve {
	return &CustomerObserve{data: data}
}

func main() {

	sms := NewSMSSubject("新订单")
	supplier := NewSupplierObserve("")
	customer := NewCustomerObserve("")

	//
	sms.AddObserve(supplier)
	sms.AddObserve(customer)
	//
	sms.Notify()

	sms.SetData("新订单2")
	sms.Notify()
}
