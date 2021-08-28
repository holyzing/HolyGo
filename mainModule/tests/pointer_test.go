package tests

import (
	"fmt"
	"reflect"
	"testing"
)

// NOTE redeclared in this block, previous declaration at ./defined_type.go
// type name int8

// NOTE var p *string
// NOTE 指针变量的值：指针指向的变量的内存地址,即 p,类型为 int
// NOTE 指针的地址值：即指针变量自己本身的内存地址， 即 &p 类型为 int
// NOTE 指针指向的值：指的是P值代表的地址上存储的值，即 *p 类型为 string
// NOTE 在语义层面上指针变量是存放其它变量（变量的地址）的一种特殊的变量，
//      编译器会对 *操作 做取出它指向的地址的值 的操作的特殊处理。

// THINK Go语言的内存模型 ？？？？？？？

/**
一、指针和引用的区别

指针（pointer）在Go语言中可以被拆分为两个核心概念：

类型指针：允许对这个指针类型的数据进行修改，传递数据可以直接使用指针，而无须拷贝数据，类型指针不能进行偏移和运算。
	    受益于这样的约束和拆分，Go语言的指针类型变量即拥有指针高效访问的特点，又不会发生指针偏移，
	    从而避免了非法修改关键性数据的问题。同时，垃圾回收也比较容易对不会发生偏移的指针进行检索和回收。
切片：由指向起始元素的原始指针、元素数量和容量组成。
     切片比原始指针具备更强大的特性，而且更为安全。切片在发生越界时，运行时会报出宕机，并打出堆栈，而原始指针只会崩溃。


(1) 引用总是指向一个对象,没有所谓的 null reference .当有可能指向一个对象也有可能不指向对象则必须使用 指针.
    由于C++ 要求 reference 总是指向一个对象所以 reference要求有初值.

    String & rs = string1;
    由于没有所谓的 null reference 所以在使用前不需要进行测试其是否有值,而使用指针则需要测试其的有效性.

(2) 指针可以被重新赋值而reference则总是指向最初或地的对象.
(3) 必须使用reference的场合. Operator[] 操作符 由于该操作符很特别地必须返回
    [能够被当做assignment 赋值对象] 的东西,所以需要给他返回一个 reference.
(4) 其实引用在函数的参数中经常使用.

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

	fmt.Printf("%p\n", pointInt)
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

	// NOTE 使用 new 创建的类型实例返回的都是一个指向实例的指针
	firstPtr := new(First)
	fmt.Printf("%T\n", firstPtr)
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
