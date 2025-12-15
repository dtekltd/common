package shortcodes

import "github.com/dtekltd/common/pkg/cms"

func Shortcodes() map[string]cms.Shortcode {
	return map[string]cms.Shortcode{
		"blog": {
			Name:     "blog",
			Template: blogTmpl,
			Handler:  BlogHandler,
		},
		"news": {
			Name:     "news",
			Template: newsTmpl,
			Handler:  NewsHandler,
		},
		"swiper": {
			Name:     "swiper",
			Template: swiperTmpl,
			Handler:  SwiperHandler,
		},
	}
}
