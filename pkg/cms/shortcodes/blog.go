package shortcodes

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/dtekltd/common/database"
	"github.com/dtekltd/common/pkg/cms"
	"github.com/dtekltd/common/system"
	"github.com/dtekltd/common/utils"
	"gorm.io/gorm"
)

var blogTmpl = `
<div class="row gy-4">
	{{ range .Posts }}
	<div class="col-lg-4">
		<article id="post-{{ .ID }}">
			<div class="post-img">
				<img src="{{ .FeatureImage }}" alt="" class="img-fluid">
			</div>
			<p class="post-category">{{ .CreatedAt }}</p>
			<h2 class="title">
				<a href="{{ .Path }}">{{ .Name }}</a>
			</h2>
			{{ with .Categories }}
			<p class="d-none post-category">
				{{ range . }}{{.Name}} {{ end }}
			</p>
			{{ end }}
			<p>{{ .Description }}</p>
			<div>
				<a href="{{ .Path }}" class="btn btn-danger">Đọc tiếp <i class="bi bi-arrow-right"></i></a>
			</div>
		</article>
	</div>
	{{ end }}
</div>
`

func BlogHandler(tmpl string, args map[string]string) string {
	t, err := template.New("[blog]").Parse(tmpl)
	if err != nil {
		system.Logger.Errorf("cms::shortcode parse template error: %v", err)
		return fmt.Sprintf("{blog error=\"%v\"}", err)
	}

	postType, ok := args["type"]
	if !ok || postType == "" {
		postType = "post"
	}

	lmt, ok := args["limit"]
	if !ok {
		lmt = "3"
	}
	limit := utils.StringToInt(lmt)

	catIDs := []uint64{}
	if cat, ok := args["cat"]; ok {
		cats := strings.Split(cat, ",")
		database.DB.Model(&cms.Category{}).Select("id").
			Where("alias IN (?)", cats).
			Find(&catIDs)
	}

	posts := []cms.Post{}
	query := database.DB.Model(&cms.Post{}).
		Preload("Categories", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "alias")
		}).
		Select("id", "name", "path", "meta", "created_at").
		Where("(type=?) AND (published_at > 0)", postType).
		Order("created_at DESC")

	if len(catIDs) > 0 {
		query.Joins("INNER JOIN cms_post_categories c ON c.post_id=cms_posts.id AND c.category_id IN (?)", catIDs)
	}
	query.Limit(limit).Find(&posts)

	items := map[uint64]any{}
	for _, post := range posts {
		items[post.ID] = cms.MapPost(&post)
		// meta := post.Meta.Data().PageMeta
		// items[post.ID] = map[string]any{
		// 	"ID":           post.ID,
		// 	"Name":         post.Name,
		// 	"Path":         post.Path,
		// 	"Meta":         meta,
		// 	"FeatureImage": meta.FeatureImage,
		// 	"Description":  meta.Description,
		// 	"Categories":   post.Categories,
		// 	"CreatedAt":    utils.FormatFullVnDate(time.Unix(int64(post.CreatedAt), 0)),
		// }
	}

	var out bytes.Buffer
	if err = t.Execute(&out, map[string]any{
		"Posts": items,
	}); err != nil {
		system.Logger.Errorf("cms::shortcode execute template error: %v", err)
		return fmt.Sprintf("{blog error=\"%v\"}", err)
	}

	return out.String()
}
