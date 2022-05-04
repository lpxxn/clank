package internal

type SchemaDescription struct {
	Kind    ServerKind          `yaml:"kind" json:"kind"`
	Port    string              `yaml:"port" json:"port"`
	Servers []ServerDescription `yaml:"servers" json:"servers"`
}

type ServerKind string

const (
	GRPC ServerKind = "grpc"
	HTTP ServerKind = "http"
)

type ServerDescription struct {
	Name    string               `yaml:"name" json:"name"`
	Methods []*MethodDescription `yaml:"methods" json:"methods"`
}

type MethodDescription struct {
	Name            string                          `yaml:"name" json:"name"`
	DefaultResponse string                          `yaml:"defaultResponse" json:"defaultResponse"`
	Conditions      []*ResponseConditionDescription `yaml:"conditions" json:"conditions"`
}

type ResponseConditionDescription struct {
	Condition string `yaml:"condition" json:"condition"`
	Response  string `yaml:"response" json:"response"`
}
