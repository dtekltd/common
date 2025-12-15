package site

import "github.com/gofiber/fiber/v2"

func NewPage(ctx *fiber.Ctx, title string) *Page {
	breadcrumbs := Breadcrumbs{
		NavItem{
			Name: title,
		},
	}
	if ctx.Path() != "/" {
		breadcrumbs = append(Breadcrumbs{
			NavItem{
				Name: "Home",
				URL:  "/",
			},
		}, breadcrumbs...)
	}
	return &Page{
		Title:       title,
		Breadcrumbs: &breadcrumbs,
	}
}
