package shortcodes

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/dtekltd/common/database"
	"github.com/dtekltd/common/pkg/cms"
	"github.com/dtekltd/common/system"
	"github.com/dtekltd/common/utils"
)

var newsTmpl = `
{{ range .Posts }}
<div class="card mb-3" style="max-width: 100%;">
	<div class="row g-0">
		<div class="col-5 col-md-3">
			<div class="rounded-start h-100" style="background-image: url('{{ .FeatureImage }}'); background-size: cover; background-position: center; background-repeat: no-repeat;">&nbsp;</div>
		</div>
		<div class="col-7 col-md-9">
			<div class="card-body">
				<h5 class="card-title">{{ .Name }}</h5>
				<p class="card-text mb-0">
					{{ .Description }}... 
					<a href="{{ .Path }}">Đọc tiếp <i class="bi bi-arrow-right"></i></a>
				</p>
				<p class="card-text"><small class="text-muted">{{ .CreatedAt }}</small></p>
			</div>
		</div>
	</div>
</div>
{{ end }}
`

func NewsHandler(tmpl string, args map[string]string) string {
	t, err := template.New("[news]").Parse(tmpl)
	if err != nil {
		system.Logger.Errorf("cms::shortcode parse template error: %v", err)
		return fmt.Sprintf("{news error=\"%v\"}", err)
	}

	postType, ok := args["type"]
	if !ok || postType == "" {
		postType = "news"
	}

	lmt, ok := args["limit"]
	if !ok {
		lmt = "3"
	}
	limit := utils.StringToInt(lmt)

	posts := []cms.Post{}
	query := database.DB.Model(&cms.Post{}).
		Select("id", "name", "path", "meta", "created_at").
		Where("(type=?) AND (published_at > 0)", postType).
		Order("created_at DESC")

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
		return fmt.Sprintf("{news error=\"%v\"}", err)
	}

	return out.String()
}
