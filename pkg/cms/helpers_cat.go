package cms

import "github.com/dtekltd/common/database"

func AllCategories() []Category {
	categories := []Category{}
	query := database.DB.Model(&Category{}).
		Select("id", "type", "name", "alias").
		Where("published_at>0 AND type=?", "category")
	query.Find(&categories)
	return categories
}
