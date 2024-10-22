package main

import "fmt"

func main() {
	adapter := newAdapter(v220{})
	p := &phone{v5: adapter}
	p.charging()
}

type v5 interface {
	use5v()
}

// -------
type phone struct {
	v5 v5
}

func newPhone(v5 v5) *phone {
	return &phone{v5: v5}
}
func (p *phone) charging() {
	p.v5.use5v()
}

// -------
type adapter struct {
	v220 v220
}

func newAdapter(v220 v220) *adapter {
	return &adapter{v220: v220}
}

func (a *adapter) use5v() {
	// 转换操作
	fmt.Println("开始执行转换操作，将220v转成5v")
	a.v220.Use220v()
}

// -------
type v220 struct {
}

func (v *v220) Use220v() {
	fmt.Println("使用220v")
}
