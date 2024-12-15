package dao

import "Redrock/message-board/model"

func CreateMessage(message model.Message) error {
	createdAt := message.CreatedAt
	updatedAt := message.UpdatedAt

	_, err := db.Exec("INSERT INTO messages (user_id,content, created_at, updated_at) VALUES (?,?, ?, ?)", message.UserID, message.Content, createdAt, updatedAt)
	return err
}
