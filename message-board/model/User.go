package model

import "time"

type User struct {
	ID        int       `json:"id"`
	Nickname  string    `json:"nickname"` //用户名
	Username  string    `json:"username"` //账号
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"` //创建时间
	UpdatedAt time.Time `json:"updated_at"` //更新时间
}
