package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type Event struct {
	Date time.Time
	Name string
}

func main() {
	file, err := os.Open("events.txt") //打开文件
	if err != nil {
		fmt.Println("无法打开该文件:", err)
		return
	}

	defer func() { //确保文件在函数结束时关闭
		if err := file.Close(); err != nil {
			fmt.Println("无法关闭该文件:", err)
		}
	}()

	var events []Event

	scanner := bufio.NewScanner(file) //逐行读取使用NewScanner
	for scanner.Scan() {
		line := scanner.Text()           //逐行读取
		part := strings.Split(line, " ") //以空格分割日期和时间
		date, err := time.Parse("2006-01-02", part[0])
		if err != nil {
			fmt.Println("无法正确转换日期:", err)
			continue
		}
		event := part[1]

		events = append(events, Event{Date: date, Name: event}) //将每一行的日期和事件添加到events中
	}

	now := time.Now()
	var closestEvent *Event
	var minDuration time.Duration

	for _, event := range events {
		duration := event.Date.Sub(now)
		if duration >= 0 && (closestEvent == nil || duration < minDuration) { //首先取第一个事件，然后找到最近的事件
			closestEvent = &event
			minDuration = duration
		}
	}
	days := int(minDuration / (time.Hour * 24)) //将时间转换为天数

	if closestEvent != nil {
		fmt.Printf("最近的一个事件是：%v\n还有%v天\n", closestEvent.Name, days)
	} else {
		fmt.Println("没有找到再未来发生事件")
	}

}
