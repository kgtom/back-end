/**
 * @Author: tom
 * @Date time: 2020-03-22 18:39
 * @description:
 * golang使用interface 实现开闭原则、依赖倒置原则
interface 设计api接口，不仅能满足当下功能，而且对未来代码有扩展，即调用未来哈

开闭原则：对系统需要添加一个新功能时，不是修改原有代码，而是通过扩展新增代码，即：对修改是关闭的，对扩展是开放的
依赖倒置原则：Dependence Inversion Principle，应该依赖于抽象，不依赖于具体实现,即：面向接口编程。

 * @return
*/

package practice

import "fmt"

//开闭原则：对系统需要添加一个新功能时，不是修改原有代码，而是通过扩展新增代码，即：对修改是关闭的，对扩展是开放的
//优点：代码可扩展、不影响之前功能
//缺点：代码写的多一点
//场景：用户资金账户：充值、支付、提现、转账,
//实现：之前用户资金账户只有充值、支付功能，逐渐新增提现、转账功能，
//使用开闭原则，抽象业务接口，新增的提现、转账增加代码，不修改之前充值、支付功能，保证之前功能不受影响。
type OperateAccount interface {
	OperateHandle() error //抽象的处理业务接口
}

//充值
type RechargeAccount struct {
}

func (r *RechargeAccount) OperateHandle() error {

	fmt.Println("recharge")
	return nil
}

//支付
type PayAccount struct {
}

func (p *PayAccount) OperateHandle() error {
	fmt.Println("pay")
	return nil
}

//提现
type WithdrawAccount struct {
}

func (w *WithdrawAccount) OperateHandler() error {
	fmt.Println("withdraw")
	return nil
}

//转账
type TransferAccount struct {
}

func (t *TransferAccount) OperateHandler() error {
	fmt.Println("transfer")
	return nil
}

func OperateImpl(o OperateAccount) {

	o.OperateHandle()
}

func InvokeOCP() {

	//充值 调用方式一
	r := &RechargeAccount{}
	r.OperateHandle()
	//充值 调用方式二
	OperateImpl(&RechargeAccount{})

	//支付
	OperateImpl(&PayAccount{})
}
