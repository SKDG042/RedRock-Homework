package main

import (
	"fmt"
)

type Student1 struct { //定义结构体
	Name  string
	Age   int
	Score float64
}

func main() {
	students := []Student1{ //创造实例
		{"程敬魏", 18, 20.09},
		{"晨月", 19, 11.10},
		{"白齐梁", 20, 23.00},
	}
	for _, student := range students { //遍历切片，打印信息
		fmt.Printf("姓名:%s \t年龄:%d\t成绩:%.2f\n", student.Name, student.Age, student.Score)
	}
}
