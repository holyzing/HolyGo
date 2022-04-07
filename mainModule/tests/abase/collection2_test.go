package abase

import (
	"fmt"
	"strings"
	"testing"
)

func TestMoreArray(t *testing.T) {
	// 数组是值类型，作为函数参数传递是值传递
	// 数组元素的类型和长度一起作为数组的类型（签名）
	// 多维数组中的元素类型是一样的
	// 内存连续，可快速迭代，数组类型信息提供每次迭代移位距离。
	var arr1 [0]int
	var arr2 = [...]int{}
	var arr3 = [0]int{}
	fmt.Println(arr1, arr2, arr3)

	// NOTE 数组中的 {} 意味着构造，其实和结构体在 ”表达形式“ 上有一定相似之处

	// var p *int
	// invalid memory address or nil pointer dereference
	// fmt.Println(*p)

	ar := [...]int{1, 2, 3, 4, 5}
	// 值拷贝
	ar1 := ar
	ar[1] = -1
	fmt.Println(ar, ar1)

	// go 中的序列（数组，字符串）并没有提供像python那么强大的 切片访问。
}

func TestMoreSlice(t *testing.T) {
	// 字面量初始化数组
	// sl1 和 sl2 共享底层数组
	arr := [...]int{1: 1, 3: 3, 5: 5, 6}
	sli := arr[:]  // 7, 7
	sl1 := arr[:2] // 2, 7
	sl1[1] = -1
	sl2 := arr[0:3] // 3, 7
	sl2[2] = -2
	println(&sli, &sl1, &sl2)
	fmt.Println(&sli[0], &sl1[0], &sl2[0])

	println("------------------------------------------")
	// 同一个数组切出来的切片，都引用该数组，所有切片的赋值引用操作都会反应在该数组上。
	sl3 := arr[0:2:3] // 2, 3
	sl4 := arr[1:2:4] // 1, 3
	sl5 := arr[1:3:4] // 2, 3
	// sl6 := arr[1:3:3]
	sl3[0] = 33
	sl4[0] = 44
	sl5[1] = 55
	// sl6 := arr[1:3:2] // invalid slice index: 3 > 2
	fmt.Println(sli, sl3, sl4, sl5, &sl3[0], &sl4[0], &sl5[0])

	println("------------------------------------------")

	// nil 切片也是可追加的,解释器可能做了特殊的处理
	var s1 []int
	var s2 = make([]int, 0)
	var s3 = []int{}

	print(len(s1), cap(s1))
	println(&s1, &s2, &s3)
	// append 可能产生新的切片，如果出现扩容，则底层引用数组会被改变
	s4 := append(s1, 0, 1, 2)
	// TODO copy 会按长度copy
	s5 := copy(s2, s4)

	fmt.Println(s1, s2, s4, s5)

	// 具体扩容策略可以参考源码 src/runtime/slice.go/growslice

	/*
		As for simple slice expressions, if a is a pointer to an array,
		a[low : high : max] is shorthand for (*a)[low : high : max]

		max：容量能达到最大的索引之处，所以容量就是 max - low
		   ：该切片一旦超出容量，则触发扩容，会基于一个新的底层数组生成新的切片，即使底层数组并没有越界
		   ：Additionally, it controls the resulting slice's capacity by setting it to max - low
		   ：如果不指定max，则容量为底层数组长度减去 low
		low：开始的索引
		high：结束的索引
	*/

	println("---------------------------------------------------------------------------------")
	arr = [...]int{1, 2, 3, 4, 5, 6, 7}
	ns1 := arr[1:2:3]
	fmt.Printf("%s len=%d cap=%d %v %X\n", "ns1", len(ns1), cap(ns1), ns1, &ns1[0])
	fmt.Println(arr, ns1)

	ns2 := append(ns1, -111)
	fmt.Printf("%s len=%d cap=%d %v %X\n", "ns2", len(ns2), cap(ns2), ns2, &ns1[0])
	fmt.Println(arr, ns1, ns2)

	ns3 := append(ns2, -222)
	fmt.Printf("%s len=%d cap=%d %v %X\n", "ns3", len(ns3), cap(ns3), ns3, &ns3[0])
	fmt.Println(arr, ns1, ns2, ns3)
	println("---------------------------------------------------------------------------------")

	ss1 := make([]int, 1, 3)
	fmt.Println(cap(ss1), len(ss1), ss1)

	func(sp []int) {
		a := append(sp, 1)
		fmt.Println(cap(a), len(a), a)
	}(ss1)
	fmt.Println(cap(ss1), len(ss1), ss1)
	fmt.Println(ss1[0:3])
}

func TestMoreMap(t *testing.T) {
	var m map[string]string = map[string]string{"a": "a"}
	fmt.Println("=====", m["b"] == "")
	val, ok := m["b"]
	if ok {
		fmt.Println(val)
	}
	fmt.Println(strings.EqualFold("a", "A"))
	// NOTE As a special case, it is legal to append a string to a byte slice, like this:
	// slice = append([]byte("hello "), "world"...)
	fmt.Println(([]byte)("a"))
}
