package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func main() {
	a := "hello"
	res := strings.HasPrefix(a, "e")
	fmt.Println(res)

	res = strings.HasSuffix(a, "o")
	fmt.Println(res)

	res = strings.ContainsAny(a, "asdo")
	fmt.Println(res)

	res = strings.Contains(a, "ll")
	fmt.Println(res)

	res2 := strings.Index(a, "ll")
	fmt.Println(res2)

	res2 = strings.LastIndex(a, "l")
	fmt.Println(res2)

	res2 = strings.LastIndexAny(a, "ol")
	fmt.Println(res2)

	str := "你好大家osgame大家们"
	old := "大家"
	new := "dajia"
	n := 2
	newStr := strings.Replace(str, old, new, n)
	fmt.Println(newStr)

	num := strings.Count(str, "大家")
	fmt.Println(num)

	str2 := "你好 osg"
	fmt.Println(len(str2), utf8.RuneCountInString(str2), len([]rune(str2)))

	fmt.Println(strings.Trim(str, "你好"))

	strSli := strings.Split(str2, "好")
	for _, v := range strSli {
		fmt.Println(v)
	}

	res3 := strings.Fields(str2)
	fmt.Println(res3, len(res3), strings.Join(res3, ","))
}
