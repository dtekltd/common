package cms

import (
	"strings"

	"github.com/dtekltd/common/pkg/site"
	"github.com/dtekltd/common/utils"
	"github.com/gofiber/fiber/v2"
)

func NewPage(ctx *fiber.Ctx, post *Post) *site.Page {
	pageMeta := &site.Meta{}
	if post.Meta != nil {
		meta := post.Meta.Data()
		if meta.PageMeta != nil {
			pageMeta = meta.PageMeta
		}
	}

	title := pageMeta.Title
	if title == "" {
		title = post.Name
	}

	breadcrumbs := site.Breadcrumbs{}
	if post.Path != "/" {
		// Home
		breadcrumbs = append(breadcrumbs, site.NavItem{
			Name: "Home",
			URL:  "/",
		})

		url := ""
		parts := strings.Split(post.Path, "/")
		if len(parts) > 1 {
			for _, part := range parts[:len(parts)-1] {
				url += "/" + part
				breadcrumbs = append(breadcrumbs, site.NavItem{
					Name: utils.UCWord(utils.ReplaceSpecialChars(part, " ", "")),
					URL:  url,
				})
			}
		}
	}
	breadcrumbs = append(breadcrumbs, site.NavItem{
		Name: title,
	})

	var content string
	if post.ContentType == "html" {
		// apply shortcode, if any
		content = parseShortcodes(post.Content)
	} else {
		content = post.Content
	}

	return &site.Page{
		Title:       title,
		Path:        post.Path,
		Content:     content,
		Breadcrumbs: &breadcrumbs,
		Meta:        pageMeta,
	}
}

func DefaultHomePage() *site.Page {
	return &site.Page{
		Title:   "Home Page",
		Content: "No Content",
		Breadcrumbs: &site.Breadcrumbs{
			site.NavItem{
				Name: "Home",
			},
		},
	}
}
