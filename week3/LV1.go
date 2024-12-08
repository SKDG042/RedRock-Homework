package main

import (
	"fmt"
	"sync"
)

func main() {
	var ch = make(chan string) //创建一个通道

	var wg sync.WaitGroup

	wg.Add(2) //等待两个goroutine完成

	go func() {
		defer wg.Done()
		for i := 1; i <= 100; i += 2 {
			ch <- "Red" //向通道发送数据同时阻塞所在函数
			fmt.Println(i)
			<-ch //当通道接收数据时继续循环
		}
	}()

	go func() {
		defer wg.Done()
		for i := 2; i <= 100; i += 2 {
			<-ch //当通道接收数据时执行
			fmt.Println(i)
			ch <- "Rock" //向通道发送数据同时阻塞所在函数
		}
	}()

	wg.Wait() //等待两个goroutine完成结束
}
