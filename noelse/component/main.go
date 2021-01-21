package main

import (
	"fmt"
	"reflect"
	"runtime"
)

/**
一个父类(抽象类)：
- 成员属性
	+ `ChildComponents`: 子组件列表
- 成员方法
	+ `Mount`: 实现添加一个子组件
	+ `Remove`: 实现移除一个子组件
	+ `Do`: 抽象方法

组件一，订单结算页面组件类(继承父类、看成一个大的组件)：
- 成员方法
	+ `Do`: 执行当前组件的逻辑，执行子组件的逻辑

组件二，地址组件(继承父类)：
- 成员方法
	+ `Do`: 执行当前组件的逻辑，执行子组件的逻辑

组件三，支付方式组件(继承父类)：
- 成员方法
	+ `Do`: 执行当前组件的逻辑，执行子组件的逻辑

组件四，店铺组件(继承父类)：
- 成员方法
	+ `Do`: 执行当前组件的逻辑，执行子组件的逻辑

组件五，商品组件(继承父类)：
- 成员方法
	+ `Do`: 执行当前组件的逻辑，执行子组件的逻辑

组件六，优惠信息组件(继承父类)：
- 成员方法
	+ `Do`: 执行当前组件的逻辑，执行子组件的逻辑

组件七，物流组件(继承父类)：
- 成员方法
	+ `Do`: 执行当前组件的逻辑，执行子组件的逻辑

组件八，发票组件(继承父类)：
- 成员方法
	+ `Do`: 执行当前组件的逻辑，执行子组件的逻辑

组件九，优惠券组件(继承父类)：
- 成员方法
	+ `Do`: 执行当前组件的逻辑，执行子组件的逻辑

组件十，礼品卡组件(继承父类)：
- 成员方法
	+ `Do`: 执行当前组件的逻辑，执行子组件的逻辑

组件十一，订单金额详细信息组件(继承父类)：
- 成员方法
	+ `Do`: 执行当前组件的逻辑，执行子组件的逻辑
组件十二，售后组件(继承父类)：
- 成员方法
	+ `Do`: 执行当前组件的逻辑，执行子组件的逻辑


@TODO 但是，golang里没有的继承的概念，要复用成员属性ChildComponents、成员方法Mount、成员方法Remove怎么办呢？我们使用合成复用的特性变相达到“继承复用”的目的，如下：


一个接口(interface)：
+ 抽象方法`Mount`: 添加一个子组件
+ 抽象方法`Remove`: 移除一个子组件
+ 抽象方法`Do`: 执行组件&子组件

一个基础结构体`BaseComponent`：
- 成员属性
	+ `ChildComponents`: 子组件列表
- 成员方法
	+ 实体方法`Mount`: 添加一个子组件
	+ 实体方法`Remove`: 移除一个子组件
	+ 实体方法`ChildrenDo`: 执行子组件

组件一，订单结算页面组件类：
- 合成复用基础结构体`BaseComponent`
- 成员方法
	+ `Do`: 执行当前组件的逻辑，执行子组件的逻辑

组件二，地址组件：
- 合成复用基础结构体`BaseComponent`
- 成员方法
	+ `Do`: 执行当前组件的逻辑，执行子组件的逻辑

组件三，支付方式组件：
- 合成复用基础结构体`BaseComponent`
- 成员方法
	+ `Do`: 执行当前组件的逻辑，执行子组件的逻辑

...略

组件十一，订单金额详细信息组件：
- 合成复用基础结构体`BaseComponent`
- 成员方法
	+ `Do`: 执行当前组件的逻辑，执行子组件的逻辑
*/

// Context 上下文
type Context struct {
}

// Component 组件接口
type Component interface {
	Mount(c Component, components ...Component) error
	Remove(c Component) error
	Do(ctx *Context) error
}

// BaseComponent 基础组件
// 实现Add:添加一个子组件
// 实现Remove:移除一个子组件
type BaseComponent struct {
	// 子组件列表
	ChildComponents []Component
}

// Mount 挂载一个子组件
func (bc *BaseComponent) Mount(c Component, components ...Component) (err error) {
	bc.ChildComponents = append(bc.ChildComponents, c)
	if len(components) == 0 {
		return
	}
	bc.ChildComponents = append(bc.ChildComponents, components...)
	return
}

// Remove 移除一个子组件
func (bc *BaseComponent) Remove(c Component) (err error) {
	if len(bc.ChildComponents) == 0 {
		return
	}
	for k, childComponent := range bc.ChildComponents {
		if c == childComponent {
			fmt.Println(runFuncName(), "移除:", reflect.TypeOf(childComponent))
			bc.ChildComponents = append(bc.ChildComponents[:k], bc.ChildComponents[k+1:]...)
		}
	}
	return
}

// Do 执行组件&子组件
func (bc *BaseComponent) Do(*Context) (err error) {
	return
}

func (bc *BaseComponent) ChildrenDo(ctx *Context) (err error) {
	// 执行子组件
	for _, childComponent := range bc.ChildComponents {
		if err = childComponent.Do(ctx); err != nil {
			return err
		}
	}
	return
}

// CheckoutPageComponent 订单结算页面组件
type CheckoutPageComponent struct {
	// 合成复用基础组件
	BaseComponent
}

// Do 执行组件&子组件
func (bc *CheckoutPageComponent) Do(ctx *Context) (err error) {
	// 当前组件的业务逻辑写这
	fmt.Println(runFuncName(), "订单结算页面组件...")

	// 执行子组件
	_ = bc.ChildrenDo(ctx)

	// 当前组件的业务逻辑写这

	return
}

// AddressComponent 地址组件
type AddressComponent struct {
	// 合成复用基础组件
	BaseComponent
}

// Do 执行组件&子组件
func (bc *AddressComponent) Do(ctx *Context) (err error) {
	// 当前组件的业务逻辑写这
	fmt.Println(runFuncName(), "地址组件...")

	// 执行子组件
	_ = bc.ChildrenDo(ctx)

	// 当前组件的业务逻辑写这

	return
}

// PayMethodComponent 支付方式组件
type PayMethodComponent struct {
	// 合成复用基础组件
	BaseComponent
}

// Do 执行组件&子组件
func (bc *PayMethodComponent) Do(ctx *Context) (err error) {
	// 当前组件的业务逻辑写这
	fmt.Println(runFuncName(), "支付方式组件...")

	// 执行子组件
	_ = bc.ChildrenDo(ctx)

	// 当前组件的业务逻辑写这

	return
}

// StoreComponent 店铺组件
type StoreComponent struct {
	// 合成复用基础组件
	BaseComponent
}

// Do 执行组件&子组件
func (bc *StoreComponent) Do(ctx *Context) (err error) {
	// 当前组件的业务逻辑写这
	fmt.Println(runFuncName(), "店铺组件...")

	// 执行子组件
	_ = bc.ChildrenDo(ctx)

	// 当前组件的业务逻辑写这

	return
}

// SkuComponent 商品组件
type SkuComponent struct {
	// 合成复用基础组件
	BaseComponent
}

// Do 执行组件&子组件
func (bc *SkuComponent) Do(ctx *Context) (err error) {
	// 当前组件的业务逻辑写这
	fmt.Println(runFuncName(), "商品组件...")

	// 执行子组件
	_ = bc.ChildrenDo(ctx)

	// 当前组件的业务逻辑写这

	return
}

// PromotionComponent 优惠信息组件
type PromotionComponent struct {
	// 合成复用基础组件
	BaseComponent
}

// Do 执行组件&子组件
func (bc *PromotionComponent) Do(ctx *Context) (err error) {
	// 当前组件的业务逻辑写这
	fmt.Println(runFuncName(), "优惠信息组件...")

	// 执行子组件
	_ = bc.ChildrenDo(ctx)

	// 当前组件的业务逻辑写这

	return
}

// ExpressComponent 物流组件
type ExpressComponent struct {
	// 合成复用基础组件
	BaseComponent
}

// Do 执行组件&子组件
func (bc *ExpressComponent) Do(ctx *Context) (err error) {
	// 当前组件的业务逻辑写这
	fmt.Println(runFuncName(), "物流组件...")

	// 执行子组件
	_ = bc.ChildrenDo(ctx)

	// 当前组件的业务逻辑写这

	return
}

// AftersaleComponent 售后组件
type AftersaleComponent struct {
	// 合成复用基础组件
	BaseComponent
}

// Do 执行组件&子组件
func (bc *AftersaleComponent) Do(ctx *Context) (err error) {
	// 当前组件的业务逻辑写这
	fmt.Println(runFuncName(), "售后组件...")

	// 执行子组件
	_ = bc.ChildrenDo(ctx)

	// 当前组件的业务逻辑写这

	return
}

// InvoiceComponent 发票组件
type InvoiceComponent struct {
	// 合成复用基础组件
	BaseComponent
}

// Do 执行组件&子组件
func (bc *InvoiceComponent) Do(ctx *Context) (err error) {
	// 当前组件的业务逻辑写这
	fmt.Println(runFuncName(), "发票组件...")

	// 执行子组件
	_ = bc.ChildrenDo(ctx)

	// 当前组件的业务逻辑写这

	return
}

// CouponComponent 优惠券组件
type CouponComponent struct {
	// 合成复用基础组件
	BaseComponent
}

// Do 执行组件&子组件
func (bc *CouponComponent) Do(ctx *Context) (err error) {
	// 当前组件的业务逻辑写这
	fmt.Println(runFuncName(), "优惠券组件...")

	// 执行子组件
	_ = bc.ChildrenDo(ctx)

	// 当前组件的业务逻辑写这

	return
}

// GiftCardComponent 礼品卡组件
type GiftCardComponent struct {
	// 合成复用基础组件
	BaseComponent
}

// Do 执行组件&子组件
func (bc *GiftCardComponent) Do(ctx *Context) (err error) {
	// 当前组件的业务逻辑写这
	fmt.Println(runFuncName(), "礼品卡组件...")

	// 执行子组件
	_ = bc.ChildrenDo(ctx)

	// 当前组件的业务逻辑写这

	return
}

// OrderComponent 订单金额详细信息组件
type OrderComponent struct {
	// 合成复用基础组件
	BaseComponent
}

// Do 执行组件&子组件
func (bc *OrderComponent) Do(ctx *Context) (err error) {
	// 当前组件的业务逻辑写这
	fmt.Println(runFuncName(), "订单金额详细信息组件...")

	// 执行子组件
	_ = bc.ChildrenDo(ctx)

	// 当前组件的业务逻辑写这

	return
}

func main() {
	// 初始化订单结算页面 这个大组件
	checkoutPage := &CheckoutPageComponent{}

	// 挂载子组件
	storeComponent := &StoreComponent{}
	skuComponent := &SkuComponent{}
	_ = skuComponent.Mount(
		&PromotionComponent{},
		&AftersaleComponent{},
	)
	_ = storeComponent.Mount(
		skuComponent,
		&ExpressComponent{},
	)

	// 挂载组件
	_ = checkoutPage.Mount(
		&AddressComponent{},
		&PayMethodComponent{},
		storeComponent,
		&InvoiceComponent{},
		&CouponComponent{},
		&GiftCardComponent{},
		&OrderComponent{},
	)

	// 移除组件测试
	// checkoutPage.Remove(storeComponent)

	// 开始构建页面组件数据
	_ = checkoutPage.Do(&Context{})
}

// 获取正在运行的函数名
func runFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}
