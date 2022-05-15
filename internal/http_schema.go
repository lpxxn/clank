package internal

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
)

type httpServerDescriptor struct {
	MethodDescriptor []*httpMethodDescriptor `yaml:"methods"`
	methodMap        map[string]*httpMethodDescriptor
}

func (h *httpServerDescriptor) Validate() error {
	if len(h.MethodDescriptor) == 0 {
		return errors.New("no methods defined")
	}
	if h.methodMap == nil {
		h.methodMap = make(map[string]*httpMethodDescriptor)
	}
	for _, m := range h.MethodDescriptor {
		if err := m.Validate(); err != nil {
			return err
		}
		h.methodMap[m.Path] = m
	}
	return nil
}

func (h *httpServerDescriptor) GetMethod(path string) *httpMethodDescriptor {
	return h.methodMap[path]
}

type httpMethodDescriptor struct {
	Name               string                           `yaml:"name,required"`
	Path               string                           `yaml:"path"`
	Method             string                           `yaml:"method"`
	DefaultResponse    string                           `yaml:"defaultResponse"`
	Conditions         ResponseConditionDescriptionList `yaml:"conditions" json:"conditions"`
	responseParameters map[string]struct{}              `yaml:"-" json:"-"`
}

var httpRegex = regexp.MustCompile(`\$(?P<parameter>(param|body|query|form)\.\w+[.\w]*)`)

func (d *httpMethodDescriptor) Validate() error {
	if d.Name == "" {
		return errors.New("name is required")
	}
	if d.Path == "" {
		return errors.New("path is required")
	}
	if d.Method == "" {
		return errors.New("method is required")
	}
	if _, ok := methodMap[d.Method]; !ok {
		return errors.New("method is invalid")
	}
	if d.DefaultResponse == "" {
		return errors.New("defaultResponse is required")
	}
	d.responseParameters = ParametersFromStr(d.DefaultResponse)
	for _, condition := range d.Conditions {
		if len(condition.Condition) == 0 {
			return fmt.Errorf("http method %s condition is empty", d.Name)
		}
		if len(condition.Response) == 0 {
			return fmt.Errorf("http method %s condition response is empty", d.Name)
		}
		condition.Parameters = ParametersFromStr(condition.Condition)
		condition.ResponseParameters = ParametersFromStr(condition.Response)
	}
	return nil
}

func ParametersFromStr(str string) map[string]struct{} {
	parameters := make(map[string]struct{})
	match := httpRegex.FindAllStringSubmatch(str, -1)
	idx := httpRegex.SubexpIndex("parameter")
	for _, matchItem := range match {
		parameters[matchItem[idx]] = struct{}{}
	}
	return parameters
}
func ParamValue(param map[string]struct{}, jBody string) (map[string]interface{}, error) {
	paramValue := map[string]interface{}{}
	for key, _ := range param {
		v := jsonIterator.Get([]byte(jBody), key)
		if v.LastError() != nil {
			return paramValue, v.LastError()
		}
		paramValue[key] = v.GetInterface()
	}
	return paramValue, nil
}

func (h *httpServerDescriptor) GetResponse(methodName string, jBody string) (string, error) {
	method := h.methodMap[methodName]
	if len(method.Conditions) == 0 {
		return h.getResponse(nil, method.DefaultResponse, jBody)
	}
	for _, condition := range method.Conditions {
		if condition.Condition == "" {
			continue
		}
		if condition.Condition == "" {
			return condition.Response, nil
		}
		conditionStr := condition.Condition
		paramValue, err := ParamValue(condition.Parameters, jBody)
		if err != nil {
			return "", err
		}
		if len(paramValue) != len(condition.Parameters) {
			return h.getResponse(nil, method.DefaultResponse, jBody)
		}
		for k, v := range paramValue {
			conditionStr = strings.ReplaceAll(conditionStr, "$"+k, fmt.Sprintf("%v", v))
		}
		log.Printf("condition: %s", conditionStr)
		return h.getResponse(condition, conditionStr, jBody)
	}
	return method.DefaultResponse, nil
}

func (h *httpServerDescriptor) getResponse(condition *ResponseConditionDescription, body string, jBody string) (string, error) {
	if condition == nil || len(condition.ResponseParameters) == 0 {
		return GenerateDefaultStringTemplate(body)
	}
	conditionStr := body
	paramValue, err := ParamValue(condition.ResponseParameters, jBody)
	if err != nil {
		return "", err
	}
	if len(paramValue) != len(condition.ResponseParameters) {
		return "", fmt.Errorf("http response %s parameters %+v is not foud", body, condition.ResponseParameters)
	}
	for k, v := range paramValue {
		conditionStr = strings.ReplaceAll(conditionStr, grpcRequestToken+"."+k, fmt.Sprintf("%v", v))
	}
	return GenerateDefaultStringTemplate(conditionStr)
}
