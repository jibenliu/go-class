package example

import "fmt"

//实现一个日志记录器，（相当于context）
type LogManager struct {
	Logging
}

func NewLogManager(logging Logging) *LogManager {
	return &LogManager{logging}
}

//抽象的日志
type Logging interface {
	Info()
	Error()
}

//实现具体的日志：文件方式日志
type FileLogging struct {
}

func (fl *FileLogging) Info() {
	fmt.Println("文件记录Info")
}
func (fl *FileLogging) Error() {
	fmt.Println("文件记录Error")
}

// 实现具体的日志：数据库方式的日志

type DbLogging struct {
}

func (dl *DbLogging) Info() {
	fmt.Println("DB记录Info")
}
func (dl *DbLogging) Error() {
	fmt.Println("DB记录Error")
}
