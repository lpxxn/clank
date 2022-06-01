package internal

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const grpcRequestToken = "$request"

type grpcSchema struct {
	// gRpc
	ImportPath   []string                  `yaml:"importPath" json:"importPath"`
	ProtoPath    []string                  `yaml:"protoPath" json:"protoPath"`
	ProtosetPath string                    `yaml:"protosetPath" json:"protosetPath"`
	Servers      GrpcServerDescriptionList `yaml:"servers" json:"servers"`
}

func (g *grpcSchema) Validate() error {
	if len(g.ProtoPath) == 0 && len(g.ProtosetPath) == 0 {
		return errors.New("grpc protoPath or protosetPath must be set")
	}
	if len(g.Servers) == 0 {
		return errors.New("grpc servers must be set")
	}

	return g.Servers.Validate()
}

func (g *grpcSchema) StartServer(port int) error {
	var serv *gRpcServer
	var err error
	if g.ProtosetPath != "" {
		serv, err = ParseServerMethodsFromProtoset(g.ProtosetPath)
	} else {
		serv, err = ParseServerMethodsFromProto(g.ImportPath, g.ProtoPath)
	}
	if err != nil {
		return err
	}
	if err := ValidateGrpcServiceInputAndOutput(g.Servers, serv); err != nil {
		return err
	}
	if err := SetOutputFunc(g.Servers, serv); err != nil {
		return err
	}
	if err := serv.StartWithPort(port); err != nil {
		return err
	}
	return nil
}

type GrpcServerDescription struct {
	Name    string                    `yaml:"name" json:"name"`
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

func (g *GrpcServerDescription) Validate() error {
	if g.Name == "" {
		return fmt.Errorf("server name is empty")
	}
	for _, m := range g.Methods {
		if err := m.Validate(); err != nil {
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
	DefaultMetaData map[string]string                `yaml:"defaultMetaData"`
	Conditions      ResponseConditionDescriptionList `yaml:"conditions" json:"conditions"`
	Parameters      map[string]string                `yaml:"-" json:"-"`
}

const grpcRequestParam = "request"

var grpcRegex = regexp.MustCompile(`\$(?P<parameter>(request|header)\.\w+[.\w-_]*)`)

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
		c.Parameters = ParametersFromStr(c.Condition, grpcRegex)
		for key, _ := range c.Parameters {
			if strings.HasPrefix(key, grpcRequestParam) {
				m.Parameters[key] = strings.Split(key, ".")[1]
			}
			m.Parameters[key] = key
		}
		c.ResponseParameters = ParametersFromStr(c.Response, grpcRegex)
	}
	return nil
}
