package cmsApi

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterHandlers(router fiber.Router) {
	router = router.Group("/cms")

	// settings
	router.Get("/settings", getSettings)
	router.Post("/settings", updateSettings)

	// upload
	router.Post("/uploads", ckUpload)

	// categories
	router1 := router.Group("/categories")
	router1.Get("", withPostType(getCategories, "category"))
	router1.Get("/:id", withPostType(getCategory, "category"))
	router1.Put("/:id", withPostType(publishCategory, "category"))
	router1.Post("", withPostType(saveCategory, "category"))

	router1 = router.Group("/tags")
	router1.Get("", withPostType(getCategories, "tag"))
	router1.Get("/:id", withPostType(getCategory, "tag"))
	router1.Put("/:id", withPostType(publishCategory, "tag"))
	router1.Post("", withPostType(saveCategory, "tag"))

	// posts
	router2 := router.Group("/posts")
	router2.Get("", withPostType(getPosts, "post"))
	router2.Get("/:id", withPostType(getPost, "post"))
	router2.Put("/:id", withPostType(publishPost, "post"))
	router2.Post("", withPostType(savePost, "post"))

	// pages
	router2 = router.Group("/pages")
	router2.Get("", withPostType(getPosts, "page"))
	router2.Get("/:id", withPostType(getPost, "page"))
	router2.Put("/:id", withPostType(publishPost, "page"))
	router2.Post("", withPostType(savePost, "page"))

	// news
	router2 = router.Group("/news")
	router2.Get("", withPostType(getPosts, "news"))
	router2.Get("/:id", withPostType(getPost, "news"))
	router2.Put("/:id", withPostType(publishPost, "news"))
	router2.Post("", withPostType(savePost, "news"))

	// news
	router2 = router.Group("/products")
	router2.Get("", withPostType(getPosts, "product"))
	router2.Get("/:id", withPostType(getPost, "product"))
	router2.Put("/:id", withPostType(publishPost, "product"))
	router2.Post("", withPostType(savePost, "product"))

	// anounces
	router.Get("/announces", withPostType(getAnnounces, "news"))
}
