package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"sync"
	"time"
)

// 一个具有层级关系的对象由一系列拥有父子关系的对象通过树形结构组成，子对象即可被串行执行，也可被并发执行

var (
	// ErrConcurrencyComponentTimeout 并发组件业务超时
	ErrConcurrencyComponentTimeout = errors.New("concurrency component timeout")
)

// Component 组件接口
type Component interface {
	// 添加一个子组件
	Mount(c Component, components ...Component) error
	// 移除一个子组件
	Remove(c Component) error
	// 执行当前组件业务:`BusinessLogicDo`和执行子组件:`ChildsDo`
	// ctx 业务上下文
	// currentComponent 当前组件
	// wg 父组件的WaitGroup对象
	// 区别1：增加了WaitGroup对象参数，目的是等待并发子组件的执行完成。
	Do(ctx *Context, currentComponent Component, wg *sync.WaitGroup) error
	// 执行当前组件业务逻辑
	// resChan 回写当前组件业务执行结果的channel
	// 区别2：增加了一个channel参数，目的是并发组件执行逻辑时引入了超时机制，需要一个channel接受组件的执行结果
	BusinessLogicDo(resChan chan interface{}) error
	// 执行子组件
	ChildrenDo(ctx *Context) error
}

// 并发子组件需要设置超时时间：防止子组件执行时间过长，解决方案关键字context.WithTimeout
// 区分普通组件和并发组件：合成复用基础组件，封装为并发基础组件
// 拥有并发子组件的父组件需要等待并发子组件执行完毕(包含超时)，解决方案关键字sync.WaitGroup
// 并发子组件执行自身业务逻辑是需检测超时：防止子组件内部执行业务逻辑时间过长，解决方案关键字select和<-ctx.Done()

// Context 业务上下文
type Context struct {
	// context.WithTimeout派生的子上下文
	TimeoutCtx context.Context
	// 超时函数
	context.CancelFunc
}

// GetContext 获取业务上下文实例
// d 超时时间
func GetContext(d time.Duration) *Context {
	c := &Context{}
	c.TimeoutCtx, c.CancelFunc = context.WithTimeout(context.Background(), d)
	return c
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

// Do 执行子组件
// ctx 业务上下文
// currentComponent 当前组件
// wg 父组件的waitGroup对象
func (bc *BaseComponent) Do(ctx *Context, currentComponent Component, _ *sync.WaitGroup) (err error) {
	//执行当前组件业务代码
	err = currentComponent.BusinessLogicDo(nil)
	if err != nil {
		return err
	}
	// 执行子组件
	return currentComponent.ChildrenDo(ctx)
}

// ChildrenDo 执行子组件
func (bc *BaseComponent) ChildrenDo(ctx *Context) (err error) {
	// 执行子组件
	for _, childComponent := range bc.ChildComponents {
		if err = childComponent.Do(ctx, childComponent, nil); err != nil {
			return err
		}
	}
	return
}

// BusinessLogicDo 当前组件业务逻辑代码填充处
func (bc *BaseComponent) BusinessLogicDo(chan interface{}) (err error) {
	// do nothing
	return
}

// BaseConcurrencyComponent 并发基础组件
type BaseConcurrencyComponent struct {
	// 合成复用基础组件
	BaseComponent
	// 当前组件是否有并发子组件
	HasChildConcurrencyComponents bool
	// 并发子组件列表
	ChildConcurrencyComponents []Component
	// wg 对象
	*sync.WaitGroup
	// 当前组件业务执行结果channel
	logicResChan chan interface{}
	// 当前组件执行过程中的错误信息
	Err error
}

// ChildrenDo 执行子组件
func (bc *BaseConcurrencyComponent) ChildrenDo(ctx *Context) (err error) {
	if bc.WaitGroup == nil {
		bc.WaitGroup = &sync.WaitGroup{}
	}
	// 执行并发子组件
	for _, childComponent := range bc.ChildConcurrencyComponents {
		bc.WaitGroup.Add(1)
		go childComponent.Do(ctx, childComponent, bc.WaitGroup)
	}
	// 执行子组件
	for _, childComponent := range bc.ChildComponents {
		if err = childComponent.Do(ctx, childComponent, nil); err != nil {
			return err
		}
	}
	if bc.HasChildConcurrencyComponents {
		// 等待并发组件执行结果
		bc.WaitGroup.Wait()
	}
	return
}

// Do 执行子组件
// ctx 业务上下文
// currentComponent 当前组件
// wg 父组件的waitGroup对象
func (bc *BaseConcurrencyComponent) Do(ctx *Context, currentComponent Component, wg *sync.WaitGroup) (err error) {
	defer wg.Done()
	// 初始化并发子组件channel
	if bc.logicResChan == nil {
		bc.logicResChan = make(chan interface{}, 1)
	}

	go currentComponent.BusinessLogicDo(bc.logicResChan)

	select {
	// 等待业务执行结果
	case <-bc.logicResChan:
		// 业务执行结果
		fmt.Println(runFuncName(), "bc.BusinessLogicDo wait.done...")
		break
	// 超时等待
	case <-ctx.TimeoutCtx.Done():
		// 超时退出
		fmt.Println(runFuncName(), "bc.BusinessLogicDo timeout...")
		bc.Err = ErrConcurrencyComponentTimeout
		break
	}
	// 执行子组件
	err = currentComponent.ChildrenDo(ctx)
	return
}

// Remove 移除一个子组件
func (bc *BaseConcurrencyComponent) Remove(c Component) (err error) {
	if len(bc.ChildComponents) == 0 {
		return
	}
	for k, childComponent := range bc.ChildComponents {
		if c == childComponent {
			fmt.Println(runFuncName(), "移除:", reflect.TypeOf(childComponent))
			bc.ChildComponents = append(bc.ChildComponents[:k], bc.ChildComponents[k+1:]...)
		}
	}
	for k, childComponent := range bc.ChildConcurrencyComponents {
		if c == childComponent {
			fmt.Println(runFuncName(), "移除:", reflect.TypeOf(childComponent))
			bc.ChildConcurrencyComponents = append(bc.ChildComponents[:k], bc.ChildComponents[k+1:]...)
		}
	}
	return
}

// MountConcurrency 挂载一个并发子组件
func (bc *BaseConcurrencyComponent) MountConcurrency(c Component, components ...Component) (err error) {
	bc.HasChildConcurrencyComponents = true
	bc.ChildConcurrencyComponents = append(bc.ChildConcurrencyComponents, c)
	if len(components) == 0 {
		return
	}
	bc.ChildConcurrencyComponents = append(bc.ChildConcurrencyComponents, components...)
	return
}

// CheckoutPageComponent 订单结算页面组件
type CheckoutPageComponent struct {
	// 合成复用基础组件
	BaseConcurrencyComponent
}

// BusinessLogicDo 当前组件业务逻辑代码填充处
func (bc *CheckoutPageComponent) BusinessLogicDo(chan interface{}) (err error) {
	// 当前组件的业务逻辑写这
	fmt.Println(runFuncName(), "订单结算页面组件...")
	return
}

// AddressComponent 地址组件
type AddressComponent struct {
	// 合成复用基础组件
	BaseConcurrencyComponent
}

// BusinessLogicDo 并发组件实际填充业务逻辑的地方
func (bc *AddressComponent) BusinessLogicDo(resChan chan interface{}) error {
	fmt.Println(runFuncName(), "地址组件...")
	fmt.Println(runFuncName(), "获取地址信息 ing...")

	// 模拟远程调用地址服务
	_, _ = http.Get("http://example.com/")

	resChan <- struct{}{} // 写入业务执行结果
	fmt.Println(runFuncName(), "获取地址信息 done...")
	return nil
}

// PayMethodComponent 支付方式组件
type PayMethodComponent struct {
	// 合成复用基础组件
	BaseConcurrencyComponent
}

// BusinessLogicDo 并发组件实际填充业务逻辑的地方
func (bc *PayMethodComponent) BusinessLogicDo(resChan chan interface{}) (err error) {
	// 当前组件的业务逻辑写这
	fmt.Println(runFuncName(), "支付方式组件...")
	fmt.Println(runFuncName(), "获取支付方式 ing...")
	// 模拟远程调用地址服务 略
	resChan <- struct{}{}
	fmt.Println(runFuncName(), "获取支付方式 done...")
	return nil
}

// StoreComponent 店铺组件
type StoreComponent struct {
	// 合成复用基础组件
	BaseComponent
}

// BusinessLogicDo 并发组件实际填充业务逻辑的地方
func (bc *StoreComponent) BusinessLogicDo(chan interface{}) (err error) {
	// 当前组件的业务逻辑写这
	fmt.Println(runFuncName(), "店铺组件...")
	return
}

// SkuComponent 商品组件
type SkuComponent struct {
	// 合成复用基础组件
	BaseComponent
}

// BusinessLogicDo 并发组件实际填充业务逻辑的地方
func (bc *SkuComponent) BusinessLogicDo(chan interface{}) (err error) {
	// 当前组件的业务逻辑写这
	fmt.Println(runFuncName(), "商品组件...")
	return
}

// PromotionComponent 优惠信息组件
type PromotionComponent struct {
	// 合成复用基础组件
	BaseComponent
}

// BusinessLogicDo 并发组件实际填充业务逻辑的地方
func (bc *PromotionComponent) BusinessLogicDo(chan interface{}) (err error) {
	// 当前组件的业务逻辑写这
	fmt.Println(runFuncName(), "优惠信息组件...")
	return
}

// ExpressComponent 物流组件
type ExpressComponent struct {
	// 合成复用基础组件
	BaseComponent
}

// BusinessLogicDo 并发组件实际填充业务逻辑的地方
func (bc *ExpressComponent) BusinessLogicDo(chan interface{}) (err error) {
	// 当前组件的业务逻辑写这
	fmt.Println(runFuncName(), "物流组件...")
	return
}

// AfterSaleComponent 售后组件
type AfterSaleComponent struct {
	// 合成复用基础组件
	BaseComponent
}

// BusinessLogicDo 并发组件实际填充业务逻辑的地方
func (bc *AfterSaleComponent) BusinessLogicDo(chan interface{}) (err error) {
	// 当前组件的业务逻辑写这
	fmt.Println(runFuncName(), "售后组件...")
	return
}

// InvoiceComponent 发票组件
type InvoiceComponent struct {
	// 合成复用基础组件
	BaseConcurrencyComponent
}

// BusinessLogicDo 并发组件实际填充业务逻辑的地方
func (bc *InvoiceComponent) BusinessLogicDo(resChan chan interface{}) (err error) {
	// 当前组件的业务逻辑写这
	fmt.Println(runFuncName(), "发票组件...")
	fmt.Println(runFuncName(), "获取发票信息 ing...")
	// 模拟远程调用地址服务 略
	resChan <- struct{}{} // 写入业务执行结果
	fmt.Println(runFuncName(), "获取发票信息 done...")
	return
}

// CouponComponent 优惠券组件
type CouponComponent struct {
	// 合成复用基础组件
	BaseConcurrencyComponent
}

// BusinessLogicDo 并发组件实际填充业务逻辑的地方
func (bc *CouponComponent) BusinessLogicDo(resChan chan interface{}) (err error) {
	// 当前组件的业务逻辑写这
	fmt.Println(runFuncName(), "优惠券组件...")
	fmt.Println(runFuncName(), "获取最优优惠券 ing...")

	// 模拟远程调用优惠券服务
	_, _ = http.Get("http://example.com/")

	// 写入业务执行结果
	resChan <- struct{}{}
	fmt.Println(runFuncName(), "获取最优优惠券 done...")
	return
}

// GiftCardComponent 礼品卡组件
type GiftCardComponent struct {
	// 合成复用基础组件
	BaseConcurrencyComponent
}

// BusinessLogicDo 并发组件实际填充业务逻辑的地方
func (bc *GiftCardComponent) BusinessLogicDo(resChan chan interface{}) (err error) {
	// 当前组件的业务逻辑写这
	fmt.Println(runFuncName(), "礼品卡组件...")
	fmt.Println(runFuncName(), "获取礼品卡信息 ing...")
	// 模拟远程调用地址服务 略
	resChan <- struct{}{} // 写入业务执行结果
	fmt.Println(runFuncName(), "获取礼品卡信息 done...")
	return
}

// OrderComponent 订单金额详细信息组件
type OrderComponent struct {
	// 合成复用基础组件
	BaseComponent
}

// BusinessLogicDo 当前组件业务逻辑代码填充处
func (bc *OrderComponent) BusinessLogicDo(chan interface{}) (err error) {
	// 当前组件的业务逻辑写这
	fmt.Println(runFuncName(), "订单金额详细信息组件...")
	return
}

// Demo 示例
func Demo() {
	// 初始化订单结算页面 这个大组件
	checkoutPage := &CheckoutPageComponent{}

	// 挂载子组件
	storeComponent := &StoreComponent{}
	skuComponent := &SkuComponent{}
	_ = skuComponent.Mount(
		&PromotionComponent{},
		&AfterSaleComponent{},
	)
	_ = storeComponent.Mount(
		skuComponent,
		&ExpressComponent{},
	)

	// ---挂载组件---

	// 普通组件
	_ = checkoutPage.Mount(
		storeComponent,
		&OrderComponent{},
	)
	// 并发组件
	_ = checkoutPage.MountConcurrency(
		&AddressComponent{},
		&PayMethodComponent{},
		&InvoiceComponent{},
		&CouponComponent{},
		&GiftCardComponent{},
	)

	// 初始化业务上下文 并设置超时时间
	ctx := GetContext(5 * time.Second)
	defer ctx.CancelFunc()
	// 开始构建页面组件数据
	_ = checkoutPage.ChildrenDo(ctx)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	Demo()
}

// 获取正在运行的函数名
func runFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}
