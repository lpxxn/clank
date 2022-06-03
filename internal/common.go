package internal

import (
	"bytes"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/Knetic/govaluate"
	"github.com/lpxxn/clank/internal/clanklog"
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

func GenerateDefaultStringTemplate(templateText string) (string, error) {
	rev, err := GenerateTemplate(templateText, nil, nil)
	return string(rev), err
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
		"RandInt64Range": func(min, max int64) int64 {
			return rand.Int63n(max-min) + min
		},
		// eg: [910709413461759461,5141766179235889031]
		"RandInt64Slice": func(size int) string {
			revSlice := make([]string, 0, size)
			for i := 0; i < size; i++ {
				revSlice = append(revSlice, fmt.Sprintf("%d", rand.Int63()))
			}
			return fmt.Sprintf("[%s]", strings.Join(revSlice, ","))
		},
		"RandInt32": func() int32 {
			return rand.Int31()
		},
		"RandInt32Range": func(min, max int32) int32 {
			return rand.Int31n(max-min) + min
		},
		"RandInt32Slice": func(size int) string {
			revSlice := make([]string, 0, size)
			for i := 0; i < size; i++ {
				revSlice = append(revSlice, fmt.Sprintf("%d", rand.Int31()))
			}
			return fmt.Sprintf("[%s]", strings.Join(revSlice, ","))
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

func ParametersFromStr(str string, regExpr *regexp.Regexp) map[string]struct{} {
	parameters := make(map[string]struct{})
	match := regExpr.FindAllStringSubmatch(str, -1)
	idx := regExpr.SubexpIndex("parameter")
	for _, matchItem := range match {
		parameters[matchItem[idx]] = struct{}{}
	}
	return parameters
}
func ParamValue(param map[string]struct{}, jBody string) (map[string]interface{}, error) {
	paramValue := map[string]interface{}{}
	for key, _ := range param {
		v := jsonIterator.Get([]byte(jBody), keysInterfaceSlice(key)...)
		if v.LastError() != nil {
			return paramValue, v.LastError()
		}
		paramValue[key] = v.GetInterface()
	}
	return paramValue, nil
}

func ReplaceParamValue(param map[string]struct{}, jBody string, originStr string) (string, error) {
	paramValue, err := ParamValue(param, jBody)
	if err != nil {
		return "", err
	}
	for k, v := range paramValue {
		originStr = strings.ReplaceAll(originStr, "$"+k, fmt.Sprintf("%v", v))
	}
	return originStr, nil
}

func ValuableBoolExpression(expressionStr string) (bool, error) {
	expression, err := govaluate.NewEvaluableExpression(expressionStr)
	if err != nil {
		return false, err
	}
	result, err := expression.Evaluate(nil)
	clanklog.Info("evaluate result", result, err)
	if err != nil {
		return false, err
	}
	return result.(bool), nil
}
