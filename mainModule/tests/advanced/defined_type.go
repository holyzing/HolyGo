package advanced

type More struct {
	A int
	b bool // 其它包无法直接访问
}

func (m *More) GetB() bool {
	return m.b
}
