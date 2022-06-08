package internal

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/lpxxn/clank/internal/clanklog"
)

type httpSchema struct {
	Server httpServerDescriptor `yaml:"server"`
}

func (h *httpSchema) Validate() error {
	if err := h.Server.Validate(); err != nil {
		return err
	}
	return nil
}

func (h *httpSchema) StartServer(port int) error {
	serv := NewHttpServer(&h.Server)
	if err := serv.MethodHandler(); err != nil {
		return err
	}
	return serv.StartWithPort(port)
}

type httpServerDescriptor struct {
	MethodDescriptor []*httpMethodDescriptor          `yaml:"methods"`
	methodMap        map[string]*httpMethodDescriptor `yaml:"-"`
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
	Name               string                           `yaml:"name"`
	Path               string                           `yaml:"path"`
	Method             string                           `yaml:"method"`
	DefaultResponse    string                           `yaml:"defaultResponse"`
	Conditions         ResponseConditionDescriptionList `yaml:"conditions" json:"conditions"`
	responseParameters map[string]struct{}              `yaml:"-" json:"-"`
	HttpCallback       HttpCallbackDescriptionList      `yaml:"httpCallback"`
}

var httpRegex = regexp.MustCompile(`\$(?P<parameter>(param|body|query|form|header)\.\w+[.\w]*)`)
var httpCallbackRegex = regexp.MustCompile(`\$(?P<parameter>(param|body|query|form|header|response|request)\.\w+[.\w]*)`)

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
	d.responseParameters = ParametersFromStr(d.DefaultResponse, httpRegex)
	for _, condition := range d.Conditions {
		if len(condition.Condition) == 0 {
			return fmt.Errorf("http method %s condition is empty", d.Name)
		}
		if len(condition.Response) == 0 {
			return fmt.Errorf("http method %s condition response is empty", d.Name)
		}
		condition.Parameters = ParametersFromStr(condition.Condition, httpRegex)
		condition.ResponseParameters = ParametersFromStr(condition.Response, httpRegex)
	}
	return d.HttpCallback.Validate()
}

func (h *httpServerDescriptor) GetResponse(methodName string, jBody string) (string, error) {
	method := h.methodMap[methodName]
	if len(method.Conditions) == 0 {
		return h.getResponse(method, nil, method.DefaultResponse, jBody)
	}
	for _, condition := range method.Conditions {
		if condition.Condition == "" {
			continue
		}
		conditionStr := condition.Condition
		paramValue, err := ParamValue(condition.Parameters, jBody)
		if err != nil {
			clanklog.Errorf("get condition param value error: %s", err)
			continue
		}
		if len(paramValue) != len(condition.Parameters) {
			return h.getResponse(method, condition, method.DefaultResponse, jBody)
		}
		for k, v := range paramValue {
			conditionStr = strings.ReplaceAll(conditionStr, "$"+k, fmt.Sprintf("%v", v))
		}
		clanklog.Infof("condition: %s", conditionStr)
		result, err := ValuableBoolExpression(conditionStr)
		if err != nil {
			return "", err
		}
		if result {
			return h.getResponse(method, condition, condition.Response, jBody)
		}
	}
	return h.getResponse(method, nil, method.DefaultResponse, jBody)
}

func (h *httpServerDescriptor) getResponse(method *httpMethodDescriptor, condition *ResponseConditionDescription, body string, jBody string) (string, error) {
	if condition == nil {
		return h.getResponseByParameters(body, jBody, method.responseParameters)
	}

	return h.getResponseByParameters(body, jBody, condition.ResponseParameters)
}

func (h *httpServerDescriptor) getResponseByParameters(body string, jBody string, parameters map[string]struct{}) (string, error) {
	if len(parameters) == 0 {
		return GenerateDefaultStringTemplate(body)
	}
	paramValue, err := ParamValue(parameters, jBody)
	if err != nil {
		return "", err
	}
	if len(paramValue) != len(parameters) {
		return "", fmt.Errorf("response parameters is not match, response: %s, param: %+v", body, parameters)
	}
	conditionStr := body
	for k, v := range paramValue {
		conditionStr = strings.ReplaceAll(conditionStr, "$"+k, fmt.Sprintf("%v", v))
	}
	return GenerateDefaultStringTemplate(conditionStr)
}

func (h *httpServerDescriptor) makeCallback(methodName string, jBody string) {
	method := h.methodMap[methodName]
	if len(method.HttpCallback) == 0 {
		return
	}
	for _, callback := range method.HttpCallback {
		delayTime := callback.DelayTime
		if delayTime <= 0 {
			delayTime = 1
		}
		time.AfterFunc(time.Duration(delayTime)*time.Second, func() {
			if err := callback.makeRequest(context.Background(), jBody); err != nil {
				clanklog.Errorf("callback err: %+v", err)
			}
		})
	}
}
