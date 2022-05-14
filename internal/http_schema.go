package internal

type httpServerDescriptor struct {
	Name             string                  `yaml:"name"`
	MethodDescriptor []*httpMethodDescriptor `yaml:"methods"`
}

type httpMethodDescriptor struct {
	Name   string `yaml:"name"`
	Path   string `yaml:"path"`
	Method string `yaml:"method"`
}
