package advanced

import (
	"fmt"
	"runtime"
	"testing"
)

// NOTE 逃逸
func TestMemoryEscape(t *testing.T) {
	fmt.Println(runtime.NumCPU())
	i := 1
	i++
	i--
}

// 函数内部的变量一般是分配在栈中的，函数返回的局部变量是局部变量的拷贝,函数结束后，栈空间被收回
// 当函数内部返回一个局部变量的地址，则当函数结束后，局部变量被回收，在函数外部访问变量的地址，会出现段错误，
// 此时局部变量需要定义在堆中。

// 但是在go语言中，函数是应该定义在堆上还是栈中，由编译器来进行逃逸分析后决定。
