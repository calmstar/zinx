package main

import (
	"bytes"
	"fmt"
)

func main() {
	buff := bytes.NewBuffer([]byte{})
	cnt, err := buff.Write([]byte("hhh"))
	buff.Write([]byte(" nihao "))
	fmt.Println(cnt, err, buff)

	buff.WriteString(" world")
	fmt.Println(string(buff.Bytes()))

	var b byte = '!'
	buff.WriteByte(b)
	fmt.Println(buff.String())

	var s rune = '中'
	buff.WriteRune(s)
	fmt.Println(buff.String())

	//f, _ := os.Create("test.log")
	//writeCnt, err := buff.WriteTo(f)
	//fmt.Println(writeCnt, err)

	read3 := make([]byte, 3)
	readCnt, err := buff.Read(read3)
	fmt.Println(err, readCnt, string(read3), buff.String())

	b2 := bytes.NewBufferString("")
	b2.WriteByte('a')
	b2.WriteString("bcdefg")
	read1, err := b2.ReadByte()
	fmt.Println(string(read1), b2)

}
