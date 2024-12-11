// 定义所属的包为 main，这是 Go 程序的入口包
package main

// 导入需要的标准库和第三方库
import (
	// context 包提供了上下文的定义，包括取消信号、截止时间和请求范围的数据
	"context"
	// database/sql 包提供了通用的数据库接口
	"database/sql"

	// 导入 Hertz 框架的应用程序包，用于处理 HTTP 请求和响应
	"github.com/cloudwego/hertz/pkg/app"
	// 导入 Hertz 框架的服务器包，用于创建服务器实例
	"github.com/cloudwego/hertz/pkg/app/server"
	// 导入 Hertz 框架的通用工具包，提供了一些实用函数
	"github.com/cloudwego/hertz/pkg/common/utils"
	// 导入 Hertz 框架的协议常量包，提供了一些常用的 HTTP 常量，如状态码
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	// 导入 MySQL 驱动包，用于连接和操作 MySQL 数据库
	_ "github.com/go-sql-driver/mysql"
	// log 包提供了简单的日志记录功能
	"log"
	// sync 包提供了基本的同步原语，如互斥锁
	"sync"
)

// 定义 Student 结构体，表示学生的信息
type Student struct {
	// 学生的姓名，JSON 字段名为 "name"
	Name string `json:"name"`
	// 学生的学号，JSON 字段名为 "id"
	ID int `json:"id"`
	// 学生的地址，JSON 字段名为 "address"
	Address string `json:"address"`
	// 学生的生日，JSON 字段名为 "birthday"
	BirthDay string `json:"birthday"`
	// 学生的性别，JSON 字段名为 "gender"
	Gender string `json:"gender"`
}

// 定义 User 结构体，表示用户的账号信息
type User struct {
	// 用户名，JSON 字段名为 "username"
	Username string `json:"username"`
	// 用户密码，JSON 字段名为 "password"
	Password string `json:"password"`
}

// 声明一个全局变量 students，类型为 map[string]Student，用于存储学生信息，键为学生姓名，值为学生结构体
var students = make(map[string]Student)

// 声明一个全局的互斥锁 mu，用于保证对共享资源的并发访问是安全的
var mu sync.Mutex

// 主函数，程序的入口点
func main() {

	// 定义数据库连接的 DSN（数据源名称），包含用户名、密码、主机地址、数据库名、字符集等信息
	dsn := "042:123123@tcp(127.0.0.1:3306)/db042?charset=utf8mb4&parseTime=True&loc=Local"

	// 使用 sql.Open 函数打开数据库连接，指定驱动名为 "mysql" 和数据源名称 dsn
	db, err := sql.Open("mysql", dsn)
	// 如果在打开数据库连接时发生错误，记录错误日志并退出程序
	if err != nil {
		log.Println(err)
		return
	}

	// 测试与数据库的连接是否正常
	if err = db.Ping(); err != nil {
		// 如果连接测试失败，记录错误日志
		log.Println(err)
	}

	// 在程序结束前关闭数据库连接，使用 defer 延迟调用
	defer func(db *sql.DB) {
		// 尝试关闭数据库连接
		err := db.Close()
		// 如果关闭时发生错误，记录错误日志
		if err != nil {
			log.Println(err)
		}
	}(db)

	// 创建一个默认的 Hertz 服务器实例
	h := server.Default()

	// 为路径 "/add" 添加一个处理 POST 请求的路由处理器函数，用于添加学生信息
	h.POST("/add", func(c context.Context, ctx *app.RequestContext) {
		// 定义一个 Student 类型的变量 student，用于存储解析后的请求数据
		var student Student
		// 从请求体中绑定并验证数据，将结果存储到 student 变量中
		if err := ctx.BindAndValidate(&student); err != nil {
			// 如果绑定或验证失败，返回 400 状态码和错误信息
			ctx.JSON(consts.StatusBadRequest, utils.H{"error": err.Error()})
			return
		}
		// 加锁，保证对共享资源的并发安全
		mu.Lock()
		// 在函数执行结束时解锁
		defer mu.Unlock()
		// 执行数据库插入操作，将新的学生信息插入到 students 表中
		_, err2 := db.Exec("INSERT INTO students (name,id,address,birthday,gender) VALUES (?,?,?,?,?)",
			student.Name, student.ID, student.Address, student.BirthDay, student.Gender)
		// 如果执行数据库操作时发生错误
		if err2 != nil {
			// 返回 500 状态码和错误信息，表示服务器内部错误
			ctx.JSON(consts.StatusInternalServerError, utils.H{"error": err2.Error()})
			return
		}
		// 将学生信息添加到内存中的 students 映射中，键为学生姓名，值为学生信息
		students[student.Name] = student
		// 返回 200 状态码和成功消息，表示成功添加学生
		ctx.JSON(consts.StatusOK, utils.H{"message": "成功添加学生"})
	})

	// 为路径 "/profile" 添加一个处理 POST 请求的路由处理器函数，用于更新学生信息
	h.POST("/profile", func(c context.Context, ctx *app.RequestContext) {
		// 定义一个 Student 类型的变量 student，用于存储解析后的请求数据
		var student Student
		// 从请求体中绑定并验证数据，将结果存储到 student 变量中
		if err := ctx.BindAndValidate(&student); err != nil {
			// 如果绑定或验证失败，返回 400 状态码和错误信息
			ctx.JSON(consts.StatusBadRequest, utils.H{"error": err.Error()})
			return
		}
		// 加锁，保证对共享资源的并发安全
		mu.Lock()
		// 在函数执行结束时解锁
		defer mu.Unlock()
		// 检查学生是否存在于内存中的 students 映射中
		if _, exists := students[student.Name]; exists {
			// 如果学生存在，更新其信息
			students[student.Name] = student
			// 返回 200 状态码和成功消息，表示成功更新学生信息
			ctx.JSON(consts.StatusOK, utils.H{"message": "成功更新学生信息"})
		} else {
			// 如果学生不存在，返回 404 状态码和错误信息
			ctx.JSON(consts.StatusNotFound, utils.H{"error": "学生不存在"})
		}
	})

	// 为路径 "/search" 添加一个处理 GET 请求的路由处理器函数，用于搜索学生信息
	h.GET("/search", func(c context.Context, ctx *app.RequestContext) {
		// 从查询参数中获取学生的姓名，例如 /search?name=张三
		name := ctx.Query("name")
		// 加锁，保证对共享资源的并发安全
		mu.Lock()
		// 在函数执行结束时解锁
		defer mu.Unlock()
		// 从内存中的 students 映射中查找是否存在指定姓名的学生
		if student, exists := students[name]; exists {
			// 如果找到学生，返回 200 状态码，并将学生信息序列化为 JSON 格式返回
			ctx.JSON(consts.StatusOK, student)
		} else {
			// 如果未找到学生，返回 404 状态码和错误信息，表示学生不存在
			ctx.JSON(consts.StatusNotFound, utils.H{"error": "学生不存在"})
		}
	})

	// 为路径 "/Delete" 添加一个处理 DELETE 请求的路由处理器函数，用于删除学生信息
	h.DELETE("/Delete", func(c context.Context, ctx *app.RequestContext) {
		// 从查询参数中获取要删除的学生姓名，例如 /Delete?name=张三
		name := ctx.Query("name")
		// 加锁，保证对共享资源的并发安全
		mu.Lock()
		// 在函数执行结束时解锁
		defer mu.Unlock()
		// 执行数据库删除操作，从 students 表中删除指定姓名的学生记录
		_, err := db.Exec("delete from students where name = ?", name)
		// 如果执行数据库操作时发生错误
		if err != nil {
			// 返回 500 状态码和错误信息，表示服务器内部错误
			ctx.JSON(consts.StatusInternalServerError, utils.H{"error": err.Error()})
			return
		}
		// 从内存中的 students 映射中删除指定的学生
		delete(students, name)
		// 返回 200 状态码和成功消息，表示成功删除学生
		ctx.JSON(consts.StatusOK, utils.H{"message": "成功删除学生"})
	})

	// 为路径 "/register" 添加一个处理 POST 请求的路由处理器函数，用于用户注册
	h.POST("/register", func(c context.Context, ctx *app.RequestContext) {
		// 定义一个 User 类型的变量 user，用于存储解析后的请求数据
		var user User
		// 从请求体中绑定并验证数据，将结果存储到 user 变量中
		if err := ctx.BindAndValidate(&user); err != nil {
			// 如果绑定或验证失败，返回 400 状态码和错误信息
			ctx.JSON(consts.StatusBadRequest, utils.H{"error": err.Error()})
			return
		}
		// 加锁，保证对数据库操作的并发安全
		mu.Lock()
		// 在函数执行结束时解锁
		defer mu.Unlock()
		// 执行数据库插入操作，将新的用户信息插入到 users 表中
		_, err2 := db.Exec("INSERT INTO users (username,password) VALUES (?,?)",
			user.Username, user.Password)
		// 如果执行数据库操作时发生错误
		if err2 != nil {
			// 返回 500 状态码和错误信息，表示服务器内部错误
			ctx.JSON(consts.StatusInternalServerError, utils.H{"error": err2.Error()})
			return
		}
		// 返回 200 状态码和成功消息，表示成功注册用户
		ctx.JSON(consts.StatusOK, utils.H{"message": "成功注册用户"})
	})

	// 为路径 "/login" 添加一个处理 POST 请求的路由处理器函数，用于用户登录
	h.POST("/login", func(c context.Context, ctx *app.RequestContext) {
		// 定义一个 User 类型的变量 user，用于存储解析后的请求数据
		var user User
		// 定义一个字符串变量 password，用于存储从数据库中查询到的密码
		var password string
		// 从请求体中绑定并验证数据，将结果存储到 user 变量中
		if err := ctx.BindAndValidate(&user); err != nil {
			// 如果绑定或验证失败，返回 400 状态码和错误信息
			ctx.JSON(consts.StatusBadRequest, utils.H{"error": err.Error()})
			return
		}
		// 加锁，保证对数据库操作的并发安全
		mu.Lock()
		// 在函数执行结束时解锁
		defer mu.Unlock()
		// 执行数据库查询操作，获取对应用户名的密码
		err := db.QueryRow("select password from users where username = ? ",
			user.Username).Scan(&password)
		// 如果查询过程中发生错误
		if err != nil {
			// 返回 404 状态码和错误信息，表示账号不存在或密码错误
			ctx.JSON(consts.StatusNotFound, utils.H{"error": "账号或密码错误"})
			return
		}
		// 比较请求中的密码和数据库中存储的密码是否一致
		if password != user.Password {
			// 如果密码不匹配，返回 401 状态码和错误信息，表示未授权
			ctx.JSON(consts.StatusUnauthorized, utils.H{"error": "账号或密码错误"})
			return
		}
		// 如果登录成功，返回 200 状态码和成功消息，表示登录成功
		ctx.JSON(consts.StatusOK, utils.H{"message": "登录成功"})
	})

	// 启动 HTTP 服务器，开始监听并处理请求
	h.Spin()
}
