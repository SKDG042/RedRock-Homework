package dao

import "time"

func AddLike(userID, messageID int) error {
	_, err := db.Exec("INSERT INTO likes (user_id, message_id,created_at) VALUES (?, ?,?)", userID, messageID, time.Now())
	return err

}
