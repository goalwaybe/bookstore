package model

// Category 分类模型
type Category struct {
	ID       int         `json:"id"`
	Name     string      `json:"name"`
	ParentID *int        `json:"parent_id"`
	Children []*Category `json:"children"`
}
