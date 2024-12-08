// 我的理解:生产者和消费者之间使用一个通道,生产者产生的数据通过通道传给消费者，消费者从通道中取出数据进行处理
// 但生产者不需要等待消费者处理完数据即可继续生产,多产出的数据会被缓存到通道中
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var ch = make(chan int, 10)
	var wg sync.WaitGroup

	wg.Add(2)
	fmt.Println("抽取你的角色")

	go func() {
		defer wg.Done()
		for i := 1; i <= 10; i++ {
			fmt.Println("抽取中...", i)
			time.Sleep(20 * time.Millisecond)
			ch <- i
		}
		close(ch)
	}()

	go func() {
		defer wg.Done()
		time.Sleep(50 * time.Millisecond)
		for num := range ch {
			switch num {
			case 1:
				fmt.Printf("TOP%d:永雏塔菲❤\n", num)
			case 2:
				fmt.Printf("TOP%d:冰糖IO❤\n", num)
			case 3:
				fmt.Printf("TOP%d:東雪莲❤\n", num)
			case 4:
				fmt.Printf("TOP%d:七海❤\n", num)
			case 5:
				fmt.Printf("TOP%d:库莉姆❤\n", num)
			case 6:
				fmt.Printf("TOP%d:阿梓从小就很可爱❤\n", num)
			case 7:
				fmt.Printf("TOP%d:尼奈❤\n", num)
			case 8:
				fmt.Printf("TOP%d:伽乐❤\n", num)
			case 9:
				fmt.Printf("TOP%d:圣嘉然❤\n", num)
			case 10:
				fmt.Println("排行结束")
			}
		}
	}()

	wg.Wait()
}
