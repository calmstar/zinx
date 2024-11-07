package main

import (
	"fmt"
	"time"
)

func main() {
	ch1 := make(chan int)

	go func() {
		time.Sleep(3 * time.Second)
		ch1 <- 1
	}()

	select {
	case <-time.After(2 * time.Second):
		fmt.Println("时间过了2秒，超时了")
	case data := <-ch1:
		fmt.Println("从ch1中读出数据, data:", data)
	}

}
