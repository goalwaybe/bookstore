package model

// Favorite 结构体
type Favorite struct {
	ID         int    `json:"id"`
	UserID     int    `json:"user_id"`
	BookID     int    `json:"book_id"`
	CreateTime string `json:"create_time"`
}
