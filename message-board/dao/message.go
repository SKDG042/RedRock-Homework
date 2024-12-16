package dao

import (
	"Redrock/message-board/model"
	"database/sql"
)

func CreateMessage(message model.Message) error {
	createdAt := message.CreatedAt
	updatedAt := message.UpdatedAt

	_, err := db.Exec("INSERT INTO messages (user_id,content,parent_id ,created_at, updated_at) VALUES (?,?, ?,?, ?)", message.UserID, message.Content, message.ParentID, createdAt, updatedAt)
	return err
}

func GetAllMessages() ([]model.Message, error) {
	rows, err := db.Query("SELECT id, user_id, content, created_at, updated_at, is_deleted, parent_id FROM messages WHERE parent_id IS NULL")
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var messages []model.Message
	for rows.Next() {
		var message model.Message
		if err := rows.Scan(&message.ID, &message.UserID, &message.Content, &message.CreatedAt, &message.UpdatedAt, &message.IsDeleted, &message.ParentID); err != nil {
			return nil, err
		}
		message.Children, err = GetChildMessages(message.ID) //获取父消息的子消息切片
		messages = append(messages, message)                 //将获取的父消息添加到messages切片中
	}
	return messages, err
}

func GetChildMessages(parentID int) ([]model.Message, error) {
	rows, err := db.Query("SELECT id, user_id, content, created_at, updated_at, is_deleted, parent_id FROM messages WHERE parent_id = ?", parentID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var children []model.Message
	for rows.Next() {
		var message model.Message
		if err := rows.Scan(&message.ID, &message.UserID, &message.Content, &message.CreatedAt, &message.UpdatedAt, &message.IsDeleted, &message.ParentID); err != nil {
			return nil, err
		}
		subChildren, err := GetChildMessages(message.ID) //开始递归获取子消息，直到没有子消息
		if err != nil {
			return nil, err
		}
		message.Children = subChildren
		children = append(children, message) //将子消息添加到父消息的Children切片中
	}
	return children, err
}

func DeleteMessage(id int) error {
	_, err := db.Exec("DELETE FROM messages WHERE id = ?", id) //根据id删除消息
	return err
}
