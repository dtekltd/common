package utils

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	idLength  = 8
	alphabets = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func NewStringPtr(str string) *string {
	return &str
}

func GenerateOTP(length int) (string, error) {
	id := make([]byte, length)
	seed := "0123456789"
	for i := range length {
		char, err := rand.Int(rand.Reader, big.NewInt(int64(len(seed))))
		if err != nil {
			return "", err
		}
		id[i] = seed[char.Int64()]
	}
	return string(id), nil
}

func GenerateRandomString(length int, seeds ...string) (string, error) {
	id := make([]byte, length)
	seed := alphabets
	if len(seeds) > 0 {
		seed = seeds[0]
	}
	for i := range length {
		char, err := rand.Int(rand.Reader, big.NewInt(int64(len(seed))))
		if err != nil {
			return "", err
		}
		id[i] = seed[char.Int64()]
	}
	return string(id), nil
}

func GenerateID() (string, error) {
	return GenerateRandomString(idLength)
}

func GeneratePassword(length int) (string, error) {
	return GenerateRandomString(length, alphabets+"!@#$%^&*()_+")
}

func UCWord(input string) string {
	titleCaser := cases.Title(language.Und)
	words := strings.Fields(input)
	for i, word := range words {
		words[i] = titleCaser.String(word)
	}
	return strings.Join(words, " ")
}

func RemoveSignChars(s string) string {
	rules := []struct {
		pattern     string
		replacement string
	}{
		{`[àáạảãâầấậẩẫăằắặẳẵ]`, "a"},
		{`[èéẹẻẽêềếệểễ]`, "e"},
		{`[ìíịỉĩ]`, "i"},
		{`[òóọỏõôồốộổỗơờớợởỡ]`, "o"},
		{`[ùúụủũưừứựửữ]`, "u"},
		{`[ỳýỵỷỹ]`, "y"},
		{`đ`, "d"},
		{`[ÀÁẠẢÃÂẦẤẬẨẪĂẰẮẶẲẴ]`, "A"},
		{`[ÈÉẸẺẼÊỀẾỆỂỄ]`, "E"},
		{`[ÌÍỊỈĨ]`, "I"},
		{`[ÒÓỌỎÕÔỒỐỘỔỖƠỜỚỢỞỠ]`, "O"},
		{`[ÙÚỤỦŨƯỪỨỰỬỮ]`, "U"},
		{`[ỲÝỴỶỸ]`, "Y"},
		{`Đ`, "D"},
	}

	for _, rule := range rules {
		s = regexp.MustCompile(rule.pattern).ReplaceAllString(s, rule.replacement)
	}

	return s
}

func RemoveSpecialChars(s, skip string) string {
	return ReplaceSpecialChars(s, "", skip)
}

func ReplaceSpecialChars(s, replace, skip string) string {
	specialChars := "!@#$%^&*()_\\-+{}|:\"'<>,.?/\\[\\];.,/*~`=+"
	if skip != "" {
		specialChars = regexp.MustCompile("["+skip+"]").ReplaceAllString(specialChars, "")
	}
	s = regexp.MustCompile("["+specialChars+"]+").ReplaceAllString(s, replace)
	if replace != "" {
		s = regexp.MustCompile("["+replace+"]+").ReplaceAllString(s, replace)
		s = strings.Trim(s, " "+replace)
	}
	return s
}

func StringWithDefault(val string, defaultVal string) string {
	if val != "" {
		return val
	}
	return defaultVal
}

func StringToInt(val string) int {
	i, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}
	return i
}

func ToString(data any) string {
	switch data.(type) {
	case int:
		return fmt.Sprintf("%d", data)
	case float32:
	case float64:
		return fmt.Sprintf("%f", data)
	case string:
		return data.(string)
	default:
	}
	text, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return ""
	}
	return string(text)
}

func HashKey(key string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(key)))
}
