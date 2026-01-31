package model

type Page struct {
	Books       []*Book `json:"books"`
	PageNo      int64   `json:"page_no"`
	PageSize    int64   `json:"page_size"`
	TotalPageNo int64   `json:"total_page_no"`
	TotalRecord int64   `json:"total_record"`
	MinPrice    string  `json:"min_price"`
	MaxPrice    string  `json:"max_price"`
	IsLogin     bool    `json:"is_login"`
	Username    string  `json:"username"`
}

// IsHasPrev 判断是否有上一页
func (p *Page) IsHasPrev() bool {
	return p.PageNo > 1
}

// IsHasNext 判断是否有下一页
func (p *Page) IsHasNext() bool {
	return p.PageNo < p.TotalPageNo
}

// GetPrevPageNo 获取上一页
func (p *Page) GetPrevPageNo() int64 {
	if p.IsHasPrev() {
		return p.PageNo - 1
	}
	return 1
}

// GetNextPageNo 获取下一页
func (p *Page) GetNextPageNo() int64 {
	if p.IsHasNext() {
		return p.PageNo + 1
	}
	return p.TotalPageNo
}
