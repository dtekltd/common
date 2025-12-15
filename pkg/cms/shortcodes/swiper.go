package shortcodes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"

	"github.com/dtekltd/common/database"
	"github.com/dtekltd/common/pkg/cms"
	"github.com/dtekltd/common/system"
)

var swiperTmpl = `
<div class="swiper init-swiper">
	<script type="application/json" class="swiper-config">
	{
		"loop": true,
		"speed": 600,
		"autoplay": {
		"delay": 5000
		},
		"slidesPerView": "auto",
		"centeredSlides": true,
		"pagination": {
		"el": ".swiper-pagination",
		"type": "bullets",
		"clickable": true
		},
		"breakpoints": {
		"320": {
			"slidesPerView": 1,
			"spaceBetween": 0
		},
		"768": {
			"slidesPerView": 3,
			"spaceBetween": 20
		},
		"1200": {
			"slidesPerView": 5,
			"spaceBetween": 20
		}
		}
	}
	</script>
	<div class="swiper-wrapper align-items-center">
		{{ range .items }}
	<div class="swiper-slide">
		<a href="{{ .url }}" class="glightbox" data-gallery="images-gallery">
			<img src="{{ .url }}" class="img-fluid" alt="{{ with .title }}{{ . }}{{ end }}">
		</a>
	</div>
	{{ end }}
	</div>
	<div class="swiper-pagination"></div>
</div>
`

func SwiperHandler(tmpl string, args map[string]string) string {
	t, err := template.New("[swiper]").Parse(tmpl)
	if err != nil {
		system.Logger.Errorf("cms::shortcode parse template error: %v", err)
		return fmt.Sprintf("{swiper error=\"%v\"}", err)
	}

	if args["id"] == "" {
		return ""
	}

	post := &cms.Post{}
	if err := database.DB.Take(post, args["id"]).Error; err != nil {
		return fmt.Sprintf("{swiper error=\"%v\"}", err)
	}

	data := map[string]any{}
	if err := json.Unmarshal([]byte(post.Content), &data); err != nil {
		return fmt.Sprintf("{swiper error=\"%v\"}", err)
	}

	var out bytes.Buffer
	if err = t.Execute(&out, data); err != nil {
		system.Logger.Errorf("cms::shortcode execute template error: %v", err)
		return fmt.Sprintf("{swiper error=\"%v\"}", err)
	}

	return out.String()
}
