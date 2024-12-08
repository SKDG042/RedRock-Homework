package main

import (
	"context"
	"encoding/json"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"os"
	"sync"
)

type Student struct {
	Name     string `json:"name"`
	ID       int    `json:"id"`
	Address  string `json:"address"`
	BirthDay string `json:"birthday"`
	Gender   string `json:"gender"`
}

var students = make(map[string]Student)
var mu sync.Mutex
var File = "students.json"

func load() error {
	file, err := os.Open(File)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	decoder := json.NewDecoder(file)
	return decoder.Decode(&students)
}

func save() error {
	file, err := os.Create(File)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)
	encoder := json.NewEncoder(file)
	return encoder.Encode(students)
}

func main() {
	if err := load(); err != nil {
		panic(err)
	} //加载文件
	h := server.Default()

	h.POST("/add", func(c context.Context, ctx *app.RequestContext) {
		var student Student
		if err := ctx.BindAndValidate(&student); err != nil {
			ctx.JSON(consts.StatusBadRequest, utils.H{"error": err.Error()})
			return
		}
		mu.Lock()
		defer mu.Unlock()
		students[student.Name] = student
		err := save() //保存文件
		if err != nil {
			return
		}
		ctx.JSON(consts.StatusOK, utils.H{"message": "成功添加学生"})
	})

	h.POST("/profile", func(c context.Context, ctx *app.RequestContext) {
		var student Student
		if err := ctx.BindAndValidate(&student); err != nil {
			ctx.JSON(consts.StatusBadRequest, utils.H{"error": err.Error()})
			return
		}
		mu.Lock()
		defer mu.Unlock()
		if _, exists := students[student.Name]; exists {
			students[student.Name] = student
			ctx.JSON(consts.StatusOK, utils.H{"message": "成功更新学生信息"})
		} else {
			ctx.JSON(consts.StatusNotFound, utils.H{"error": "学生不存在"})
		}
		err := save()
		if err != nil {
			return
		}
	})

	h.GET("/search", func(c context.Context, ctx *app.RequestContext) {
		name := ctx.Query("name")
		mu.Lock()
		defer mu.Unlock()
		if student, exists := students[name]; exists {
			ctx.JSON(consts.StatusOK, student)
		} else {
			ctx.JSON(consts.StatusNotFound, utils.H{"error": "学生不存在"})
		}
		err := save()
		if err != nil {
			return
		}
	})

	h.Spin()
}
