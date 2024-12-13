package main

import (
	"context"
	"database/sql"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	_ "github.com/go-sql-driver/mysql"
	"log"
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

func main() {

	dsn := "042:123123@tcp(127.0.0.1:3306)/db042?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Println(err)
		return
	}

	if err = db.Ping(); err != nil {
		log.Println(err)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Println(err)
		}
	}(db)

	h := server.Default()

	h.POST("/add", func(c context.Context, ctx *app.RequestContext) {
		var student Student
		if err := ctx.BindAndValidate(&student); err != nil {
			ctx.JSON(consts.StatusBadRequest, utils.H{"error": err.Error()})
			return
		}
		mu.Lock()
		defer mu.Unlock()
		_, err2 := db.Exec("INSERT INTO students (name,id,address,birthday,gender) VALUES (?,?,?,?,?)",
			student.Name, student.ID, student.Address, student.BirthDay, student.Gender)
		if err2 != nil {
			ctx.JSON(consts.StatusInternalServerError, utils.H{"error": err2.Error()})
			return
		}
		students[student.Name] = student
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

	})

	h.DELETE("/Delete", func(c context.Context, ctx *app.RequestContext) {
		name := ctx.Query("name")
		mu.Lock()
		defer mu.Unlock()
		_, err := db.Exec("delete from students where name = ?", name)
		if err != nil {
			ctx.JSON(consts.StatusInternalServerError, utils.H{"error": err.Error()})
			return
		}
		ctx.JSON(consts.StatusOK, utils.H{"message": "成功删除学生"})
	})

	h.Spin()
}
