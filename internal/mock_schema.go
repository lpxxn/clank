package internal

import (
	"errors"
	"fmt"
	"regexp"

	jsonIter "github.com/json-iterator/go"
)

const requestToken = "$request"

var jsonIterator = jsonIter.ConfigCompatibleWithStandardLibrary

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
	if s.Kind == GRPC {
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
	}
	return nil
}

type ServerKind string

const (
	GRPC ServerKind = "grpc"
	HTTP ServerKind = "http"
)

type GrpcServerDescription struct {
	Name    string                    `yaml¡:"name" json:"name"`
	Methods GrpcMethodDescriptionList `yaml:"methods" json:"methods"`
}

func (g *GrpcServerDescription) GetMethod(methodName string) (*GrpcMethodDescription, error) {
	for _, method := range g.Methods {
		if method.Name == methodName {
			return method, nil
		}
	}
	return nil, fmt.Errorf("method %s not found", methodName)
}

type ServerDescriptionInterface interface {
	Validate() error
}

type ServerList []ServerDescriptionInterface

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

func (s ServerList) Validate() error {
	for _, item := range s {
		err := item.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

type GrpcServerDescriptionList []*GrpcServerDescription

func (s GrpcServerDescriptionList) Validate() error {
	for _, item := range s {
		if err := item.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (s GrpcServerDescriptionList) GetMethod(servName string, methodName string) (*GrpcMethodDescription, error) {
	for _, item := range s {
		if item.Name == servName {
			return item.GetMethod(methodName)
		}
	}
	return nil, fmt.Errorf("server: %s method %s not found", servName, methodName)
}

func (s GrpcServerDescriptionList) ToInterface() []ServerDescriptionInterface {
	result := make([]ServerDescriptionInterface, 0, len(s))
	for _, item := range s {
		result = append(result, item)
	}
	return result
}

func (s *GrpcServerDescription) Validate() error {
	if s.Name == "" {
		return fmt.Errorf("server name is empty")
	}
	for _, m := range s.Methods {
		if err := m.Validate(); err != nil {
			return err
		}
	}
	return nil
}

type GrpcMethodDescriptionList []*GrpcMethodDescription

func (m GrpcMethodDescriptionList) Validate() error {
	for _, method := range m {
		if err := method.Validate(); err != nil {
			return err
		}
	}
	return nil
}

type GrpcMethodDescription struct {
	Name            string                           `yaml:"name" json:"name"`
	DefaultResponse string                           `yaml:"defaultResponse" json:"defaultResponse"`
	Conditions      ResponseConditionDescriptionList `yaml:"conditions" json:"conditions"`
	Parameters      map[string]string                `yaml:"-" json:"-"`
}

var re = regexp.MustCompile(`\$request.(?P<parameter>[.\w]+)`)

func (m *GrpcMethodDescription) Validate() error {
	if m.Name == "" {
		return errors.New("method name is empty")
	}
	if m.DefaultResponse == "" {
		return errors.New("method default response is empty")
	}
	m.Parameters = map[string]string{}
	for _, c := range m.Conditions {
		if c.Condition == "" || c.Response == "" {
			return errors.New("condition or response is empty")
		}
		match := re.FindAllStringSubmatch(c.Condition, -1)
		idx := re.SubexpIndex("parameter")
		c.Parameters = map[string]struct{}{}
		for _, matchItem := range match {
			m.Parameters[requestToken+"."+matchItem[idx]] = matchItem[idx]
			c.Parameters[matchItem[idx]] = struct{}{}
			fmt.Println(matchItem[idx])
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
