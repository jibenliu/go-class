package main

import (
	"fmt"
	"reflect"
	"runtime"
)
// refer : https://mp.weixin.qq.com/s/f2rhliZFvHcuykllFPxL2A
/**
	// ------------这里实现一个具体的“主题”------------

具体订单取消的动作实现“主题”(被观察者)接口`Observable`。得到一个具体的“主题”:

- 订单取消的动作的“主题”结构体`ObservableConcrete`
    +  成员属性`observerList []ObserverInterface`:订阅者列表
    +  具体方法`Attach`: 增加子逻辑
    +  具体方法`Detach`: 删除子逻辑
    +  具体方法`Notify`: 通知子逻辑

// ------------这里实现所有具体的“订阅者”------------

子逻辑实现“订阅者”接口`ObserverInterface`:

- 具体“订阅者”也就是子逻辑`OrderStatus`
    +  实现方法`Do`: 修改订单状态
- 具体“订阅者”也就是子逻辑`OrderStatusLog`
    +  实现方法`Do`: 记录订单状态变更日志
- 具体“订阅者”也就是子逻辑`CouponRefund`
    +  实现方法`Do`: 退优惠券
- 具体“订阅者”也就是子逻辑`PromotionRefund`
    +  实现方法`Do`: 还优惠活动资格
- 具体“订阅者”也就是子逻辑`StockRefund`
    +  实现方法`Do`: 还库存
- 具体“订阅者”也就是子逻辑`GiftCardRefund`
    +  实现方法`Do`: 还礼品卡
- 具体“订阅者”也就是子逻辑`WalletRefund`
    +  实现方法`Do`: 退钱包余额
- 具体“订阅者”也就是子逻辑`DeliverBillStatus`
    +  实现方法`Do`: 修改发货单状态
- 具体“订阅者”也就是子逻辑`DeliverBillStatusLog`
    +  实现方法`Do`: 记录发货单状态变更日志
- 具体“订阅者”也就是子逻辑`Refund`
    +  实现方法`Do`: 生成退款单
- 具体“订阅者”也就是子逻辑`Invoice`
    +  实现方法`Do`: 生成发票-红票
- 具体“订阅者”也就是子逻辑`Email`
    +  实现方法`Do`: 发邮件
- 具体“订阅者”也就是子逻辑`Sms`
    +  实现方法`Do`: 发短信
- 具体“订阅者”也就是子逻辑`WechatNotify`
    +  实现方法`Do`: 发微信消息
*/

// ObserverInterface 定义一个观察者的接口
type ObserverInterface interface {
	// 自身的业务
	Do(o Observable) error
}

// Observable 被观察者
type Observable interface {
	Attach(observer ...ObserverInterface) Observable
	Detach(observer ObserverInterface) Observable
	Notify() error
}

// ObservableConcrete 一个具体的 订单状态变化的被观察者
type ObservableConcrete struct {
	observerList []ObserverInterface
}

// Attach 注册观察者
// @param $observer ObserverInterface 观察者列表
func (o *ObservableConcrete) Attach(observer ...ObserverInterface) Observable {
	o.observerList = append(o.observerList, observer...)
	return o
}

// Detach 注销观察者
// @param $observer ObserverInterface 待注销的观察者
func (o *ObservableConcrete) Detach(observer ObserverInterface) Observable {
	if len(o.observerList) == 0 {
		return o
	}
	for k, observerItem := range o.observerList {
		if observer == observerItem {
			fmt.Println(runFuncName(), "注销:", reflect.TypeOf(observer))
			o.observerList = append(o.observerList[:k], o.observerList[k+1:]...)
		}
	}
	return o
}

// Notify 通知观察者
func (o *ObservableConcrete) Notify() (err error) {
	for _, observer := range o.observerList {
		if err = observer.Do(o); err != nil {
			return err
		}
	}
	return nil
}

// OrderStatus 修改订单状态
type OrderStatus struct {
}

func (observer *OrderStatus) Do(Observable) (err error) {
	fmt.Println(runFuncName(), "修改订单状态...")
	return
}

// OrderStatusLog 记录订单状态变更日志
type OrderStatusLog struct {
}

//Do 具体业务
func (observer *OrderStatusLog) Do(Observable) (err error) {
	fmt.Println(runFuncName(), "记录订单状态变更日志...")
	return
}

// CouponRefund 退优惠券
type CouponRefund struct {
}

// Do 具体业务
// PromotionRefund 还优惠活动资格
type PromotionRefund struct {
}

// Do 具体业务
func (observer *PromotionRefund) Do(Observable) (err error) {
	// code...
	fmt.Println(runFuncName(), "还优惠活动资格...")
	return
}

// StockRefund 还库存
type StockRefund struct {
}

// Do 具体业务
func (observer *StockRefund) Do(Observable) (err error) {
	// code...
	fmt.Println(runFuncName(), "还库存...")
	return
}

// GiftCardRefund 还礼品卡
type GiftCardRefund struct {
}

// Do 具体业务
func (observer *GiftCardRefund) Do(Observable) (err error) {
	// code...
	fmt.Println(runFuncName(), "还礼品卡...")
	return
}

// WalletRefund 退钱包余额
type WalletRefund struct {
}

// Do 具体业务
func (observer *WalletRefund) Do(Observable) (err error) {
	// code...
	fmt.Println(runFuncName(), "退钱包余额...")
	return
}

// DeliverBillStatus 修改发货单状态
type DeliverBillStatus struct {
}

// Do 具体业务
func (observer *DeliverBillStatus) Do(Observable) (err error) {
	// code...
	fmt.Println(runFuncName(), "修改发货单状态...")
	return
}

// DeliverBillStatusLog 记录发货单状态变更日志
type DeliverBillStatusLog struct {
}

// Do 具体业务
func (observer *DeliverBillStatusLog) Do(Observable) (err error) {
	// code...
	fmt.Println(runFuncName(), "记录发货单状态变更日志...")
	return
}

// Refund 生成退款单
type Refund struct {
}

// Do 具体业务
func (observer *Refund) Do(Observable) (err error) {
	// code...
	fmt.Println(runFuncName(), "生成退款单...")
	return
}

// Invoice 生成发票-红票
type Invoice struct {
}

// Do 具体业务
func (observer *Invoice) Do(Observable) (err error) {
	// code...
	fmt.Println(runFuncName(), "生成发票-红票...")
	return
}

// Email 发邮件
type Email struct {
}

// Do 具体业务
func (observer *Email) Do(Observable) (err error) {
	// code...
	fmt.Println(runFuncName(), "发邮件...")
	return
}

// Sms 发短信
type Sms struct {
}

// Do 具体业务
func (observer *Sms) Do(Observable) (err error) {
	// code...
	fmt.Println(runFuncName(), "发短信...")
	return
}

// WechatNotify 发微信消息
type WechatNotify struct {
}

// Do 具体业务
func (observer *WechatNotify) Do(Observable) (err error) {
	// code...
	fmt.Println(runFuncName(), "发微信消息...")
	return
}
func (observer *CouponRefund) Do(Observable) (err error) {
	// code...
	fmt.Println(runFuncName(), "退优惠券...")
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
	// 创建 未支付取消订单 “主题”
	fmt.Println("----------------------- 未支付取消订单 “主题”")
	orderUnPaidCancelSubject := &ObservableConcrete{}
	orderUnPaidCancelSubject.Attach(
		&OrderStatus{},
		&OrderStatusLog{},
		&CouponRefund{},
		&PromotionRefund{},
		&StockRefund{},
	)

	_ = orderUnPaidCancelSubject.Notify()

	// 创建 超时关单 “主题”
	fmt.Println("----------------------- 超时关单 “主题”")
	orderOverTimeSubject := &ObservableConcrete{}
	orderOverTimeSubject.Attach(
		&OrderStatus{},
		&OrderStatusLog{},
		&CouponRefund{},
		&PromotionRefund{},
		&StockRefund{},
		&Email{},
		&Sms{},
		&WechatNotify{},
	)
	_ = orderOverTimeSubject.Notify()

	// 创建 已支付取消订单 “主题”
	fmt.Println("----------------------- 已支付取消订单 “主题”")
	orderPaidCancelSubject := &ObservableConcrete{}
	orderPaidCancelSubject.Attach(
		&OrderStatus{},
		&OrderStatusLog{},
		&CouponRefund{},
		&PromotionRefund{},
		&StockRefund{},
		&GiftCardRefund{},
		&WalletRefund{},
		&Refund{},
		&Invoice{},
		&Email{},
		&Sms{},
		&WechatNotify{},
	)
	_ = orderPaidCancelSubject.Notify()

	// 创建 取消发货单 “主题”
	fmt.Println("----------------------- 取消发货单 “主题”")
	deliverBillCancelSubject := &ObservableConcrete{}
	deliverBillCancelSubject.Attach(
		&OrderStatus{},
		&OrderStatusLog{},
		&DeliverBillStatus{},
		&DeliverBillStatusLog{},
		&StockRefund{},
		&GiftCardRefund{},
		&WalletRefund{},
		&Refund{},
		&Invoice{},
		&Email{},
		&Sms{},
		&WechatNotify{},
	)
	_ = deliverBillCancelSubject.Notify()

	// 创建 拒收 “主题”
	fmt.Println("----------------------- 拒收 “主题”")
	deliverBillRejectSubject := &ObservableConcrete{}
	deliverBillRejectSubject.Attach(
		&OrderStatus{},
		&OrderStatusLog{},
		&DeliverBillStatus{},
		&DeliverBillStatusLog{},
		&StockRefund{},
		&GiftCardRefund{},
		&WalletRefund{},
		&Refund{},
		&Invoice{},
		&Email{},
		&Sms{},
		&WechatNotify{},
	)
	_ = deliverBillRejectSubject.Notify()

	// 未来可以快速的根据业务的变化 创建新的主题 从而快速构建新的业务接口
	fmt.Println("----------------------- 未来的扩展...")
}
