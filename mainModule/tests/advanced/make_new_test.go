package advanced

import (
	"fmt"
	"mainModule/tests"
	"testing"
)

func TestPrivateOrPublic(t *testing.T) {
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

	// ??? 所谓引用类型是指一般存储在堆上的,变量存储的是堆上地址的变量
	// Go: slice map chan ptr func interface{}

	ap := new([3]int)
	av := *ap // invalid operation: av == nil (mismatched types [3]int and nil)

	sp := new([]int)
	sv := *sp

	cp := new(chan int) // buffer 为0
	cv := *cp

	mp := new(map[string]int)
	mv := *mp // invalid argument mv (type map[string]int) for cap

	pp := new(*int)
	pv := *pp

	fmt.Println(av, ap, sp, sv, cp, cv, mp, mv, pp, pv)
	fmt.Println(sp == nil, sv == nil, cp == nil, cv == nil, mp == nil, mv == nil, pp == nil, pv == nil)
	fmt.Println(len(ap), cap(ap), len(av), cap(av))
	fmt.Println(len(sv), cap(sv), len(cv), cap(cv), len(mv))

	/**
	// NOTE: UNDERLYING TYPE runtime/*

	type slice struct {
		array unsafe.Pointer
		len   int
		cap   int
	}

	type hmap struct {
		// Note: the format of the hmap is also encoded in cmd/compile/internal/gc/reflect.go.
		// Make sure this stays in sync with the compiler's definition.
		count     int // # live cells == size of map.  Must be first (used by len() builtin)
		flags     uint8
		B         uint8  // log_2 of # of buckets (can hold up to loadFactor * 2^B items)
		noverflow uint16 // approximate number of overflow buckets; see incrnoverflow for details
		hash0     uint32 // hash seed

		buckets    unsafe.Pointer // array of 2^B Buckets. may be nil if count==0.
		oldbuckets unsafe.Pointer // previous bucket array of half the size, non-nil only when growing
		nevacuate  uintptr        // progress counter for evacuation (buckets less than this have been evacuated)
		extra *mapextra // optional fields
	}

	type hchan struct {
		qcount   uint           // total data in the queue
		dataqsiz uint           // size of the circular queue
		buf      unsafe.Pointer // points to an array of dataqsiz elements
		elemsize uint16
		closed   uint32
		elemtype *_type // element type
		sendx    uint   // send index
		recvx    uint   // receive index
		recvq    waitq  // list of recv waiters
		sendq    waitq  // list of send waiters

		// lock protects all fields in hchan, as well as several
		// fields in sudogs blocked on this channel.
		//
		// Do not change another G's status while holding this lock
		// (in particular, do not ready a G), as this can deadlock
		// with stack shrinking.
		lock mutex
	}
	*/

	println("-----------------------------------------------------------------------------------")

	// NOTE Go声明一个变量,即使不为其赋值,也会为为其开辟存储零值的内存,
	// 但是声明一个指针变量,其指向内存的零值为 nil

	var s1 []int
	var s2 = new([]int)
	var v2 = *s2
	var s3 *[]int

	fmt.Println(s1, s2, v2, s3, s1 == nil, s2 == nil, v2 == nil, s3 == nil, &s1)

	v1 := append(s1, 1)
	v2 = append(*s2, 1) // 返回的可能是经过扩容的,也可能是未经过扩容的
	println(v1[0], v2[0])

	println("-----------------------------------------------------------------------------------")

	chm := make(chan int)
	slm := make([]int, 0)
	mpm := make(map[string]int)

	fmt.Println(chm, chm == nil, slm, slm == nil, mpm, mpm == nil)
	fmt.Println(len(slm), cap(slm), len(mpm), len(chm), cap(chm))

	// TODO 引出堆栈分配
}
