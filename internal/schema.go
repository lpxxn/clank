package internal

import (
	"context"
	"errors"
	"fmt"
	"strings"

	jsonIter "github.com/json-iterator/go"
	"github.com/lpxxn/clank/internal/clanklog"
)

var jsonIterator = jsonIter.ConfigCompatibleWithStandardLibrary

type ServerKind string

const (
	GRPC ServerKind = "grpc"
	HTTP ServerKind = "http"
)

type ServerDescriptionInterface interface {
	Validate() error
	StartServer(port int) error
}

type SchemaDescriptionBase struct {
	Kind ServerKind `yaml:"kind" json:"kind"`
	Port int        `yaml:"port" json:"port"`
}

type SchemaDescription struct {
	SchemaDescriptionBase
	Server ServerDescriptionInterface `yaml:"servers" json:"servers"`
}

func (s SchemaDescription) Validate() error {
	if s.Kind != GRPC && s.Kind != HTTP {
		return errors.New("kind must be GRPC OR HTTP")
	}
	err := s.Server.Validate()
	if err != nil {
		return err
	}

	return nil
}

func (s SchemaDescription) ValidateAndStartServer() error {
	if err := s.Validate(); err != nil {
		return err
	}
	return s.Server.StartServer(s.Port)
}

func (s *SchemaDescription) Unmarshal(d []byte) error {
	kind := ServerKind(jsonIter.Get(d, "kind").ToString())
	if kind == GRPC {
		param := struct {
			SchemaDescriptionBase
			grpcSchema `yaml:"servers" json:"servers"`
		}{}
		if err := jsonIter.Unmarshal(d, &param); err != nil {
			return err
		}
		s.SchemaDescriptionBase = param.SchemaDescriptionBase
		s.Server = &param.grpcSchema

	} else if kind == HTTP {

	}

	return nil
}
func (s *SchemaDescription) UnmarshalYAML(unmarshal func(interface{}) error) error {
	b := SchemaDescriptionBase{}
	if err := unmarshal(&b); err != nil {
		return err
	}
	kind := b.Kind
	if kind == GRPC {
		param := &grpcSchema{}
		if err := unmarshal(param); err != nil {
			return err
		}
		s.Server = param
	} else if kind == HTTP {
		param := &httpSchema{}
		if err := unmarshal(param); err != nil {
			return err
		}
		s.Server = param
	}
	s.SchemaDescriptionBase = b
	return nil
}

type ResponseConditionDescriptionList []*ResponseConditionDescription

type ResponseConditionDescription struct {
	Condition          string              `yaml:"condition" json:"condition"`
	Response           string              `yaml:"response" json:"response"`
	Parameters         map[string]struct{} `yaml:"-" json:"-"`
	ResponseParameters map[string]struct{} `yaml:"-" json:"-"`
}

type CallbackRequest interface {
	MakRequest() error
}

type CallbackDescriptionBase struct {
	Kind ServerKind `yaml:"kind" json:"kind"`
}

type CallbackDescription struct {
	CallbackDescriptionBase
	Request []CallbackRequest `yaml:"request" json:"request"`
}

type HttpCallbackDescription struct {
	Method           string            `yaml:"method"`
	URL              string            `yaml:"url" json:"url"`
	Header           map[string]string `yaml:"header" json:"header"`
	Body             string            `yaml:"body" json:"body"`
	DelayTime        int64             `yaml:"delayTime" json:"delayTime"`
	urlParameters    map[string]struct{}
	headerParameters map[string]struct{}
	bodyParameters   map[string]struct{}
}

func (s *CallbackDescription) UnmarshalYAML(unmarshal func(interface{}) error) error {
	b := CallbackDescriptionBase{}
	if err := unmarshal(&b); err != nil {
		return err
	}
	kind := b.Kind
	if kind == GRPC {

	} else if kind == HTTP {

	}
	return nil
}

type HttpCallbackDescriptionList []*HttpCallbackDescription

func (c HttpCallbackDescriptionList) Validate() error {
	for _, item := range c {
		if err := item.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (h *HttpCallbackDescription) Validate() error {
	if h.Method == "" {
		return errors.New("callback method is required")
	}
	if h.URL == "" {
		return errors.New("callback URL is required")
	}
	if h.Body == "" {
		return errors.New("callback method is required")
	}
	if _, ok := methodMap[h.Method]; !ok {
		return errors.New("callback method is invalid")
	}
	h.urlParameters = ParametersFromStr(h.URL, httpRegex)
	for _, v := range h.Header {
		for key, item := range ParametersFromStr(v, httpRegex) {
			h.headerParameters[key] = item
		}
	}
	h.bodyParameters = ParametersFromStr(h.Body, httpRegex)
	return nil
}

func (h *HttpCallbackDescription) makeRequest(ctx context.Context, jBody string) error {
	var err error
	url := h.URL
	url, err = ReplaceParamValue(h.urlParameters, jBody, url)
	if err != nil {
		clanklog.Errorf("callback get url param value error: %+v", err)
		return err
	}
	headerValue, err := ParamValue(h.headerParameters, jBody)
	if err != nil {
		clanklog.Errorf("callback get header param value error: %+v", err)
		return err
	}
	header := map[string]string{}
	for headerKey, headerV := range h.Header {
		for k, v := range headerValue {
			headerV = strings.ReplaceAll(headerV, "$"+k, fmt.Sprintf("%v", v))
		}
		header[headerKey] = headerV
	}
	body := h.Body
	body, err = ReplaceParamValue(h.bodyParameters, jBody, body)
	if err != nil {
		clanklog.Errorf("callback get body param value error: %s", err)
		return err
	}

	resp, err := NewHttpRequestWithHeader(ctx, h.Method, url, []byte(body), header)
	if err != nil {
		clanklog.Errorf("callback url: %s request error: %+v", url, err)
		return err
	}
	clanklog.Infof("callback url: %s response body: %s", url, string(resp))
	return nil
}
