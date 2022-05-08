package internal

import (
	"bytes"
	"math/rand"
	"regexp"
	"strings"
	"text/template"
	"time"
)

const (
	LeftDelim  = "{{"
	RightDelim = "}}"
)

var templateRe = regexp.MustCompile(`.*{{.*}}.*`)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GenerateDefaultTemplate(templateText string) ([]byte, error) {
	return GenerateTemplate(templateText, nil, nil)
}

func GenerateTemplate(templateText string, templateData interface{}, params map[string]interface{}) ([]byte, error) {
	if !templateRe.MatchString(templateText) {
		return []byte(templateText), nil
	}
	t, err := template.New("tableTemplate").Delims(LeftDelim, RightDelim).Funcs(template.FuncMap{
		"Replace": func(old, new, src string) string {
			return strings.ReplaceAll(src, old, new)
		},
		"ToLower": func(s string) string {
			return strings.ToLower(s)
		},
		"ToUpper": func(s string) string {
			return strings.ToUpper(s)
		},
		"RandInt64": func() int64 {
			return rand.Int63()
		},
		"RandInt32": func() int32 {
			return rand.Int31()
		},
		"RandFixLenString": func(length int) string {
			b := make([]rune, length)
			for i := range b {
				b[i] = letters[rand.Intn(len(letters))]
			}
			return string(b)
		},
		"RandString": func(minLength, maxLength int) string {
			length := rand.Intn(maxLength-minLength) + minLength
			b := make([]rune, length)
			for i := range b {
				b[i] = letters[rand.Intn(len(letters))]
			}
			return string(b)
		},
		"Add": func(a, b int) int {
			return a + b
		},
		"Now": func() string {
			return time.Now().Format(time.RFC3339)
		},
		"NowDate": func() string {
			return time.Now().Format("2006-01-02")
		},
		"param": func(name string) interface{} {
			if v, ok := params[name]; ok {
				return v
			}
			return ""
		},
	}).Parse(templateText)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, templateData); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}