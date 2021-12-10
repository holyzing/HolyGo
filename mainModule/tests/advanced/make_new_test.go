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

	// 二者均是内建函数，不属于关键字
	// make 创建内建的引用类型，并返回一个类型实例，而不是类型实例的指针。

	// slice 是一个包含 data、cap 和 len 的结构体 reflect.SliceHeader；
	// hash 是一个指向 runtime.hmap 结构体的指针；
	// ch 是一个指向 runtime.hchan 结构体的指针；

	// NOTE As for simple slice expressions, if a is a pointer to an array,
	// a[low : high : max] is shorthand for (*a)[low : high : max]
	// 但是 slice 和 map 和chan 则不能通过指针直接访问

	chn := new(chan int) // buffer 为0
	sln := new([]int)
	mpn := new(map[string]int)

	fmt.Println(*chn == nil, *sln == nil, *mpn == nil)
	fmt.Println(chn, *chn, sln, *sln, mpn, *mpn)
	fmt.Println(len(*sln), cap(*sln), len(*mpn), len(*chn))

	*sln = append(*sln, 1)
	// (*mpn)["a"] = 2  // assignment to entry in nil map
	println(sln, *sln)
	println("-----------------------------------------------------------------------------------")

	chm := make(chan int)
	slm := make([]int, 0)
	mpm := make(map[string]int)

	fmt.Println(chm, slm, mpm)
	fmt.Println(len(slm), cap(slm), len(mpm), len(chm), cap(chm))

	println("-----------------------------------------------------------------------------------")

	// TODO 引出堆栈分配
}
