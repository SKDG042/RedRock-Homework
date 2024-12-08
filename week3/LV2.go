package main

import (
	"fmt"
	"sync"
	"time"
)

// 定义全局变量
var ch chan int               // 用于goroutine间通信
var startTime time.Time       // 计时器开始时间
var elapsedTime time.Duration // 计时器已运行时间
var running bool              // 计时器是否在运行
var flags []time.Duration     // 存储flag时间
var wg sync.WaitGroup         // 等待所有goroutine结束后退出

// Timer 函数，接收通道中的命令并执行相应的操作
func Timer(ch chan int) {
	for {
		switch <-ch {
		case 0: // 重置计时器
			startTime = time.Time{}
			elapsedTime = 0
			running = false
			fmt.Println("计时器已重置")
		case 1: // 开始计时
			if !running {
				startTime = time.Now()
				running = true
				fmt.Println("计时器已开始")
			} else {
				fmt.Println("计时器已在运行")
			}
		case 2: // 暂停/继续计时
			if running {
				elapsedTime += time.Since(startTime)
				running = false
				fmt.Println("计时器已暂停, 已运行时间：", elapsedTime)
			} else {
				startTime = time.Now()
				running = true
				fmt.Println("计时器已继续")
			}
		case 3: // 设置flag
			if running {
				flagtime := time.Since(startTime) + elapsedTime
				flags = append(flags, flagtime)
				fmt.Println("创建了一个flag:", flagtime)
			} else {
				fmt.Println("计时器未在运行")
			}
		case 4: // 显示所有flag
			if len(flags) == 0 {
				fmt.Println("您还没有设置flag")
			} else {
				fmt.Println("所有flag：")
				for i, flag := range flags {
					fmt.Println("flag", i+1, ":", flag)
				}
			}
		case 5: // 清空所有flag
			flags = nil
			fmt.Println("已清空所有flag")

		case 6: // 退出计数器
			fmt.Println("计时器已退出")
			close(ch)
			wg.Done()
			return

		default:
			fmt.Println("输入错误,请输入一个1-5的数字")
			continue
		}
	}
}

// Input 函数，接收用户输入的命令并发送到通道
func Input(ch chan int) {
	for {
		var cmd int
		_, err := fmt.Scan(&cmd)
		if err != nil {
			fmt.Println("输入错误,请输入一个1-6的数字")
			continue
		}
		ch <- cmd
		if cmd == 5 {
			wg.Done()
		}
	}
}
func main() {
	fmt.Println("请输入命令：0-重置计时器，1-开始计时，2-暂停/继续计时,3-设置flag,4-显示所有flag,5-清空所有flag,6-退出计数器")
	// 初始化通道
	ch = make(chan int)
	wg.Add(2) //添加两个goroutine
	// 启动 Timer goroutine
	go Timer(ch)
	// 启动 Input goroutine
	go Input(ch)

	wg.Wait() // 等待所有goroutine结束后退出
}
