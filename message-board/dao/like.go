package dao

import (
	"log"
	"time"
)

func AddLike(userID, messageID int) error {
	_, err := db.Exec("INSERT INTO likes (user_id, message_id,created_at,updated_at) VALUES (?,?, ?,?)", userID, messageID, time.Now(), time.Now())
	return err
}

func DeleteLike(userID, messageID int) error {
	_, err := db.Exec("DELETE FROM likes WHERE user_id = ? AND message_id = ?", userID, messageID)
	_, err2 := db.Exec("UPDATE likes SET updated_at = ? WHERE user_id = ? AND message_id = ?", time.Now(), userID, messageID)
	if err2 != nil {
		log.Println(err2)
	}

	return err
}
