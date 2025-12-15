package cms

import (
	"time"

	"github.com/dtekltd/common/database"
	"github.com/dtekltd/common/pkg/users"
	"github.com/dtekltd/common/types"
	"github.com/dtekltd/common/utils"
)

func MapFullPost(post *Post) types.Params {
	data := MapPost(post)
	mapPostAuthor(post.CreatedBy, &data)
	mapPostCategories(post.ID, &data)
	relatedPosts(post, data)
	return data
}

func MapPost(post *Post) types.Params {
	meta := post.Meta.Data()
	title := meta.PageMeta.Title
	if title == "" {
		title = post.Name
	}
	return types.Params{
		"Title":        title,
		"Type":         post.Type,
		"Path":         post.Path,
		"Name":         post.Name,
		"Content":      post.Content,
		"FeatureImage": meta.PageMeta.FeatureImage,
		"IntroText":    meta.PageMeta.IntroText,
		"Description":  meta.PageMeta.Description,
		"CreatedAt":    utils.FormatFullVnDate(time.Unix(int64(post.CreatedAt), 0)),
		// "Meta":      meta.PageMeta,
		// "CreatedAt":    time.Unix(int64(post.CreatedAt), 0).Format("02/01/2006"),
		// "UpdatedAt":   time.Unix(int64(post.UpdatedAt), 0).Format("02/01/2006"),
		// "PublishedAt": time.Unix(int64(post.PublishedAt), 0).Format("02/01/2006"),
	}
}

func mapPostAuthor(accID uint64, data *types.Params) {
	author := &users.Account{}
	database.DB.Model(author).Select("id", "name").Find(author, accID)
	(*data)["Author"] = author
}

func mapPostCategories(postID uint64, data *types.Params) {
	// post's categories
	postCategories := []Category{}
	database.DB.Model(&Category{}).Where("id IN (?)", database.DB.Model(&PostCategory{}).
		Select("category_id").Where("post_id=?", postID)).
		Find(&postCategories)

	cats := []Category{}
	tags := []Category{}
	if len(postCategories) > 0 {
		for _, cat := range postCategories {
			if cat.Type == "category" {
				cats = append(cats, cat)
			} else {
				tags = append(tags, cat)
			}
		}
	}
	(*data)["Categories"] = cats
	(*data)["Tags"] = tags
}

func relatedPosts(post *Post, data types.Params) {
	relatedPosts := []Post{}
	database.DB.Model(&Post{}).
		Select("id", "name", "path", "meta", "created_at").
		Where("(id <> ?) AND (published_at>0) AND (type=?)", post.ID, post.Type).
		Limit(6).
		Order("created_at DESC").Find(&relatedPosts)

	if len(relatedPosts) > 0 {
		posts := []types.Params{}
		for _, post := range relatedPosts {
			posts = append(posts, MapPost(&post))
		}
		data["RelatedPosts"] = posts
	}
}
