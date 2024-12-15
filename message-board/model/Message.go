package model

import "time"

type Message struct {
	ID        int       `json:"id"`         //消息的ID
	UserID    int       `json:"user_id"`    //用户的ID
	Content   string    `json:"content"`    //消息内容
	CreatedAt time.Time `json:"created_at"` //创建时间
	UpdatedAt time.Time `json:"updated_at"` //更新时间
	IsDeleted bool      `json:"is_deleted"` //是否删除,0表示未删除，1表示已删除
	ParentID  int       `json:"parent_id"`  //父消息的ID
}
