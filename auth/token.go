package auth

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func ExtractToken(ctx *fiber.Ctx, tokenLookups ...string) string {
	if len(tokenLookups) == 0 {
		// default, only header
		tokenLookups = []string{"header:Authorization"}
	}

	for _, tokenLookup := range tokenLookups {
		// Initialize
		extractors := make([]func(c *fiber.Ctx) string, 0)
		rootParts := strings.Split(tokenLookup, ",")

		for _, rootPart := range rootParts {
			parts := strings.Split(strings.TrimSpace(rootPart), ":")

			switch parts[0] {
			case "header":
				extractors = append(extractors, jwtFromHeader(parts[1], "Bearer"))
			case "query":
				extractors = append(extractors, jwtFromQuery(parts[1]))
			case "param":
				extractors = append(extractors, jwtFromParam(parts[1]))
			case "cookie":
				extractors = append(extractors, jwtFromCookie(parts[1]))
			}
		}

		for _, extractor := range extractors {
			token := extractor(ctx)
			if token != "" {
				return token
			}
		}
	}

	return ""
}

// jwtFromHeader returns a function that extracts token from the request header.
func jwtFromHeader(header string, authScheme string) func(c *fiber.Ctx) string {
	return func(c *fiber.Ctx) string {
		auth := c.Get(header)
		l := len(authScheme)
		if len(auth) > l+1 && strings.EqualFold(auth[:l], authScheme) {
			return auth[l+1:]
		}
		return ""
	}
}

// jwtFromQuery returns a function that extracts token from the query string.
func jwtFromQuery(param string) func(c *fiber.Ctx) string {
	return func(c *fiber.Ctx) string {
		return c.Query(param)
	}
}

// jwtFromParam returns a function that extracts token from the url param string.
func jwtFromParam(param string) func(c *fiber.Ctx) string {
	return func(c *fiber.Ctx) string {
		return c.Params(param)
	}
}

// jwtFromCookie returns a function that extracts token from the named cookie.
func jwtFromCookie(name string) func(c *fiber.Ctx) string {
	return func(c *fiber.Ctx) string {
		return c.Cookies(name)
	}
}
