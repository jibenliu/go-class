package main

//业务梳理
//我们以某团的订单支付页面为例，页面上的每一个支付选项都是一个支付策略。如下：
//
//用户可以使用：
//		美团支付(策略)
//		微信支付(策略)
//		支付宝支付(策略)

/***--------------------伪代码如下：-----------***/
//// 定义一个支付接口
//- `PaymentInterface`
//    + 抽象方法`Pay(ctx *Context) error`: 当前支付方式的支付逻辑
//	+ 抽象方法`Refund(ctx *Context) error`: 当前支付方式的退款逻辑
//
//// 定义具体的支付方式 实现接口`PaymentInterface`
//
//- 具体的微信支付方式`WechatPay`
//    +  实现方法`Pay`: 支付逻辑
//	+  实现方法`Refund`: 支付逻辑
//- 具体的支付宝支付网页版方式`AliPayWap`
//    +  实现方法`Pay`: 支付逻辑
//	+  实现方法`Refund`: 支付逻辑
//- 具体的支付宝支付网页版方式`BankPay`
//    +  实现方法`Pay`: 支付逻辑
//	+  实现方法`Refund`: 支付逻辑
//
//// 客户端代码
//通过接口参数pay_type的值判断是哪种支付方式策略

import (
	"fmt"
	"runtime"
)

//------------------------------------------------------------
//我的代码没有`else`系列
//策略模式
//@auhtor TIGERB<https://github.com/TIGERB>
//------------------------------------------------------------

const (
	// ConstWechatPay 微信支付
	ConstWechatPay = "wechat_pay"
	// ConstAliPayWap 支付宝支付 网页版
	ConstAliPayWap = "AliPayWapwap"
	// ConstBankPay 银行卡支付
	ConstBankPay = "quickbank"
)

// Context 上下文
type Context struct {
	// 用户选择的支付方式
	PayType string `json:"pay_type"`
}

// PaymentInterface 支付方式接口
type PaymentInterface interface {
	Pay(ctx *Context) error    // 支付
	Refund(ctx *Context) error // 退款
}

// WechatPay 微信支付
type WechatPay struct {
}

// Pay 当前支付方式的支付逻辑
func (p *WechatPay) Pay(*Context) (err error) {
	// 当前策略的业务逻辑写这
	fmt.Println(runFuncName(), "使用微信支付...")
	return
}

// Refund 当前支付方式的支付逻辑
func (p *WechatPay) Refund(*Context) (err error) {
	// 当前策略的业务逻辑写这
	fmt.Println(runFuncName(), "使用微信退款...")
	return
}

// AliPayWap 支付宝网页版
type AliPayWap struct {
}

// Pay 当前支付方式的支付逻辑
func (p *AliPayWap) Pay(*Context) (err error) {
	// 当前策略的业务逻辑写这
	fmt.Println(runFuncName(), "使用支付宝网页版支付...")
	return
}

// Refund 当前支付方式的支付逻辑
func (p *AliPayWap) Refund(*Context) (err error) {
	// 当前策略的业务逻辑写这
	fmt.Println(runFuncName(), "使用支付宝网页版退款...")
	return
}

// BankPay 银行卡支付
type BankPay struct {
}

// Pay 当前支付方式的支付逻辑
func (p *BankPay) Pay(*Context) (err error) {
	// 当前策略的业务逻辑写这
	fmt.Println(runFuncName(), "使用银行卡支付...")
	return
}

// Refund 当前支付方式的支付逻辑
func (p *BankPay) Refund(*Context) (err error) {
	// 当前策略的业务逻辑写这
	fmt.Println(runFuncName(), "使用银行卡退款...")
	return
}

// 获取正在运行的函数名
func runFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}

func main() {
	// 相对于被调用的支付策略 这里就是支付策略的客户端

	// 业务上下文
	ctx := &Context{
		PayType: "wechat_pay",
	}

	// 获取支付方式
	var instance PaymentInterface
	switch ctx.PayType {
	case ConstWechatPay:
		instance = &WechatPay{}
	case ConstAliPayWap:
		instance = &AliPayWap{}
	case ConstBankPay:
		instance = &BankPay{}
	default:
		panic("无效的支付方式")
	}

	// 支付
	_ = instance.Pay(ctx)
}
