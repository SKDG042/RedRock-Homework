package main

import (
	"fmt"
)

type Student struct { //定义结构体 student
	Name  string
	Age   int
	score []int
}

type Classroom struct { //定义结构体 class
	Classname string
	Students  []*Student
}

func AddStudent(c *Classroom, s *Student) { //编写函数添加学生信息
	c.Students = append(c.Students, s)
}

func UpdateScore(s *Student, score int) { //编写函数追加成绩
	s.score = append(s.score, score)
}

func CalculateAverage(s *Student) float64 { //编写函数计算平均成绩
	var sum int
	for _, s := range s.score { //遍历计算总分
		sum += s
	}
	var average float64
	average = float64(sum) / float64(len(s.score)) //计算平均分
	return average
}

func main() {
	class := Classroom{Classname: "我们后端都是香香软软的小可耐"} //创建亲亲爱爱的后端班级
	student1 := &Student{Name: "程敬魏", Age: 18}
	student2 := &Student{Name: "晨月", Age: 19}
	student3 := &Student{Name: "白齐梁", Age: 20}

	AddStudent(&class, student1)
	AddStudent(&class, student2)
	AddStudent(&class, student3)

	UpdateScore(student1, 18)
	UpdateScore(student1, 20)
	UpdateScore(student1, 9)
	UpdateScore(student2, 19)
	UpdateScore(student2, 11)
	UpdateScore(student2, 10)
	UpdateScore(student3, 20)
	UpdateScore(student3, 23)
	UpdateScore(student3, 0)

	fmt.Println(class.Classname)
	for _, s := range class.Students { //遍历打印平均成绩
		average := CalculateAverage(s)
		fmt.Printf("姓名:%s \t年龄:%d\t平均成绩:%.2f\n", s.Name, s.Age, average)
	}
}
