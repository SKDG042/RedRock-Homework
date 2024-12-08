package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Author struct {
	Name string
	Bio  string
}

type Post struct {
	Title   string
	Content string
	Author  Author
	Tags    []string
}

func main() {
	author := Author{
		Name: "042",
		Bio:  "215375205",
	}

	post := Post{
		Title:   "加入后端谢谢喵~",
		Content: "灌注后端谢谢喵~",
		Author:  author,
		Tags:    []string{"后端", "Golang"},
	}

	postJson, err := json.MarshalIndent(post, "", "    ") //美化输出
	if err != nil {
		fmt.Println("序列化失败", err)
		return
	}

	jsonStr := string(postJson)
	fmt.Printf(jsonStr)
	fmt.Println()

	// 查找 "Author" 的位置
	authorIndex := strings.Index(jsonStr, "\"Author\":")
	if authorIndex == -1 {
		fmt.Println("没有找到 Author 字段")
		return
	}

	// 提取 "Author" 对象的 JSON 字符串
	start := strings.Index(jsonStr[authorIndex:], "{") + authorIndex
	end := strings.Index(jsonStr[start:], "}") + start + 1
	authorJson := jsonStr[start:end]

	//fmt.Println(authorJson) 最后只能这样删去Author得到json的正确序列(哭

	var newAuthor Author

	err = json.Unmarshal([]byte(authorJson), &newAuthor)
	if err != nil {
		fmt.Println("反序列化失败", err)
		return
	}

	fmt.Println("反序列化后的 Author 实例：")
	fmt.Printf("Name: %s\nBio: %s\n", newAuthor.Name, newAuthor.Bio)
}
