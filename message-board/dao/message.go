package dao

import (
	"Redrock/message-board/model"
	"database/sql"
)

func CreateMessage(message model.Message) error {
	createdAt := message.CreatedAt
	updatedAt := message.UpdatedAt

	_, err := db.Exec("INSERT INTO messages (user_id,content, created_at, updated_at) VALUES (?,?, ?, ?)", message.UserID, message.Content, createdAt, updatedAt)
	return err
}

func GetAllMessages() ([]model.Message, error) {
	rows, err := db.Query("SELECT * FROM messages")
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
		messages = append(messages, message)
	}
	return messages, err
}

func DeleteMessage(id int) error {
	_, err := db.Exec("DELETE FROM messages WHERE id = ?", id) //根据id删除消息
	return err
}
