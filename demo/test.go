package main

import "fmt"

// 可比较的结构体
type ComparableStruct struct {
	a int
	b string
}

// 不可比较的结构体，包含不可比较的切片类型
type NonComparableStruct struct {
	c []int
	d map[string]int
}

func main() {

	// 示例：不可比较的类型
	fmt.Println("\n=== 不可比较的类型 ===")
	// 切片是不可比较的
	slice1 := [3]int{1, 2, 3}
	slice2 := [3]int{1, 2, 4}
	fmt.Println("slice1 == slice2:", slice1 == slice2) // 这行会导致编译错误

	// 不可比较的结构体，因为包含不可比较的切片类型
	//nonCompStruct1 := NonComparableStruct{c: []int{1, 2, 3}, d: map[string]int{"a": 1}}
	//nonCompStruct2 := NonComparableStruct{c: []int{1, 2, 3}, d: map[string]int{"a": 1}}
	//fmt.Println("nonCompStruct1 == nonCompStruct2:", nonCompStruct1 == nonCompStruct2) // 这行会导致编译错误
	//
	//cc := map[string]int{"a": 1}
	//dd := map[string]int{"a": 1}
	//fmt.Println(cc == dd)
}
