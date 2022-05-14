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
}

type SchemaDescriptionBase struct {
	Kind ServerKind `yaml:"kind" json:"kind"`
	Port int        `yaml:"port" json:"port"`

	// gRpc
	ImportPath   []string `yaml:"importPath" json:"importPath"`
	ProtoPath    []string `yaml:"protoPath" json:"protoPath"`
	ProtosetPath string   `yaml:"protosetPath" json:"protosetPath"`
}

type SchemaDescription struct {
	SchemaDescriptionBase
	Servers ServerList `yaml:"servers" json:"servers"`
}

func (s SchemaDescription) Validate() error {
	if s.Kind != GRPC && s.Kind != HTTP {
		return errors.New("kind must be GRPC OR HTTP")
	}
	err := s.Servers.Validate()
	if err != nil {
		return err
	}

	return nil
}

func (s SchemaDescription) ValidateAndStartServer() error {
	if err := s.Validate(); err != nil {
		return err
	}
	switch s.Kind {
	case GRPC:
		serv, err := ParseServerMethodsFromProto(s.ImportPath, s.ProtoPath)
		if err != nil {
			return err
		}
		if err := ValidateServiceInputAndOutput(s.Servers, serv); err != nil {
			return err
		}
		if err := SetOutputFunc(s.Servers, serv); err != nil {
			return err
		}
		if err := serv.StartWithPort(s.Port); err != nil {
			return err
		}
	case HTTP:

	}
	return nil
}

func (s *SchemaDescription) Unmarshal(d []byte) error {
	kind := ServerKind(jsonIter.Get(d, "kind").ToString())
	if kind == GRPC {
		param := struct {
			SchemaDescriptionBase
			Servers GrpcServerDescriptionList `yaml:"servers" json:"servers"`
		}{}
		if err := jsonIter.Unmarshal(d, &param); err != nil {
			return err
		}
		s.SchemaDescriptionBase = param.SchemaDescriptionBase
		s.Servers = make(ServerList, 0, len(param.Servers))
		for _, server := range param.Servers {
			s.Servers = append(s.Servers, server)
		}
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
		param := struct {
			Servers GrpcServerDescriptionList `yaml:"servers" json:"servers"`
		}{}
		if err := unmarshal(&param); err != nil {
			return err
		}
		s.Servers = make(ServerList, 0, len(param.Servers))
		for _, server := range param.Servers {
			s.Servers = append(s.Servers, server)
		}
	} else if kind == HTTP {

	}
	s.SchemaDescriptionBase = b
	return nil
}

type ServerList []ServerDescriptionInterface

func (s ServerList) Validate() error {
	for _, item := range s {
		err := item.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

type ResponseConditionDescriptionList []*ResponseConditionDescription

type ResponseConditionDescription struct {
	Condition  string              `yaml:"condition" json:"condition"`
	Response   string              `yaml:"response" json:"response"`
	Parameters map[string]struct{} `yaml:"-" json:"-"`
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
