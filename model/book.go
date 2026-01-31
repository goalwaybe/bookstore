package model

// Book 结构体
type Book struct {
	ID           int     `json:"id"`
	Title        string  `json:"title"`
	Author       string  `json:"author"`
	Price        float64 `json:"price"`
	Sales        int     `json:"sales"`
	Stock        int     `json:"stock"`
	ImgPath      string  `json:"img_path"`
	CategoryID   int     `json:"category_id"`
	CategoryName string  `json:"category_name"`
}
