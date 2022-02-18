package abase

// should not use dot imports (ST1001)go-staticcheck
import (
	"fmt"
	"reflect"
	"testing"
	"unicode/utf8"
)

var PERIOD = 1

var VERSION string = "0.0.0"

// syntax error: non-declaration statement outside function body
// hha := 10

// NOTE Go 语言中的变量
func TestVarGo(t *testing.T) {

	// NOTE 特殊变量 _
	// _, temp := 5, 6
	// cannot use _ as value, 已经被定义的 变量 _ 是一个具有只写属性的特殊变量,意味这抛弃,discarded
	// fmt.Println(_, temp)

	// NOTE 变量值交换
	a, b := 1, 2
	b, a = a, b
	p := &b // pointer

	println(&a, *p)

	var name string
	fmt.Println(name)
	name = "hello"
	fmt.Println(name)

	var age = 18
	var age2 int = 12
	fmt.Println(age, age2)

	// 只能用于函数体内的变量声明方式
	isExist := false
	fmt.Println(isExist)

	var t1, t2 int
	t1 = 1
	t2 = 2
	fmt.Println(t1, t2)
	t1, t2 = 3, 4
	fmt.Println(t1, t2)

	var a1, a2 = 1, "s"
	println(a1, a2)

	// no new variables on left side of :=
	// a1 := 2
	a1, a3 := 2, true
	print(a3)

	// 变量不能被重复定义
	var (
		v1 int = 2
		v2 string
	)
	v2 = "haha"
	println(v1, v2)
}

// NOTE Go语言中的常量
func TestConstGo(t *testing.T) {
	// const declaration cannot have type without expression
	// missing value in const declaration
	// const age int
	const age = 10
	const age2 int = 20
	println(age, age2)

	const a, b, c = 1, "a", true
	println(a, b, c)

	// 元组声明中只能声明一个类型，元组赋值可以赋不同类型的值
	// const a1 int, b1 string, c1 bool = 1, "a", true
	const a1, b1, c1 bool = true, false, true
	print(a1, b1, c1)

	const (
		// syntax error: unexpected :=, expecting =
		// v1 := 2
		v1 = 2
		v2
		v3      = "a"
		v4 bool = true
	)
	println(v1, v2, v3, v4)

	// NOTE 常量中的数据类型只可以是布尔型、数字型（整数型、浮点型和复数）和字符串型
	//	    不曾使用的常量，在编译的时候，是不会报错的
	// NOTE 数字常量不会分配存储空间，无须像变量那样通过内存寻址来取值，因此无法获取地址
	// ???  显示指定类型的时候，必须确保常量左右值类型一致，需要时可做显示类型转换。
	// ???  这与变量就不一样了，变量是可以是不同的类型值,会做隐式类型转换 ？？？？？

	// NOTE 特殊常量
	const (
		t1 = iota // 默认为 0
		t2        // 将上一行的表达式赋值给这一行
		t3 = "a"
		t4 = iota
		// 每当 iota 在新的一行被使用时，它的值都会自动加 1
	)
	println(t1, t2, t3, t4) // 0, 1, 2

	const (
		aa = iota //0
		bb        //1
		cc        //2
		dd = "ha" //独立值，iota += 1
		ee        //"ha"   iota += 1
		ff = 100  //iota +=1
		gg        //100  iota +=1
		hh = iota //7,恢复计数
		ii        //8
	)
	fmt.Println(aa, bb, cc, dd, ee, ff, gg, hh, ii)

	// ??? 如果中断iota自增，则必须显式恢复。且后续自增值按行序递增 ？？？
	// ??? 自增默认是int类型，可以自行进行显示指定类型 ？？？
}

// NOTE GO语言中的数据类型
func TestTypeGo(t *testing.T) {
	/**
	----
		bool

		NOTE 最高位表示符号
		int8   有符号 8  位整型 (-128 到 127)
		int16  有符号 16 位整型 (-32768 到 32767)
		int32  有符号 32 位整型 (-2147483648 到 2147483647)
		int64  有符号 64 位整型 (-9223372036854775808 到 9223372036854775807)

		NOTE 8位都用于表示数值
		uint8  无符号 8  位整型 (0 到 255)
		uint16 无符号 16 位整型 (0 到 65535)
		uint32 无符号 32 位整型 (0 到 4294967295)
		uint64 无符号 64 位整型 (0 到 18446744073709551615)

		NOTE int 和　uint　根据底层平台，表示32或64位整数。除非需要使用特定大小的整数，否则通常应该使用int来表示整数。

		float32: IEEE-754 32位浮点型数
		float64: IEEE-754 64位浮点型数
		complex64:  32 位实数和虚数
		complex128: 64 位实数和虚数

		byte: 类似 uint8
		rune: 类似 int32
		uint: 32 或 64 位
		int:  与 uint 一样大小

		uintptr: 无符号整型，用于存放一个指针

		string
			字符串就是一串固定长度的字符连接起来的字符序列。
			Go的字符串是由单个字节连接起来的。
			Go语言的字符串的字节使用UTF-8编码标识Unicode文本

		Type(Value)
		常数：在有需要的时候，会自动转型
		变量：需要手动转型 T(V)

		复合类型
		1、指针类型（Pointer）
		2、数组类型
		3、结构化类型(struct)
		4、Channel 类型
		5、函数类型
		6、切片类型
		7、接口类型（interface）
		8、Map 类型

	*/

	// TODO type rune = int32 ???

	// NOTE rune 是 int32的别名，几乎在所有方面等同于int32, 它用来区分字符值和整数值
	// NOTE byte 等同于int8，常用来处理ascii字符
	// NOTE rune 等同于int32,常用来处理unicode或utf-8字符
	r := []rune("dasdasd")
	fmt.Println(r, reflect.TypeOf(r))

	// 中文字符在unicode下占2个字节，在utf-8编码下占3个字节，而golang默认编码正好是utf-8
	var str = "hello 你好"
	//golang中string底层是通过byte数组实现的，len 函数返回的是字符串的字节数
	fmt.Println("len(str):", len(str))

	//golang中的unicode/utf8包提供了用utf-8获取长度的方法
	fmt.Println("RuneCountInString:", utf8.RuneCountInString(str))

	//通过rune类型处理unicode字符
	fmt.Println("rune:", len([]rune(str)))

	// TODO 类型转换测试

}

// NOTE Go语言中的操作符
func TestOperatorGo(t *testing.T) {
	// NOTE Go 语言中的操作符没有重载 比如 java, python中作为字符串的连字符,列表的拼接符

	// 算术运算符: + - * / %(求余) ++ --
	// 关系运算符: == != > < >= <=
	// 逻辑运算符: && || !
	// 位运算符: & | ^ &^ << >>
	// = += -= *= /= %= <<= >>= &= ^= |=

	// 优先级	   运算符
	//   7		~ ! ++ --
	//   6		* / % << >> & &^
	//   5		+ - ^
	//   4		== != < <= >= >
	//   3		<-
	//   2		&&
	//   1		||

	i := 1
	i++
	// NOTE 不支持 ++i
	// ++i

}

// NOTE Go语言的流程控制语句
func TestControlGo(t *testing.T) {
	// 普通 If 语句
	if 3 > 2 {
		print("3 > 2")
	} else if 1 < 2 {
		print("1 < 2")
	} else {
		print("没什么花样")
	}

	// if 语句变体 (num if语句的局部变量,只能在if语句及其分支中使用)
	if num := 10; num%2 == 0 {
		fmt.Println(num, "is even")
	} else {
		fmt.Println(num, "is odd")
	}
	// fmt.Println(num)

	// switch 语句, 匹配表达式的值,该值可以是任意类型的,
	// 但是用例类型必须与其保持一致,这也是顺理成章的.
	num := 10
	switch num {
	// case "a":
	// cannot use "a" (type untyped string) as type int
	// case 3 > 2:
	// cannot use 3 > 2 (type untyped bool) as type int
	case 1:
		fmt.Println("this is")
		fmt.Println("1")
		// break
		// redundant break statement (S1023)
	case 2, 3:
		fmt.Println("从 2,3 中匹配到")
	default:
		fmt.Println("没有匹配到任何值")
	}

	// 匹配 true
	switch false {
	// case 1:
	// cannot use 1 (type untyped int) as type bool
	case 2 > 3:
		fmt.Println("2 > 3")
	}

	// 默认匹配 true
	var grade = "dasd"
	switch {
	case grade == "A":
		fmt.Println("A")
	}

	// switch 变体
	switch x := 5; x {
	default:
		fmt.Println("default is a head and not break")
	// default:
	// multiple defaults in switch (first at ./init.go:271:2)
	// case 5:
	// duplicate case 5 in switch, case后的常量值不能重复, case后可以有多个常量值
	case 5:
		fmt.Println("匹配到了 x, 改变x, 让其继续向下匹配")
		x += 5
		fallthrough
	case 10:
		// fallthrough
		// fallthrough statement out of place
		fmt.Println("匹配到了 10")
		fallthrough
	case 11:
		fmt.Printf("匹配到了 11")
		// fallthrough
		// cannot fallthrough final case in switch
	}

	// NOTE switch 用于变量类型匹配
	var x interface {
	}
	switch i := x.(type) {
	case nil:
		fmt.Printf("type is %T", i)
	case int, int8:
		fmt.Printf("int or int8")
	case func(int):
		fmt.Println("func(int)")
	case func(string, int) int, string:
		fmt.Println("func(string, int) int, string")
	default:
		fmt.Println("未知型")
	}

	// use of .(type) outside type switch
	// if x.(type) == int {
	// }

	// var num = 10
	// cannot type switch on non-interface value num (type int)
	// switch num.(type) {
	// }
}

// NOTE Go语言的循环语句
func TestCircleGo(t *testing.T) {
	// Go 没有 while 语句
	// for init; condition; post { }
	for i := 10; i > 1; i -= 1 {
		fmt.Println(i)
	}

	// NOTE 相当于 while True
	for { // ==> for ;; {}
		fmt.Println("l")
		break
	}

	// expected boolean or range expression, found assignment (
	// missing parentheses around composite literal?)
	// NOTE bool 表达式不能被省略
	// for i := 3{
	// }
	// for i+=1{
	// }

	// NOTE 相当于 while
	for 5 > 3 {
		break
	}

	fmt.Println("---------------------------------")
	var a int
	//for 中的局部赋值表达式 所“定义”的变量是for 局部的
	for a := 5; a > 3; a -= 1 {
		fmt.Println(a)
	}
	fmt.Println(a)

	fmt.Println("---------------------------------")
	// if switch 局部代码块 也可改变外层变量的值。
	// NOTE C语言 的作用域边界是哪里定义哪里生效。
	for a = 5; a > 3; a -= 1 {
		fmt.Println(a)
	}
	fmt.Println(a)

	if true {
		g := true
		print(g)
	} else {
		// undefined: g
		// fmt.Println(g)
		fmt.Println("if 代码块 和 else 代码块也是不同的局部代码块")
	}
	// undefined: g
	// fmt.Println(g)

	fmt.Println("---------------------------------")
	// NOTE 迭代器
	// NOTE range 会根据所接的容器进行 一个变量的产出
	numbers := [6]int{1, 2, 3, 4, 5, 6}
	for i, num := range numbers {
		if i == 3 {
			continue
		}
		if i == 6 {
			break
		}
		fmt.Println("index:", i, "num:", num)
	}
}

// NOTE Go语言的goto
func TestGotoGo(t *testing.T) {
	// goto：可以无条件地转移到过程中指定的行。
	// NOTE：这逼玩意真难用，有啥用 ？？
	// NOTE: 如何跳出 lable ？？？
	// NOTE: 如何停止 lable 代码块的编写

	// i := 0
	// ll:
	// 	print("---------------")
	// 	for ; i < 5; i += 1 {
	// 		if i == 2 {
	// 			println(2)
	// 			goto ll
	// 		} else if i == 3 {
	// 			println(3)
	// 			break
	// 		} else {
	// 			println(i)
	// 		}
	// 	}
	// ------------------------------

	// 	err := firstCheckError()
	// 	if err != nil {
	// 		goto onExit
	// 	}
	// 	err = secondCheckError()
	// 	if err != nil {
	// 		goto onExit
	// 	}
	// 	fmt.Println("done")
	// 	return
	// onExit:
	// 	fmt.Println(err)
	// 	exitProcess()

}
