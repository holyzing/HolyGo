package tests

import (
	"fmt"
	"reflect"
	"testing"
)

// NOTE redeclared in ehis block
// type name int8

type First struct {
	a int
	b bool
	name
}

// NOTE GO语言中的指针
func TestPointer(t *testing.T) {
	/**
	指针是存储另一个变量的内存地址的变量,变量是一种使用方便的占位符，用于引用计算机内存地址。

	& :取值符
	*/
	println("-------------------------------------------")
	var pointInt *int
	i := 1

	// NOTE cannot take the address of 1
	// pointInt = &1
	pointInt = &i

	fmt.Println(pointInt, *pointInt, reflect.TypeOf(pointInt))

	var a = First{1, false, 2}
	var b *First = &a
	var c = new(First)
	fmt.Println(a.a, a.b, a.name, c)
	println(&a, b, &b, b.a, (*b).a)

	// ??? 对于结构体变量来说  v.a <==> (*v).a ????

	println("-------------------------------------------")

	var j *int
	println(j)
	if j == nil {
		fmt.Println(j)
	}

	// NOTE Go不支持指针算法。

}

const MAX int = 3

func TestMorePointer(t *testing.T) {
	// THINK 变量的地址 和变量指向的地址 ???????  变量的内存和变量的地址 ??????
	b := 255
	a := &b // 指针变量 a 的(内存)值 是变量 b（内存）值的内存地址。
	*a++

	fmt.Println("address of b is", a)
	fmt.Println("value of b is", *a)

	f1 := func(a *int) {
		*a = 1
	}
	f1(a)
	fmt.Println(b)

	arr := [3]int{}
	f2 := func(s *[3]int) {
		s[0] = 1    // TODO 变量的内存 ？？？
		(*s)[1] = 2 // TODO 取变量 ？？？
	}
	f2(&arr)
	fmt.Println(arr)

	// 虽然将指针传递给一个数组作为参数的函数并对其进行修改，但这并不是实现这一目标的惯用方法,
	// 一般会使用切片,因为不同的切片的底层数组可以是同一数组，尽管它们是 值传递

	// Go语言中的指针不支持 指针运算
	// const MIN := 3  // NOTE 常量不支持 := 赋值

	arr2 := []int{10, 100, 200}
	var i int
	var ptr [MAX]*int // NOTE 指针数组

	for i = 0; i < MAX; i++ {
		ptr[i] = &arr2[i] /* 整数地址赋值给指针数组 */
	}

	for i = 0; i < MAX; i++ {
		fmt.Printf("a[%d] = %d\n", i, *ptr[i])
	}
	// ??? 指针的指针 ？？？
}
