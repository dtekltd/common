package cms

import (
	"maps"
	"strings"

	"github.com/dtekltd/common/system"
)

type ShortcodeHandler func(tmpl string, args map[string]string) string

type Shortcode struct {
	Name     string           `json:"name"`
	Template string           `json:"template"`
	Handler  ShortcodeHandler `json:"-"`
}

var shortcodes = map[string]Shortcode{}

func RegisterShortCode(sc Shortcode) {
	shortcodes[sc.Name] = sc
}

func RegisterShortCodes(scs map[string]Shortcode) {
	maps.Copy(shortcodes, scs)
}

// Parses shortcodes within a string
func parseShortcodes(content string) string {
	for name, sc := range shortcodes {
		for {
			start := strings.Index(content, "["+name)
			if start == -1 {
				break
			}

			end := strings.Index(content[start:], "]")
			if end == -1 {
				break
			}

			fullTag := content[start : start+end+1]
			args := make(map[string]string)

			if strings.Contains(fullTag, " ") {
				parts := strings.Split(fullTag[len(name)+2:end], " ")
				for _, part := range parts {
					if strings.Contains(part, "=") {
						kv := strings.SplitN(part, "=", 2)
						if len(kv) == 2 {
							args[kv[0]] = strings.Trim(kv[1], "\"") // Remove quotes
						}
					}
				}
			}

			system.Logger.Info("shortcode:", name, fullTag, args)
			replacement := sc.Handler(sc.Template, args)
			content = strings.Replace(content, fullTag, replacement, 1)
		}
	}

	return content
}
