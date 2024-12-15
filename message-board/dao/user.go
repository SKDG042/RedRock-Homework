package dao

import (
	"Redrock/message-board/model"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var db *sql.DB

func init() {
	dsn := "042:123123@tcp(127.0.0.1:3306)/message_board?charset=utf8mb4&parseTime=True&loc=UTC"
	db, _ = sql.Open("mysql", dsn)
	if err := db.Ping(); err != nil {
		log.Println(err)
	}
}

func CreateUser(user model.User) error {
	createdAt := user.CreatedAt.Format("2006-01-02 15:04:05")
	updatedAt := user.UpdatedAt.Format("2006-01-02 15:04:05")

	_, err := db.Exec("INSERT INTO users (nickname, username , password,created_at,updated_at) VALUES (?,?,?,?,?)", user.Nickname, user.Username, user.Password, createdAt, updatedAt)
	return err
}

func GetUser(username string) (*model.User, error) {
	var user model.User
	err := db.QueryRow("SELECT username, password FROM users WHERE username = ?", username).Scan(&user.Username, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, err
}
