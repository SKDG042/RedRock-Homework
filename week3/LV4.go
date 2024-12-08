package main

import (
	"fmt"
	"sync"
	"time"
)

var inventory = 50 //库存量
var mu sync.Mutex
var wg sync.WaitGroup

func seckill(ch chan int) {
	defer wg.Done() //确保seckill结束后计数器减一
	for range ch {
		mu.Lock() //加互斥锁
		if inventory > 0 {
			inventory--
			fmt.Println("秒杀成功!剩余库存:", inventory)
		} else {
			fmt.Println("库存不足")
		}
		time.Sleep(50 * time.Millisecond)
		mu.Unlock() //解锁
	}
}

func main() {
	ch := make(chan int, 10)

	for i := 0; i < 3; i++ { //开启i个goroutine模拟秒杀
		wg.Add(1) //每开启一个goroutine，计数器加一
		go seckill(ch)
	}

	for i := 0; i < 70; i++ { //模拟i个用户请求
		ch <- i
	}
	close(ch) //关闭通道防止死锁
	wg.Wait() //等待所有goroutine结束后退出
}
