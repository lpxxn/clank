package internal

import (
	"errors"
	"fmt"
	"regexp"
)

const grpcRequestToken = "$request"

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

func (s GrpcServerDescriptionList) ToInterface() []ServerDescriptionInterface {
	result := make([]ServerDescriptionInterface, 0, len(s))
	for _, item := range s {
		result = append(result, item)
	}
	return result
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

var grpcParamRegex = regexp.MustCompile(`\$request.(?P<parameter>[.\w]+)`)

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
		match := grpcParamRegex.FindAllStringSubmatch(c.Condition, -1)
		idx := grpcParamRegex.SubexpIndex("parameter")
		c.Parameters = map[string]struct{}{}
		for _, matchItem := range match {
			m.Parameters[grpcRequestToken+"."+matchItem[idx]] = matchItem[idx]
			c.Parameters[matchItem[idx]] = struct{}{}
			fmt.Println(matchItem[idx])
		}
	}
	return nil
}
