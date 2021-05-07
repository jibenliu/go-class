package main

// StateInfo 状态的基础信息和默认实现
type StateInfo struct {
	// 状态名
	name string
}

// Name 状态名
func (s *StateInfo) Name() string {
	return s.name
}

// 提供给内部设置名字
func (s *StateInfo) setName(name string) {
	s.name = name
}

// EnableSameTransit 允许同状态转移
func (s *StateInfo) EnableSameTransit() bool {
	return false
}

// OnBegin 默认将状态开启时实现
func (s *StateInfo) OnBegin() {

}

// OnEnd 默认将状态结束时实现
func (s *StateInfo) OnEnd() {

}

// CanTransitTo 默认可以转移到任何状态
func (s *StateInfo) CanTransitTo(name string) bool {
	return true
}
