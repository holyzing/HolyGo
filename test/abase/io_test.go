package abase

import (
	"fmt"
	"testing"
)

// NOTE GO语言中的IO
func TestIOGo(t *testing.T) {
	// ??? 无法做到 在声明时，就做到元组的赋值 ???
	// var num int
	// var err error
	num, err := fmt.Print("asas")
	println(num, err)

	// fmt包实现了类似C语言printf和scanf的格式化I/O

	// Print call has possible formatting directive %d
	// fmt.Print("num %d", 5)

	fmt.Printf("num %d", 5)
	// 去官网查询 格式化的占位符
	// https://golang.google.cn/pkg/fmt/
	// https://golang.google.cn/pkg/bufio/

}
