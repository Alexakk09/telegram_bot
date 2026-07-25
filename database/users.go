package database

import "log"

func SaveUserCity(chatID int64, city string) {

	query := `
	INSERT OR REPLACE INTO users(chat_id, city)
	VALUES (?, ?)
	`

	_, err := DB.Exec(query, chatID, city)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("City saved successfully!")

}

func GetUserCity(chatID int64) (string, bool) {

	var city string

	query := `
	SELECT city
	FROM users
	WHERE chat_id = ?
	`

	err := DB.QueryRow(query, chatID).Scan(&city)

	if err != nil {
		return "", false
	}

	return city, true

}