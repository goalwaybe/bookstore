package model

import "time"

// Review 评论模型
type Review struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	BookID     int       `json:"book_id"`
	Rating     int       `json:"rating"`
	Comment    string    `json:"comment"`
	CreateTime time.Time `json:"create_time"`
}
