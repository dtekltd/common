package cms

import (
	"strings"

	"github.com/dtekltd/common/database"
	"github.com/dtekltd/common/system"
	"github.com/gofiber/fiber/v2"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
)

// Handler defines a function to serve HTTP requests.
type PostProcessor = func(*fiber.Ctx, *Post) error

func PageHandler(fallback fiber.Handler, postProcessor PostProcessor) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		post := &Post{}
		path := strings.Trim(ctx.Path(), "/")

		if path == "" {
			// home page
			// page, post = NewHomePage(ctx)
			settings := Settings()
			if err := database.DB.Take(post, settings.HomePageID).Error; err != nil {
				system.Logger.Errorf("Home page post #%d does not exist", settings.HomePageID)
			}
		} else if err := database.DB.Take(post, "published_at>0 AND path=?", path).Error; err != nil {
			return fallback(ctx)
		}

		view := post.Type
		layout := ctx.Query("layout", "main")
		data := fiber.Map{}

		if post.ID != 0 {
			if postProcessor != nil {
				if err := postProcessor(ctx, post); err != nil {
					system.Logger.Errorf("Post processor error: %s", err.Error())
				}
			}

			// process post
			switch post.Type {
			case "page":
				processPageData(ctx, post, data)
			case "post":
				fallthrough
			case "news":
				processPostData(post, data)
			}

			data["Page"] = NewPage(ctx, post)
			meta := post.Meta.Data()

			if meta.View != "" {
				view = meta.View
			}
			if meta.Layout != "" {
				layout = meta.Layout
			}
		} else {
			data["Page"] = DefaultHomePage()
		}

		return ctx.Render("views/"+view, data, "layouts/"+layout)
	}
}

func processPageData(ctx *fiber.Ctx, post *Post, data fiber.Map) {
	if post.ContentType == "markdown" {
		// render markdown to HTML
		post.Content = markdownToHTML(post.Content)
	}
	// special blog pages
	if postType := Settings().CustomParams.Get("blogPages." + post.Path); postType != nil {
		data["Blog"] = GetBlogPosts(ctx, postType.(string))
	}
}

func processPostData(post *Post, data fiber.Map) {
	data["Categories"] = AllCategories()
	data["Post"] = MapFullPost(post)
}

func markdownToHTML(md string) string {
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)
	output := markdown.ToHTML([]byte(md), nil, renderer)
	return string(output)
}
