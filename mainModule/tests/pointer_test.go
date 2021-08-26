package tests

import (
	"fmt"
	"reflect"
	"testing"
)

// NOTE redeclared in this block
// type name int8

// NOTE var p *string
// NOTE 指针变量的值：指针指向的变量的内存地址,即 p,类型为 int
// NOTE 指针地址的值：即指针变量自己本身的内存地址， 即 &p 类型为 int
// NOTE 指针指向的值：指的是P值代表的地址上存储的值，即 *p 类型为 string
// NOTE 在语义层面上指针变量是存放其它变量（变量的地址）的一种特殊的变量，
//      解释器会对 *操作 做取出它指向的地址的值 的操作的特殊处理。

// THINK 变量的值和变量指向的地址值 ？？？
// THINK slice 底层是一个指向数组的指针变量 ？？？
// THINK 变量是一种引用？引用某个地址的值 ？？？

/**
一、指针和引用的区别

(1)引用总是指向一个对象,没有所谓的 null reference .所有当有可能指向一个对象也有可能不指向对象则必须使用 指针.
   由于C++ 要求 reference 总是指向一个对象所以 reference要求有初值.

   String & rs = string1;
   由于没有所谓的 null reference 所以在使用前不需要进行测试其是否有值,而使用指针则需要测试其的有效性.

(2)指针可以被重新赋值而reference则总是指向最初或地的对象.
(3)必须使用reference的场合. Operator[] 操作符 由于该操作符很特别地必须返回
   [能够被当做assignment 赋值对象] 的东西,所以需要给他返回一个 reference.
(4)其实引用在函数的参数中经常使用.
*/

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
