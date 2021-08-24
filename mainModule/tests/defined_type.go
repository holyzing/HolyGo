package tests

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
	fmt.Printf("%d, %t, %d", sec.A, sec.b, sec.name)
}
