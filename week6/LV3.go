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

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
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
	} //确保数据库连接成功

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Println(err)
		}
	}(db) //确保函数结束时关闭数据库

	h := server.Default()

	h.POST("/add", func(c context.Context, ctx *app.RequestContext) {
		var student Student
		if err := ctx.BindAndValidate(&student); err != nil { //绑定并验证
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

	h.PUT("/profile", func(c context.Context, ctx *app.RequestContext) {
		var student Student
		if err := ctx.BindAndValidate(&student); err != nil {
			ctx.JSON(consts.StatusBadRequest, utils.H{"error": err.Error()})
			return
		}
		mu.Lock()
		defer mu.Unlock()
		if _, exists := students[student.Name]; exists { //检查请求的学生是否存在
			_, err := db.Exec("UPDATE students SET name = ?,id = ?,address = ?,birthday = ?,gender = ? WHERE name = ?")
			if err != nil {
				ctx.JSON(consts.StatusInternalServerError, utils.H{"error": err.Error()})
				return
			}
			ctx.JSON(consts.StatusOK, utils.H{"message": "成功更新学生信息"})
		} else {
			ctx.JSON(consts.StatusNotFound, utils.H{"error": "学生不存在"})
		}
	})

	h.PUT("/Update", func(c context.Context, ctx *app.RequestContext) {

		name := ctx.Query("name")   //按照名字获取学生参数
		field := ctx.Query("field") //获取想要修改的字段
		value := ctx.Query("value") //获取想要修改的值

		mu.Lock()
		defer mu.Unlock()

		_, err = db.Exec("UPDATE students SET ? = ? WHERE name = ?", field, value, name)
		if err != nil {
			ctx.JSON(consts.StatusInternalServerError, utils.H{"error": err.Error()})
			return
		}
		ctx.JSON(consts.StatusOK, utils.H{"message": "成功更新学生信息"})

	})

	h.GET("/search", func(c context.Context, ctx *app.RequestContext) {
		var student Student
		name := ctx.Query("name") //按照名字获取学生参数

		mu.Lock()
		defer mu.Unlock()

		row, err := db.Query("SELECT * FROM students WHERE name = ?", name)
		if err != nil {
			ctx.JSON(consts.StatusInternalServerError, utils.H{"error": err.Error()})
			return
		}

		for row.Next() {
			if err := row.Scan(&student.Name, &student.ID, &student.Address, &student.BirthDay, &student.Gender); err != nil {
				ctx.JSON(consts.StatusInternalServerError, utils.H{"error": err.Error()})
				return
			}
		}
		ctx.JSON(consts.StatusOK, utils.H{"student": student})
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

	h.POST("/register", func(c context.Context, ctx *app.RequestContext) {
		var user User
		if err := ctx.BindAndValidate(&user); err != nil {
			ctx.JSON(consts.StatusBadRequest, utils.H{"error": err.Error()})
			return
		}
		mu.Lock()
		defer mu.Unlock()
		_, err2 := db.Exec("INSERT INTO users (username,password) VALUES (?,?)", //连接users表并插入数据
			user.Username, user.Password)
		if err2 != nil {
			ctx.JSON(consts.StatusInternalServerError, utils.H{"error": err2.Error()})
			return
		}
		ctx.JSON(consts.StatusOK, utils.H{"message": "成功注册用户"})
	})

	h.POST("/login", func(c context.Context, ctx *app.RequestContext) {
		var user User
		var password string
		if err := ctx.BindAndValidate(&user); err != nil { //绑定user
			ctx.JSON(consts.StatusBadRequest, utils.H{"error": err.Error()})
			return
		}
		mu.Lock()
		defer mu.Unlock()
		err := db.QueryRow("select password from users where username = ? ", //查询users表中的密码
			user.Username).Scan(&password) //将查询结果赋值给password
		if err != nil {
			ctx.JSON(consts.StatusInternalServerError, utils.H{"error": err.Error()}) //查询失败
			return
		}
		if password != user.Password {
			ctx.JSON(consts.StatusUnauthorized, utils.H{"error": "账号或密码错误"}) //如果密码不匹配
			return
		}
		ctx.JSON(consts.StatusOK, utils.H{"message": "登录成功"})
	})

	h.Spin()
}
