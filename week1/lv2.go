// 黄色警告我真处理不来（
package main

import (
	"fmt"
)

func main() {
	fmt.Println("欢迎使用Go语言计算器！")
	fmt.Println("请输入两个数字和一个操作符，进行四则运算")
	fmt.Println("输入exit退出程序")
	var a, b float64
	var x string
	for {
		fmt.Println("请输入第一个整数")
		if _, err := fmt.Scanf("%v", &a); err != nil {
			fmt.Println("输入错误，请输入一个整数")
			var discard string // 清理缓存，要不然会输出两次(
			fmt.Scanln(&discard)
			continue
		}
		fmt.Println("请输入操作符")
		fmt.Scanf("%s", &x)
		fmt.Println("请输入第二个整数")
		if _, err := fmt.Scanf("%v", &b); err != nil {
			fmt.Println("输入错误，请输入一个整数")
			var discard string
			fmt.Scanln(&discard)
			continue
		}
		switch x {
		case "+":
			fmt.Printf("%v + %v = %v\n", a, b, a+b)
		case "-":
			fmt.Printf("%v - %v = %v\n", a, b, a-b)
		case "*":
			fmt.Printf("%v * %v = %v\n", a, b, a*b)
		case "/":
			if b == 0 {
				fmt.Println("除数不能为0")
			} else {
				fmt.Printf("%v / %v = %v\n", a, b, a/b)
			}
		default:
			fmt.Println("操作符错误") // 试了半天只会在这里判断操作符错误
			continue
		}
		fmt.Println("是否继续?(exit退出)")
		var y string
		fmt.Scanf("%s", &y)
		if y == "exit" {
			fmt.Println("感谢使用！再见！")
			break
		}
	}
}
