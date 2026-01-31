package dao

import (
	"bookstore/config"
	"bookstore/model"
)

func GetAllCategories() ([]*model.Category, error) {
	sqlStr := "SELECT id, name, parent_id FROM categories"
	rows, err := config.DB.Query(sqlStr)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var categories []*model.Category
	for rows.Next() {
		var c model.Category
		rows.Scan(&c.ID, &c.Name, &c.ParentID)
		categories = append(categories, &c)
	}
	return categories, nil
}

// BuildCategoryTree 构建树形结构
func BuildCategoryTree() ([]*model.Category, error) {
	categories, err := GetAllCategories()
	if err != nil {
		return nil, err
	}
	// 用 map 保存 id -> category 的引用，方便组装树
	categoryMap := make(map[int]*model.Category)

	for _, c := range categories {
		categoryMap[c.ID] = c
	}

	var roots []*model.Category
	for _, c := range categories {
		if c.ParentID == nil {
			// 没有父级 => 根节点
			roots = append(roots, c)
		} else {
			// 有父级 => 挂到父节点的 Children
			parent, exists := categoryMap[*c.ParentID]
			if exists {
				parent.Children = append(parent.Children, c)
			}
		}
	}
	return roots, nil
}
