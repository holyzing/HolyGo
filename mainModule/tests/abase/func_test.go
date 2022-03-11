package abase

import (
	"fmt"
	"reflect"
	"testing"
)

// NOTE 返回类型中带变量名的函数
func test1(x, y int) (add, multi int) {
	add = x + y
	multi = x * y
	return // a 可省略, 返回的依旧是 a
}

// NOTE 可变长参数
func test2(a ...int) string {
	fmt.Println(reflect.TypeOf(a)) // []int
	return "successful"
}

// NOTE 值传递和引用传递
func test3() {
	// NOTE 值传递:
	// NOTE 非复合类型和数组都是值传递
	// NOTE 切片和数组其实是值传递,但是其底层的数组和哈希的地址是一样的,所以可以认为是引用传递

	slice := []int{1, 2}
	println(&slice, &slice[0])
	funcValue := func(x []int) {
		println(&x, &x[0])
	}
	funcValue(slice)

	println("--------------------------------------")
	// NOTE 引用传递:
	// NOTE 直接传递变量的指针类型, 指针类型的值就是一字节的内存地址, 但是传递的依旧是 指针类型变量的地址.
	// NOTE 这样做的一大好处是,直接操作指向的内存而不是拷贝一份后传递,降低内存存储以及拷贝等的资源消耗.
	//      传指针使得多个函数能操作同一个对象。
	//	    传指针比较轻量级 (8bytes),只是传内存地址，我们可以用指针传递体积大的结构体。
	//      如果用参数值传递的话, 在每次copy上面就会花费相对较多的系统开销（内存和时间）。
	// ???  注：若函数需改变slice的长度，则仍需要取地址传递指针 ???
	i := 1
	ip := &i
	println("outer:", i, ip, &ip)

	funcPointer := func(i *int) {
		// NOTE 显然是无法改变变量所指向的内存的地址
		println("inner:", *i, i, &i)
		*i = 2
	}
	funcPointer(ip)
	println("outer:", i, ip, &ip)
}

// expected declaration, found println
// println("---------------------------------------------------------------------------")

//NOTE Go语言中的函数
func TestFunc(t *testing.T) {
	/*
		函数名称: 函数名和参数列表一起构成了函数签名。
		参数列表:
			参数列表指定的是参数类型、顺序、及参数个数。
			形式参数：定义函数时，用于接收外部传入的数据，叫做形式参数，简称形参
			实际参数：调用函数时，传给形参的实际的数据，叫做实际参数，简称实参。

		函数调用：
			A：函数名称必须匹配
			B：实参与形参必须一一对应：顺序，个数，类型

		返回类型:
			函数返回一列值。return_types 是该列值的数据类型。
			有些功能不需要返回值，这种情况下 return_types 不是必须的。
			上面返回值声明了两个变量output1和output2，如果你不想声明也可以，直接就两个类型。
			如果只有一个返回值且不声明返回值变量，那么可以省略包括返回值的括号
			没有返回类型声明, "最后的" return 可省略, redundant return statement

		变量作用域: 所有代码块 "{...}" 内定义的变量,只能在代码块中起作用

	*/
	a, m := test1(1, 2)
	println(a, m)
	test2(1)
	test3()
}

// NOTE 这是一个类
type Person struct {
	name string
	age  int
}

// NOTE 这是一个方法, 函数中不能定义方法 ???
func (p Person) eat(food string) {
	fmt.Println("name:", p.name, "age:", p.age, "eating food:", food)
}

func redress(l int, err error) {
	fmt.Println("参数输出的长度为", l)
}

// NOTE Go 语言中的延迟
func TestDefer(t *testing.T) {
	// NOTE 延迟到某个动作(比如 return) 后执行.
	/*
		可以在函数中添加多个defer语句。当函数执行到最后时，这些defer语句会按照逆序执行，最后该函数返回。
		当在进行一些打开资源的操作时，遇到错误需要提前返回，返回前需要关闭资源，不然造成资源泄露等问题。

		1- 一个函数体内的多个 defer 的调用,是按照 先进后出 即入栈出栈 的顺序.
		2- 当函数执行中 遇到return 或者 出现未处理跑出异常的时候,开始执行 defer
	*/

	// NOTE IO 操作最后才会执行 资源关闭等操作

	a, b := 1, "a"
	defer println(a)
	println(b)

	p := Person{
		name: "ll",
		age:  10,
	}
	defer p.eat("apple")

	fmt.Println("Show Person")

	c, ok := fmt.Println("延迟执行函数的参数")

	var d int

	var e int
	fmt.Println(e)
	e = 2
	e = 4
	defer redress(c, ok)
	defer redress(d, ok)
	defer redress(e, ok)

	c = 0
	d = 1
	e = 3
	// fmt.Println(ok)
	// NOTE:1 在延迟函数执行之前且作为延迟函数的参数 的变量的初始化赋值操作(如果没有则是默认初始化),
	//      会在执行延迟函数之前完成赋值操作,但是赋值的表达式是会提前执行的.
	// NOTE:2 其实上边的NOTE:1解释是错误的,正确的是,defer 会记录延迟执行时作为函数参数的变量的值.

	fmt.Println("-------------------------------------------", c, d, e)

	// NOTE 原理: 当一个函数有多个延迟调用时，多个延迟调用被添加到一个堆栈中，
	//           并最后按照 Last In First Out（LIFO）的顺序进行出栈

	name := "Naveen"
	fmt.Printf("Orignal String: %s\n", string(name))
	fmt.Printf("Reversed String: ")
	for _, v := range name {
		defer fmt.Printf("%c", v)
	}
}

func func4(i int, fc func(int, int) (int, string)) (int, string) {
	// THINK 函数式编程，在函数内部是不能声明函数的 ？？？
	// func innerFunc(){
	// }
	// NOTE 所谓高阶函数和回调函数
	return fc(i, 1)
}

// NOTE Go 语言中的高阶函数
func TestHighFunc(t *testing.T) {
	// 函数名，指向函数体的内存地址，而且Go语言是支持函数式编程的。
	// 1.将匿名函数作为另一个函数的参数，回调函数
	// 2.将匿名函数作为另一个函数的返回值，可以形成闭包结构。

	var i int = func(i int) int {
		fmt.Println("匿名函数的一次调用", i)
		return i
	}(5)

	func1 := func(name string) int {
		fmt.Println("匿名函数赋值给变量", name)
		return i
	}
	func1("自律")

	func2 := func(a, b int) (int, string) {
		return 3, "a"
	}

	var func3 func(int, int) (int, string) = func(a, b int) (int, string) {
		return 3, "a"
	}

	func2(1, 2)
	func3(1, 2)
	func4(1, func3)
}

// NOTE Go 语言中的闭包
func TestClosure(t *testing.T) {
	// 一个外层函数中，有内层函数，该内层函数中，会操作外层函数的局部变量(
	// 外层函数中的参数，或者外层函数中直接定义的变量)，并且该外层函数的返回值就是这个内层函数。
	// 这个内层函数和外层函数的局部变量，统称为闭包结构

	// 闭包结构：局部变量的生命周期会发生改变，正常的局部变量随着函数调用而创建，随着函数的结束而销毁。
	// 但是闭包结构中的外层函数的局部变量并不会随着外层函数的结束而销毁，因为内层函数还要继续使用。

	increment := func() func() int { //外层函数
		//1.定义了一个局部变量
		i := 0
		//2.定义了一个匿名函数，给变量自增并返回
		fun := func() int { //内层函数
			i++
			return i
		}
		//3.返回该匿名函数
		return fun
	}

	res1 := increment()      //res1 = fun
	fmt.Printf("%T\n", res1) //func() int
	fmt.Println(&res1)
	println(&res1) // ??? 函数变量的地址 ？
	println(res1)  // ??? 函数体的地址 ？
	println("-----------------------------------------")
	v1 := res1()
	fmt.Println(v1) //1
	v2 := res1()
	fmt.Println(v2) //2
	fmt.Println(res1())
	fmt.Println(res1())
	fmt.Println(res1())
	fmt.Println(res1())
	println("-----------------------------------------")

	res2 := increment()
	fmt.Println(&res2)
	fmt.Println(res2())
	fmt.Println(res2())

	fmt.Println(res1())

}

// TODO 注意区分声明与定义，在Go语言中，函数内是不能声明函数的 ？？？？只能在函数能将一个函数定义为变量
// TODO 线上出现Bug 之后如何定位Bug ？？？
// https://zhuanlan.zhihu.com/p/159135741

// -------------------------------------------------------------------------------------------------
// Go-Kit EndPoint 闭包

func aaa() (done func(), err error) {
	return func() { print("aaa: done") }, nil
}

func bbb() (done func(), _ error) {
	done, err := aaa()
	return func() { print("bbb: surprise!"); done() }, err
}

func TestClosure2(t *testing.T) {
	done, _ := bbb()
	done()
}
