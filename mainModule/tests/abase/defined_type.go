package abase

import (
	"fmt"
)

type name int8

type second struct {
	a int
	b bool
	name
}

type Second struct {
	A int
	b bool
	name
}

func Show(sec Second) {
	s := second{}
	fmt.Println(s.a, s.b, s.name)
	fmt.Printf("%d, %t, %d", sec.A, sec.b, sec.name)
}
