package advanced

import (
	"fmt"
	"mainModule/collects"
	"mainModule/tests"
	"testing"
)

func TestPrivateOrPublic(t *testing.T) {
	collects.Cao()
	// NOTE 只能访问到结构体的公有的属性
	se := tests.Second{}
	fmt.Println(se.A)
	m := More{A: 12, b: false}
	println(m.A, m.b)
}

func TestMakeNew(t *testing.T) {
	// make 的作用是初始化内置的数据结构，也就是我们在前面提到的切片、哈希表和 Channel；
	// new 的作用是根据传入的类型分配一片内存空间并返回指向这片内存空间的指针；
	// 从编译期间和运行时两个不同阶段理解这两个关键字的原理

	println("--------------------------------------------------------------------------------------")

}
