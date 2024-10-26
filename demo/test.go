package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	const Max = 100000
	const NumReceivers = 10
	const NumSenders = 1000

	wg := sync.WaitGroup{}
	wg.Add(NumReceivers)

	// 1. 业务通道
	dataCh := make(chan int)

	// 2. 管理通道：必须是无缓冲通道
	// 其发送者是：额外启动的媒介协程
	// 其接收者是：业务通道的所有发送者和接收者
	stopCh := make(chan struct{})

	// 3. 媒介通道：必须是缓冲通道
	// 其发送者是：业务通道的所有接收者(实际情况是谁控制退出，就是谁)
	// 其接收者是：媒介协程（唯一）
	toStop := make(chan string, 1)

	var stoppedBy string

	// 媒介协程
	go func() {
		stoppedBy = <-toStop
		close(stopCh)
	}()

	// 业务通道发送者
	for i := 0; i < NumSenders; i++ {
		go func(id string) {
			for {
				// 提前检查管理通道是否关闭
				// 让业务通道发送者早尽量退出
				select {
				case <-stopCh:
					return
				default:
				}

				value := rand.Intn(Max)
				select {
				case <-stopCh:
					return
				case dataCh <- value:
				}
			}
		}(strconv.Itoa(i))
	}

	// 业务通道的接收者
	for i := 0; i < NumReceivers; i++ {
		go func(id string) {
			defer wg.Done()

			for {
				// 提前检查管理通道是否关闭
				// 让业务通道接收者早尽量退出
				select {
				case <-stopCh:
					return
				default:
				}

				select {
				case <-stopCh:
					return
				case value := <-dataCh:
					// 一旦满足某个条件，就通过媒介通道发消息给媒介协程
					// 以关闭管理通道的形式，广播给所有业务通道的协程退出
					if value == 6666 {
						// 务必使用 select，两个目的：
						// 1、防止协程阻塞
						// 2、防止向已关闭的通道发送数据导致panic
						select {
						case toStop <- "接收者#" + id:
						default:
						}
						return
					}

				}
			}
		}(strconv.Itoa(i))
	}

	wg.Wait()
	fmt.Println("被" + stoppedBy + "终止了")
}
