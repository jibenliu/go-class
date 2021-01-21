package main

import (
	"fmt"
	"runtime"
)

/**
一个父类(抽象类)：

- 成员属性
	+ `nextHandler`: 下一个等待被调用的对象实例
- 成员方法
	+ 实体方法`SetNext`: 实现把下一个对象的实例绑定到当前对象的`nextHandler`属性上
	+ 抽象方法`Do`: 当前对象业务逻辑入口
	+ 实体方法`Run`: 实现调用当前对象的`Do`，`nextHandler`不为空则调用`nextHandler.Do`

子类一(参数校验)
- 继承抽象类父类
- 实现抽象方法`Do`：具体的参数校验逻辑

子类二(获取地址信息)
- 继承抽象类父类
- 实现抽象方法`Do`：具体获取地址信息的逻辑

子类三(获取购物车数据)
- 继承抽象类父类
- 实现抽象方法`Do`：具体获取购物车数据的逻辑

......略

子类X(以及未来会增加的逻辑)
- 继承抽象类父类
- 实现抽象方法`Do`：以及未来会增加的逻辑

@TODO 但是，golang里没有的继承的概念，要复用成员属性nextHandler、成员方法SetNext、成员方法Run怎么办呢？我们使用合成复用的特性变相达到“继承复用”的目的，如下：

一个接口(interface)：

- 抽象方法`SetNext`: 待实现把下一个对象的实例绑定到当前对象的`nextHandler`属性上
- 抽象方法`Do`: 待实现当前对象业务逻辑入口
- 抽象方法`Run`: 待实现调用当前对象的`Do`，`nextHandler`不为空则调用`nextHandler.Do`

一个基础结构体：

- 成员属性
	+ `nextHandler`: 下一个等待被调用的对象实例
- 成员方法
	+ 实体方法`SetNext`: 实现把下一个对象的实例绑定到当前对象的`nextHandler`属性上
	+ 实体方法`Run`: 实现调用当前对象的`Do`，`nextHandler`不为空则调用`nextHandler.Do`

子类一(参数校验)
- 合成复用基础结构体
- 实现抽象方法`Do`：具体的参数校验逻辑

子类二(获取地址信息)
- 合成复用基础结构体
- 实现抽象方法`Do`：具体获取地址信息的逻辑

子类三(获取购物车数据)
- 合成复用基础结构体
- 实现抽象方法`Do`：具体获取购物车数据的逻辑

......略

子类X(以及未来会增加的逻辑)
- 合成复用基础结构体
- 实现抽象方法`Do`：以及未来会增加的逻辑
*/

// Context Context
type Context struct {
}

// Handler 处理
type Handler interface {
	// 自身的业务
	Do(c *Context) error
	// 设置下一个对象
	SetNext(h Handler) Handler
	// 执行
	Run(c *Context) error
}

// Next 抽象出来的 可被合成复用的结构体
type Next struct {
	// 下一个对象
	nextHandler Handler
}

// SetNext 实现好的 可被复用的SetNext方法
// 返回值是下一个对象 方便写成链式代码优雅
// 例如 nullHandler.SetNext(argumentsHandler).SetNext(signHandler).SetNext(frequentHandler)
func (n *Next) SetNext(h Handler) Handler {
	n.nextHandler = h
	return h
}

//run执行
func (n *Next) Run(c *Context) (err error) {
	// 由于go无继承的概念 这里无法执行当前handler的Do
	// n.Do(c)
	if n.nextHandler != nil {
		// 合成复用下的变种
		// 执行下一个handler的Do
		if err = (n.nextHandler).Do(c); err != nil {
			return
		}
		// 执行下一个handler的Run
		return (n.nextHandler).Run(c)
	}
	return
}

// NullHandler 空Handler
// 由于go无继承的概念 作为链式调用的第一个载体 设置实际的下一个对象
type NullHandler struct {
	// 合成复用Next的`nextHandler`成员属性、`SetNext`成员方法、`Run`成员方法
	Next
}

// Do 空Handler的Do
func (h *NullHandler) Do(*Context) (err error) {
	// 空Handler 这里什么也不做 只是载体 do nothing...
	return
}

// ArgumentsHandler 校验参数的handler
type ArgumentsHandler struct {
	// 合成复用Next
	Next
}

// Do 校验参数的逻辑
func (h *ArgumentsHandler) Do(*Context) (err error) {
	fmt.Println(runFuncName(), "校验参数成功...")
	return
}

// AddressInfoHandler 地址信息handler
type AddressInfoHandler struct {
	// 合成复用Next
	Next
}

// Do 校验参数的逻辑
func (h *AddressInfoHandler) Do(*Context) (err error) {
	fmt.Println(runFuncName(), "获取地址信息...")
	fmt.Println(runFuncName(), "地址信息校验...")
	return
}

// CartInfoHandler 获取购物车数据handler
type CartInfoHandler struct {
	// 合成复用Next
	Next
}

// Do 校验参数的逻辑
func (h *CartInfoHandler) Do(*Context) (err error) {
	fmt.Println(runFuncName(), "获取购物车数据...")
	return
}

// StockInfoHandler 商品库存handler
type StockInfoHandler struct {
	// 合成复用Next
	Next
}

// Do 校验参数的逻辑
func (h *StockInfoHandler) Do(*Context) (err error) {
	fmt.Println(runFuncName(), "获取商品库存信息...")
	fmt.Println(runFuncName(), "商品库存校验...")
	return
}

// PromotionInfoHandler 获取优惠信息handler
type PromotionInfoHandler struct {
	// 合成复用Next
	Next
}

// Do 校验参数的逻辑
func (h *PromotionInfoHandler) Do(*Context) (err error) {
	fmt.Println(runFuncName(), "获取优惠信息...")
	return
}

// ShipmentInfoHandler 获取运费信息handler
type ShipmentInfoHandler struct {
	// 合成复用Next
	Next
}

// Do 校验参数的逻辑
func (h *ShipmentInfoHandler) Do(*Context) (err error) {
	fmt.Println(runFuncName(), "获取运费信息...")
	return
}

// PromotionUseHandler 使用优惠信息handler
type PromotionUseHandler struct {
	// 合成复用Next
	Next
}

// Do 校验参数的逻辑
func (h *PromotionUseHandler) Do(*Context) (err error) {
	fmt.Println(runFuncName(), "使用优惠信息...")
	return
}

// StockSubtractHandler 库存操作handler
type StockSubtractHandler struct {
	// 合成复用Next
	Next
}

// Do 校验参数的逻辑
func (h *StockSubtractHandler) Do(*Context) (err error) {
	fmt.Println(runFuncName(), "扣库存...")
	return
}

// CartDelHandler 清理购物车handler
type CartDelHandler struct {
	// 合成复用Next
	Next
}

// Do 校验参数的逻辑
func (h *CartDelHandler) Do(*Context) (err error) {
	fmt.Println(runFuncName(), "清理购物车...")
	// err = fmt.Errorf("CartDelHandler.Do fail")
	return
}

// DBTableOrderHandler 写订单表handler
type DBTableOrderHandler struct {
	// 合成复用Next
	Next
}

// Do 校验参数的逻辑
func (h *DBTableOrderHandler) Do(*Context) (err error) {
	fmt.Println(runFuncName(), "写订单表...")
	return
}

// DBTableOrderSkusHandler 写订单商品表handler
type DBTableOrderSkusHandler struct {
	// 合成复用Next
	Next
}

// Do 校验参数的逻辑
func (h *DBTableOrderSkusHandler) Do(*Context) (err error) {
	fmt.Println(runFuncName(), "写订单商品表...")
	return
}

// DBTableOrderPromotionsHandler 写订单优惠信息表handler
type DBTableOrderPromotionsHandler struct {
	// 合成复用Next
	Next
}

// Do 校验参数的逻辑
func (h *DBTableOrderPromotionsHandler) Do(*Context) (err error) {
	fmt.Println(runFuncName(), "写订单优惠信息表...")
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
	// 初始化空handler
	nullHandler := &NullHandler{}

	// 链式调用 代码是不是很优雅
	// 很明显的链 逻辑关系一览无余
	nullHandler.SetNext(&ArgumentsHandler{}).
		SetNext(&AddressInfoHandler{}).
		SetNext(&CartInfoHandler{}).
		SetNext(&StockInfoHandler{}).
		SetNext(&PromotionInfoHandler{}).
		SetNext(&ShipmentInfoHandler{}).
		SetNext(&PromotionUseHandler{}).
		SetNext(&StockSubtractHandler{}).
		SetNext(&CartDelHandler{}).
		SetNext(&DBTableOrderHandler{}).
		SetNext(&DBTableOrderSkusHandler{}).
		SetNext(&DBTableOrderPromotionsHandler{})
	//无限扩展代码...

	// 开始执行业务
	if err := nullHandler.Run(&Context{}); err != nil {
		// 异常
		fmt.Println("Fail | Error:" + err.Error())
		return
	}
	// 成功
	fmt.Println("Success")
	return
}
