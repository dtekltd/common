package strcase

import "github.com/dtekltd/common/utils"

// KebabCase converts a string into kebab case.
func KebabCase(s string) string {
	s = utils.RemoveSpecialChars(s, "")
	s = utils.RemoveSignChars(s)
	return delimiterCase(s, '-', false)
}

// UpperKebabCase converts a string into kebab case with capital letters.
func UpperKebabCase(s string) string {
	s = utils.RemoveSpecialChars(s, "")
	s = utils.RemoveSignChars(s)
	return delimiterCase(s, '-', true)
}
