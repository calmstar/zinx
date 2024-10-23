package main

import "fmt"

func main() {
	s1 := &student{name: "zhangsan"}
	s2 := &student{name: "lisi"}
	s3 := &student{name: "wangwu"}

	fmt.Println(s1.name, s2.name, " 加入观察班长信号队列中")
	observerList := []observer{s1, s2}
	m := &monitor{observerList}
	fmt.Println(s3.name, " 加入观察班长信号队列中")
	m.add(s3)

	// 老师来了
	fmt.Println("老师来了")
	m.notify()

	fmt.Println("========")
	fmt.Println("踢出观察者：", s2.name)
	m.del(s2)
	fmt.Println("老师又来了")
	m.notify()
}

/**
触发：教师来了，
被观察者：班长，站岗看老师是否来了
观察者：学生，观察班长是否通知他们了
*/

// 抽象层 ----------
type observer interface {
	doSomething()
}

type beObserver interface {
	add(o observer)
	del(o observer)
	notify()
}

// 基础层----------
// 班长：充当被观察者（被学生观察的人）
type monitor struct {
	oList []observer
}

func (m *monitor) notify() {
	if m.oList == nil {
		return
	}

	for _, o := range m.oList {
		o.doSomething()
	}
}

func (m *monitor) add(o observer) {
	m.oList = append(m.oList, o)
}

func (m *monitor) del(o observer) {
	needDelIdx := -1
	for i, v := range m.oList {
		if v == o {
			needDelIdx = i
			break
		}
	}
	if needDelIdx != -1 {
		m.oList = append(m.oList[:needDelIdx], m.oList[needDelIdx+1:]...)
	}
}

// 学生：充当观察者，观察班长给出的通知信号，然后做点动作
type student struct {
	name string
}

func (s *student) doSomething() {
	fmt.Println(s.name, ", 停止嬉皮笑脸")
}
