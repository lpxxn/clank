package internal

import (
	"errors"

	jsonIter "github.com/json-iterator/go"
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

//type ServerList []ServerDescriptionInterface
//
//func (s ServerList) Validate() error {
//	for _, item := range s {
//		err := item.Validate()
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}

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
	URL       string            `yaml:"url" json:"url"`
	Form      map[string]string `yaml:"form" json:"form"`
	Header    map[string]string `yaml:"header" json:"header"`
	Body      string            `yaml:"body" json:"body"`
	DelayTime int64             `yaml:"delayTime" json:"delayTime"`
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

type CallbackDescriptionList []*CallbackDescription
