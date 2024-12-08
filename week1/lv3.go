// 罗马数字转整数，不是整数转罗马！！！
package main

import (
	"fmt"
)

func romanToInt(s string) int { //定义函数计算罗马数字
	roman := map[string]int{ //定义哈希表储存罗马数字
		"I": 1,
		"V": 5,
		"X": 10,
		"L": 50,
		"C": 100,
		"D": 500,
		"M": 1000,
	}
	result := 0
	n := len(s)              //计算输入的罗马数字长度
	for i := 0; i < n; i++ { //逐个计算每位罗马数字
		if i < n-1 && roman[string(s[i])] < roman[string(s[i+1])] { //if确保i不是最后一位数字，&&判断当前位罗马数字是否小于下一位
			result -= roman[string(s[i])] //如果小于则减去当前位数字
		} else {
			result += roman[string(s[i])] //如果大于或等于则加上当前位数字
		}
	}
	return result //返回计算结果
}

func main() {
	x := ""
	for {
		fmt.Println("请输入一个罗马数字")
		fmt.Scan(&x) //跟lv2一样不会在这判断(
		y := romanToInt(x)
		if y > 0 { //判断输入是否错误
			fmt.Println(y)
			break
		} else {
			fmt.Println("输入错误,请重新输入")
		}

	}

}
