package abase

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

// ??? 面向对象编程的成员方法访问在Go中是行不同的 ??? 只有 "包." 没有 "对象." ???

func TypeCheck(t interface{}) {
	// use of .(type) outside type switchgo
	// return t.(type)
}

//NOTE Go语言中的字符串
func TestString(t *testing.T) {
	// NOTE Go中的字符串是一个字节的切片。可以通过将其内容封装在 "" 中来创建字符串。
	// NOTE Go中的字符串是Unicode兼容的，并且是UTF-8编码的。
	println("---------------------------------------------------------")
	fmt.Println(strings.Count("dasdasdasd", "a"), strings.Count("dasdasdasd", ""))

	s := "Hello World"
	for i := 0; i < len(s); i++ {
		fmt.Printf("%d ", s[i])
	}
	// fmt.Printf("\n")
	println()
	for i := 0; i < len(s); i++ {
		fmt.Printf("%c", s[i])
	}
	println("\n---------------------------------------------------------")

	var str interface{} = "怎么#就切不#开呢？"
	switch strType := str.(type) {
	case string:
		subs := strings.Split(strType, "#")
		fmt.Println("长度：", len(subs))
	default:
		fmt.Println("类型：", strType)
	}

	// TODO RUNE type == int8
	s1 := 'a'
	s2 := "a"
	s3 := `a`
	fmt.Println(s1, s2, s3)

}

//NOTE Go语言中的数组
func TestArray(t *testing.T) {
	/*
		数组一旦被定义 则大小与元素类型不可改变, 一般数组元素的地址空间是连续的
		Go中的数组是值类型，而不是引用类型。这意味着当它们被分配给一个新变量时，
		将把原始数组的副本分配给新变量。如果对新变量进行了更改，则不会在原始数组中反映。
	*/

	// 1- 指定长度并完成了默认初始化
	var nums [3]int              // 指定类型，默认值 为 0
	var ages = [3]string{"a"}    // 不指定类型，则必须初始化, 默认值 ""
	names := [3]string{"2", "a"} // 省略 var, 不指定类型
	fmt.Println(nums, ages, names)

	// 2- 整体赋值，语法结构上携带类型[3]int （长度也是类型的一部分）
	nums = [3]int{1: 2}
	names = [3]string{0: "a", 2: "c"}
	fmt.Println(nums, names)

	// 3- 不指定长度，编译器自己根据初始化长度决定数组长度
	arr1 := [...]int{1, 2, 5: 7}
	arr1[2] = 3
	fmt.Println(arr1)
	// arr2 := [...]int{0: 1, 2, 1: 3}  // duplicate index in array literal: 1go

	// TODO 如何获取一个基础类型的长度 ？？？？
	fmt.Println(len(arr1), len("dasdsa"))

	// ------------------------------------------------------------
	// NOTE 使用反射库获取变量的类型
	tp := reflect.TypeOf(names[0])
	fmt.Println(tp)
	fmt.Println(reflect.TypeOf(tp))

	// NOTE 使用类型断言表达式判断变量的类型及其值的拷贝
	value, ok := interface{}(names[0]).(string)
	fmt.Println(value, ok)
	fmt.Println(&value, &names[0])

	// ------------------------------------------------------------
	// NOTE 不可变的空容器的内存地址都是同一个地址，这种容器存在没有实际的意义，
	//      为了节约内存所以采用同一个地址。 比如 map[string]struct{}.
	//      声明这样一个 map ??? 类型来标记某个 key 是否存在 ???. 在 key 值很多的情况下,
	//      要比 map[string]bool 之类的结构节约很多内存, 同时也减小 GC 压力.
	var (
		a [0]float32
		b struct{}
		c [0]struct {
			Value int64
		}
		d [10]struct{}
		e = new([10]struct{}) // new 返回的就是指针
		f byte
	)
	fmt.Printf("%p, %p, %p, %p, %p, %p", &a, &b, &c, &d, e, &f)

	// NOTE GO 语言中没有 char 类型，用 byte 代表
	var bt = [5]int{'A', 'a', '1'}
	fmt.Println(bt, reflect.TypeOf('1')) // int32

	// NOTE 和 python 一样可以通过  _ 忽略拆包中用不到的值
	for _, v := range bt { // => range(bt)
		println(v)
	}
	for i := 0; i < len(a); i++ {
		fmt.Printf("%d th element of a is %.2f\n", i, a[i])
	}

	// ------------------------------------------------------------
	// NOTE 多维数组的元素类型也都是一致的。
	var aa [3][4]int
	// illegal types for operand: print
	// println(aa)
	fmt.Println(aa)
	for i := 0; i < len(aa); i++ {
		for j := 0; j < len(aa[i]); j++ {
			fmt.Print(aa[i][j], "\t")
		}
		fmt.Println()
	}

	fmt.Println("---------------------")
	for _, arr := range aa {
		for _, val := range arr {
			fmt.Print(val, "\t")
		}
		fmt.Println()
	}

	// NOTE 数组的传递是值拷贝，且数组的长度是其类型的一部分
	aaa := [...]int{1, 2, 3}
	bbb := aaa
	bbb[1] = 0
	fmt.Println(&aaa == &bbb)

	// NOTE  如果要实现引用传递,则使用指针
	// invalid indirect of aaa (type [3]int)
	// TODO ccc := *aaa  // 指针的使用也有位置限制,只能作为函数参数 ????? func(arr *[3]int).

	// TODO 使用 GO 语言实现 二分查找 冒泡排序 快速排序 插入排序 选择排序 希尔排序 堆排序 归并排序
}

//NOTE Go语言中的切片,列表
func TestSlice(t *testing.T) {
	/**
	Go 语言切片是对数组的抽象封装。 Go 数组的长度不可改变，在特定场景中这样的集合不太适用。
	强悍的内置类型，切片("动态数组"), 应用而生,与数组相比切片的长度是不固定的，可对元素个数进行增减.

	切片是一种方便、灵活且强大的包装器。切片本身没有任何数据，它们只是对现有数组的引用。

	从概念上面来说slice像一个结构体，这个结构体包含了三个元素：
	type slice struct {
		array unsafe.Pointer // 被引用的数组中的起始元素地址, 指针，指向数组中slice指定的开始位置
		len   int            // 长度，即slice的长度
		cap   int            // 最大长度，也就是slice开始位置到数组的最后位置的长度
	}
	对 slice 的读写, 实际上操作的都是它所指向的数组.
	NOTE 切片作为函数参数的传递就是 引用传递,所谓引用传递,指的是局部变量和外部变量指向同一个编译器中可变的内存地址.
	NOTE 当多个切片共享相同的底层数组时，每个元素所做的更改将在数组中反映出来。
	*/

	// ??? 切片可以理解为是一个 特殊的指针
	ar := [...]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	// create a slice from array
	sl := ar[2:3:5]
	fmt.Println(reflect.TypeOf(ar), reflect.TypeOf(sl), sl[0])

	// NOTE 切片的长度是切片中元素的数量。切片的容量是从创建切片的索引开始的底层数组中元素的数量。
	// NOTE 访问切片时, index 值不能超过 len(slice)
	fmt.Println("cap-array:", cap(ar), "len-array:", len(ar))
	fmt.Println("cap-slice:", cap(sl), "len-slice:", len(sl))
	// NOTE 创建切片(slice[start:end])时, start 和 end 指定的区间不能超过 cap(slice) 范围

	println("-----------------------------------------------------")
	var slx []int
	// NOTE 默认值为 nil, 可见指针类型 和 slice 和 map 和 chan 都可以认为是 引用类型,
	// NOTE 换句话说 默认值 为nil的 类型都是引用类型
	var sl1 []int = nil
	// NOTE append 触发 slice 扩容 (想想也不是操作系统或者是编译器守护slice扩容的)
	sl2 := append(sl1, 1)
	fmt.Println(cap(sl1), len(sl1), cap(sl2), len(sl2), sl1, sl2, slx)
	println(&sl1, &sl2)

	println("-----------------------------------------------------")
	slice1 := []int{1, 2, 3, 4, 5}
	slice2 := append(slice1[:3], 6, 7)
	slice3 := append(slice1[:3], 6, 7, 8)
	fmt.Println(slice1, slice2, slice3)
	println(&slice1, slice1[0], &slice1[0])
	println(&slice2, slice2[0], &slice2[0])
	println(&slice3, slice3[0], &slice3[0])

	// append: 如果没超出,则在原来的slice上追加, 返回一个引用地址为追加之后的元素的起始地址
	//         如果超出,make一个长度为 最大长度*2 的 新slice,然后在拷贝,追加,原来的slice 没有变化.
	// NOTE: 无论超出还是没超出, slice 的本质是一个 结构体, 无论是切片操作还是 append 操作,都会生成一个新的
	// NOTE: slice 结构体变量,它们不是同一个结构体,内存地址不同, 但是,它们所封装的数组的起始地址可能是一样的.

	/*
		超出: 先make一个更长的slice,然后把整个slice都copy到新slice中, 再进行append.
		没超: 直接以len(slice)为起始点进行追加,len(slice)会随着append操作不断扩大,直到达到cap(slice)进行扩充.

		建议: 使用者尽可能的避免让 append 自动扩充内存. 一个是因为扩充时会出现一次内存拷贝, 二是因为 append 并不
		     知道需要扩充多少, 为了避免频繁扩充, 它会扩充到 2 * cap(slice) 长度. 而有时我们并不需要那么多内存.
			 所以在使用 slice 时, 最好不要不 make, 直接 append 让其自己扩充;
			 而是先 make([]int, 0, capValue) 准备一块内存, capValue 需要自己估计下, 尽可能确保足够用就好.
	*/
	println("-----------------------------------------------------")
	slice := make([]int, 5, 10)
	// slice[8]  越界
	println(len(slice), cap(slice)) // 长度为5, 容量为10
	newSlice := slice[:8]
	println(len(newSlice), cap(newSlice)) // 长度为8, 容量为10
	println("-----------------------------------------------------")
	// NOTE 将一个已存在的数组整体作为切片的底层数组
	numa := [3]int{78, 79, 80}
	nums1 := numa[:]
	nums1[0] = 100
	nums2 := numa[:]
	nums2[1] = 101
	fmt.Println(numa, nums1, nums2)
	println(&numa, &nums1, &nums2)
	println(&numa[0], &nums1[0], &nums2[0])
	println("-----------------------------------------------------")
}

//NOTE Go语言中的字典,映射
func TestMap(t *testing.T) {
	/**
	1- Map 是无序的，遍历map时无法决定它的返回顺序，因为Map是使用hash表来实现的.
	2- map是无序的，每次打印出来的map都会不一样，它不能通过index获取，而必须通过key获取
	3- map的长度是不固定的，也就是和slice一样，也是一种引用类型
	4- 内置的len函数同样适用于map，返回map拥有的key的数量
	5- key 的类型必须是 comparable 的
	*/

	var slice []int
	var ch chan int
	var dict map[int]string
	println(&slice, &ch, &dict)
	println("-------------------------------------------------------------------")
	// NOTE 可变长的空集合的地址是不一样的,不可变长的空容器的地址是一样的
	var arr = [...]int{}
	var stru struct{}
	println(&arr, &stru)
	var arr2 [0]int
	var stru2 struct{}
	println(&arr2, &stru2)
	println("-------------------------------------------------------------------")

	rating := map[string]float32{"C": 5, "Go": 4.5, "Python": 4.5, "C++": 2}
	fmt.Println(rating)
	var countryCapitalMap map[string]string
	print(countryCapitalMap)
	// countryCapitalMap["a"] = "lulu"
	// panic: assignment to entry in nil map [recovered]
	countryCapitalMap = make(map[string]string)

	countryCapitalMap["France"] = "Paris"
	countryCapitalMap["Italy"] = "Rome"
	countryCapitalMap["Japan"] = "Tokyo"
	countryCapitalMap["India"] = "New Delhi"

	// key 访问
	v := countryCapitalMap["Italy"]
	fmt.Println(v)
	v, ok := countryCapitalMap["Italy"]
	fmt.Println(v, ok)
	v, ok = countryCapitalMap["haha"]
	fmt.Println(v == "", ok) // NOTE 返回值类型的零值

	// exist 检测
	captial, ok := countryCapitalMap["United States"]
	if ok {
		fmt.Println("Capital of United States is", captial)
	} else {
		fmt.Println("Capital of United States is not present")
	}
	println("-------------------------------------------------------------------")

	// Key 遍历
	for country := range countryCapitalMap {
		fmt.Println("Capital of", country, "is", countryCapitalMap[country])
	}

	// Java Entry | Python Items 遍历
	for key, value := range countryCapitalMap {
		fmt.Println("key:", key, "value:", value)
	}

	// NOTE delete(map, key) 不返回任何值, 即使删除一个不存在的 key 也不会抛出异常
	delete(countryCapitalMap, "Japan")
	delete(countryCapitalMap, "lulu")
	fmt.Println(countryCapitalMap)
	println("-------------------------------------------------------------------")

	newCountryCapitalMap := countryCapitalMap
	fmt.Println(countryCapitalMap)
	newCountryCapitalMap["newInsert"] = "newInsert"
	fmt.Println(countryCapitalMap, newCountryCapitalMap)

	var newCountryCapitalMap2 map[string]string = countryCapitalMap

	// NOTE 全是值拷贝,只不过底层数据结构的,指针变量指向的 "内存值" 没变  // [key] 操作符重载
	println(&newCountryCapitalMap, &newCountryCapitalMap2, &countryCapitalMap)

	// NOTE map不能使用==操作符进行比较。==只能用来检查map是否为空。
	// 否则会报错：invalid operation: map1 == map2 (map can only be comparedto nil)
}

func TestMapPrinciple(t *testing.T) {
	/**
		数组：数组里的值指向一个链表
		链表：目的解决hash冲突的问题，并存放键值

	key
	|      key通过hash函数得到key的hash    |
	+------------------------------------+
	|       key的hash通过取模或者位操作     |
	|          得到key在数组上的索引        |
	+------------------------------------+
	|         通过索引找到对应的链表         |
	+------------------------------------+
	|       遍历链表对比key和目标key        |
	+------------------------------------+
	|              相等则返回value         |
	+------------------+-----------------+
	value
	*/

}
