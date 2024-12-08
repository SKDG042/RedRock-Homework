package main

import (
	"fmt"
	"time"
)

type Task interface { //定义一个名为 Task 的接口
	Execute() error //执行任务并用error返回错误信息（如果有）
}

type PrintTask struct { //创建一个PrintTask结构体用于打印消息
	Message string
}

type CalculationTask struct { //创建一个结构体用于计算两个数和
	A, B int
}

type SleepTask struct { //创建一个结构体用于使程序休眠指定的秒
	Duration time.Duration
}

func (m PrintTask) Execute() error { //实现结构体的Execute方法 打印Message
	fmt.Println(m.Message)
	return nil //返回nil表示没有错误
}

func (n CalculationTask) Execute() error { //计算并打印两个数的和
	fmt.Println(n.A + n.B)
	return nil
}

func (t SleepTask) Execute() error { //使程序休眠指定的秒数
	time.Sleep(t.Duration * time.Second)
	return nil
}

type Scheduler struct { //创造一个任务调度器
	Tasks []Task //用于存储所有任务
}

func (t *Scheduler) AddTask(task Task) { //添加任务
	t.Tasks = append(t.Tasks, task)
}

func (t *Scheduler) RunAll() { //遍历所有任务并执行
	for _, task := range t.Tasks {
		err := task.Execute()
		if err != nil { //错误处理
			fmt.Println("返回错误信息（如果有）", err)
		}
	}
}

func main() {
	scheduler := Scheduler{} //创建一个任务调度器

	scheduler.AddTask(PrintTask{Message: "灌注后端谢谢喵~"}) //添加任务
	scheduler.AddTask(CalculationTask{A: 0, B: 42})
	scheduler.AddTask(SleepTask{Duration: 2})

	scheduler.RunAll() //执行所有任务

	fmt.Println("给预习作业跪了") //打印GG 任务结束
}
