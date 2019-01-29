package main

import (
	"fmt"
)

// 状态模式：满足一定条件下自动更换
// 策略模式：让用户指定更换的策略
//策略接口
type PayOrderAmounter interface {

	PayAmount(amount float64) float64
}

//策略实现类

//正常支付的订单，没有折扣
type PayNormalOrder struct {

	OrderId float64
}

func newPayNoramlOrder()*PayNormalOrder  {
	p:=new(PayNormalOrder)
	p.OrderId=1001
	return p
}
func (p *PayNormalOrder)PayAmount(amount float64)float64  {
	fmt.Println("PayNormalOrder OrderId:",p.OrderId,"价格：",amount)
	return amount
}

//打折支付的订单
type PayRebateOrder struct {
	OrderId float64
	Rebate float64 //折扣
}

func (p *PayRebateOrder)PayAmount(amount float64)float64  {
	fmt.Println("PayRebateOrder OrderId:",p.OrderId,"价格:",p.Rebate*amount)
	return  p.Rebate*amount
}

func newPayRebateOrder(rebate float64)*PayRebateOrder  {
	p:=new(PayRebateOrder)
	p.OrderId=1002
	p.Rebate=rebate
	return p
}

//策略上下文

type PayAmountContext struct {

	Pay PayOrderAmounter
}

//调用方法一：使用简单工厂模式生成策略类

func NewPayAmountContext(cashType string)*PayAmountContext  {
	p:=new(PayAmountContext)
	switch cashType {
		case "八折":
			p.Pay=newPayRebateOrder(0.8)
		default:
			p.Pay=newPayNoramlOrder()
	}
	return p
}

func (c *PayAmountContext)PayAmount(amount float64)float64  {
	return  c.Pay.PayAmount(amount)
}

//调用方法二：使用New实体生成策略类
func NewPayContextByEntity(cashType PayOrderAmounter)* PayAmountContext  {
	p:=new(PayAmountContext)
	p.Pay=cashType
	return p
}

func main() {

	//方法一 工厂模式：
	p:=NewPayAmountContext("八折")
	p.PayAmount(100)

	p2:=NewPayAmountContext("")
	p2.PayAmount(100)

	fmt.Println("---------------------------------")
	//方法二 单独调用：
	NewPayContextByEntity(newPayNoramlOrder()).Pay.PayAmount(100)
	NewPayContextByEntity(newPayRebateOrder(0.8)).Pay.PayAmount(100)

}
