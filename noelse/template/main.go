package main

import (
	"fmt"
	"runtime"
)

//refer : https://mp.weixin.qq.com/s?__biz=MzAxMTA4Njc0OQ==&mid=2651439531&idx=3&sn=45bf4655646496a588584e80f0ff40b5&chksm=80bb1f59b7cc964f752d949cc04d99c8bc5f8d41b97bfa2e54176bf7d266d895a55e46e84422&scene=21#wechat_redirect
/**
一个抽象类
- 具体共有方法`Run`，里面定义了算法的执行步骤
- 具体私有方法，不会发生变化的具体方法
- 抽象方法，会发生变化的方法

子类一(按时间抽奖类型)
- 继承抽象类父类
- 实现抽象方法

子类二(按抽奖次数抽奖类型)
- 继承抽象类父类
- 实现抽象方法

子类三(按数额范围区间抽奖)
- 继承抽象类父类
- 实现抽象方法

@TODO 但是golang里面没有继承的概念，我们就把对抽象类里抽象方法的依赖转化成对接口interface里抽像方法的依赖，同时也可以利用合成复用的方式“继承”模板:

抽象行为的接口`BehaviorInterface`(包含如下需要实现的方法)
- 其他参数校验的方法`checkParams`
- 获取node奖品信息的方法`getPrizesByNode`

抽奖结构体类
- 具体共有方法`Run`，里面定义了算法的执行步骤
- 具体私有方法`checkParams` 里面的逻辑实际依赖的接口BehaviorInterface.checkParams(ctx)的抽象方法
- 具体私有方法`getPrizesByNode` 里面的逻辑实际依赖的接口BehaviorInterface.getPrizesByNode(ctx)的抽象方法
- 其他具体私有方法，不会发生变化的具体方法

实现`BehaviorInterface`的结构体一(按时间抽奖类型)
- 实现接口方法

实现`BehaviorInterface`的结构体二(按抽奖次数抽奖类型)
- 实现接口方法

实现`BehaviorInterface`的结构体三(按数额范围区间抽奖)
- 实现接口方法
*/

const (
	// ConstActTypeTime 按时间抽奖类型
	ConstActTypeTime int32 = 1
	// ConstActTypeTimes 按抽奖次数抽奖
	ConstActTypeTimes int32 = 2
	// ConstActTypeAmount 按数额范围区间抽奖
	ConstActTypeAmount int32 = 3
)

// Context 上下文
type Context struct {
	ActInfo *ActInfo
}

// ActInfo 上下文
type ActInfo struct {
	// 活动抽奖类型1: 按时间抽奖 2: 按抽奖次数抽奖 3:按数额范围区间抽奖
	ActivityType int32
	// 其他字段略
}

// BehaviorInterface 不同抽奖类型的行为差异的抽象接口
type BehaviorInterface interface {
	// 其他参数校验(不同活动类型实现不同)
	checkParams(ctx *Context) error
	// 获取node奖品信息(不同活动类型实现不同)
	getPrizesByNode(ctx *Context) error
}

// TimeDraw 具体抽奖行为
// 按时间抽奖类型 比如红包雨
type TimeDraw struct{}

// checkParams 其他参数校验(不同活动类型实现不同)
func (draw TimeDraw) checkParams(*Context) (err error) {
	fmt.Println(runFuncName(), "按时间抽奖类型:特殊参数校验...")
	return
}

// getPrizesByNode 获取node奖品信息(不同活动类型实现不同)
func (draw TimeDraw) getPrizesByNode(*Context) (err error) {
	fmt.Println(runFuncName(), "do nothing(抽取该场次的奖品即可，无需其他逻辑)...")
	return
}

// TimesDraw 具体抽奖行为
// 按抽奖次数抽奖类型 比如答题闯关
type TimesDraw struct{}

// checkParams 其他参数校验(不同活动类型实现不同)
func (draw TimesDraw) checkParams(*Context) (err error) {
	fmt.Println(runFuncName(), "按抽奖次数抽奖类型:特殊参数校验...")
	return
}

// getPrizesByNode 获取node奖品信息(不同活动类型实现不同)
func (draw TimesDraw) getPrizesByNode(*Context) (err error) {
	fmt.Println(runFuncName(), "1. 判断是该用户第几次抽奖...")
	fmt.Println(runFuncName(), "2. 获取对应node的奖品信息...")
	fmt.Println(runFuncName(), "3. 复写原所有奖品信息(抽取该node节点的奖品)...")
	return
}

// AmountDraw 具体抽奖行为
// 按数额范围区间抽奖 比如订单金额刮奖
type AmountDraw struct{}

// checkParams 其他参数校验(不同活动类型实现不同)
func (draw *AmountDraw) checkParams(*Context) (err error) {
	fmt.Println(runFuncName(), "按数额范围区间抽奖:特殊参数校验...")
	return
}

// getPrizesByNode 获取node奖品信息(不同活动类型实现不同)
func (draw *AmountDraw) getPrizesByNode(*Context) (err error) {
	fmt.Println(runFuncName(), "1. 判断属于哪个数额区间...")
	fmt.Println(runFuncName(), "2. 获取对应node的奖品信息...")
	fmt.Println(runFuncName(), "3. 复写原所有奖品信息(抽取该node节点的奖品)...")
	return
}

// Lottery 抽奖模板
type Lottery struct {
	// 不同抽奖类型的抽象行为
	concreteBehavior BehaviorInterface
}

// Run 抽奖算法
// 稳定不变的算法步骤
func (lottery *Lottery) Run(ctx *Context) (err error) {
	// 具体方法：校验活动编号(serial_no)是否存在、并获取活动信息
	if err = lottery.checkSerialNo(ctx); err != nil {
		return err
	}

	// 具体方法：校验活动、场次是否正在进行
	if err = lottery.checkStatus(); err != nil {
		return err
	}

	// ”抽象方法“：其他参数校验
	if err = lottery.checkParams(ctx); err != nil {
		return err
	}

	// 具体方法：活动抽奖次数校验(同时扣减)
	if err = lottery.checkTimesByAct(); err != nil {
		return err
	}

	// 具体方法：活动是否需要消费积分
	if err = lottery.consumePointsByAct(); err != nil {
		return err
	}

	// 具体方法：场次抽奖次数校验(同时扣减)
	if err = lottery.checkTimesBySession(); err != nil {
		return err
	}

	// 具体方法：获取场次奖品信息
	if err = lottery.getPrizesBySession(); err != nil {
		return err
	}

	// ”抽象方法“：获取node奖品信息
	if err = lottery.getPrizesByNode(ctx); err != nil {
		return err
	}

	// 具体方法：抽奖
	if err = lottery.drawPrizes(); err != nil {
		return err
	}

	// 具体方法：奖品数量判断
	if err = lottery.checkPrizesStock(); err != nil {
		return err
	}

	// 具体方法：组装奖品信息
	if err = lottery.packagePrizeInfo(); err != nil {
		return err
	}
	return
}

// checkSerialNo 校验活动编号(serial_no)是否存在
func (lottery *Lottery) checkSerialNo(ctx *Context) (err error) {
	fmt.Println(runFuncName(), "校验活动编号(serial_no)是否存在、并获取活动信息...")
	// 获取活动信息伪代码
	ctx.ActInfo = &ActInfo{
		// 假设当前的活动类型为按抽奖次数抽奖
		ActivityType: ConstActTypeTimes,
	}

	// 获取当前抽奖类型的具体行为
	switch ctx.ActInfo.ActivityType {
	case 1:
		// 按时间抽奖
		lottery.concreteBehavior = &TimeDraw{}
	case 2:
		// 按抽奖次数抽奖
		lottery.concreteBehavior = &TimesDraw{}
	case 3:
		// 按数额范围区间抽奖
		lottery.concreteBehavior = &AmountDraw{}
	default:
		return fmt.Errorf("不存在的活动类型")
	}
	return
}

// checkStatus 校验活动、场次是否正在进行
func (lottery *Lottery) checkStatus() (err error) {
	fmt.Println(runFuncName(), "校验活动、场次是否正在进行...")
	return
}

// checkParams 其他参数校验(不同活动类型实现不同)
// 不同场景变化的算法 转化为依赖抽象
func (lottery *Lottery) checkParams(ctx *Context) (err error) {
	// 实际依赖的接口的抽象方法
	return lottery.concreteBehavior.checkParams(ctx)
}

// checkTimesByAct 活动抽奖次数校验
func (lottery *Lottery) checkTimesByAct() (err error) {
	fmt.Println(runFuncName(), "活动抽奖次数校验...")
	return
}

// consumePointsByAct 活动是否需要消费积分
func (lottery *Lottery) consumePointsByAct() (err error) {
	fmt.Println(runFuncName(), "活动是否需要消费积分...")
	return
}

// checkTimesBySession 活动抽奖次数校验
func (lottery *Lottery) checkTimesBySession() (err error) {
	fmt.Println(runFuncName(), "活动抽奖次数校验...")
	return
}

// getPrizesBySession 获取场次奖品信息
func (lottery *Lottery) getPrizesBySession() (err error) {
	fmt.Println(runFuncName(), "获取场次奖品信息...")
	return
}

// getPrizesByNode 获取node奖品信息(不同活动类型实现不同)
// 不同场景变化的算法 转化为依赖抽象
func (lottery *Lottery) getPrizesByNode(ctx *Context) (err error) {
	// 实际依赖的接口的抽象方法
	return lottery.concreteBehavior.getPrizesByNode(ctx)
}

// drawPrizes 抽奖
func (lottery *Lottery) drawPrizes() (err error) {
	fmt.Println(runFuncName(), "抽奖...")
	return
}

// checkPrizesStock 奖品数量判断
func (lottery *Lottery) checkPrizesStock() (err error) {
	fmt.Println(runFuncName(), "奖品数量判断...")
	return
}

// packagePrizeInfo 组装奖品信息
func (lottery *Lottery) packagePrizeInfo() (err error) {
	fmt.Println(runFuncName(), "组装奖品信息...")
	return
}

func main() {
	(&Lottery{}).Run(&Context{})
}

// 获取正在运行的函数名
func runFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}
