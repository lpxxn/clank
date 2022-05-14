package internal

import "errors"

type httpServerDescriptor struct {
	MethodDescriptor []*httpMethodDescriptor `yaml:"methods"`
}

func (h *httpServerDescriptor) Validate() error {
	if len(h.MethodDescriptor) == 0 {
		return errors.New("no methods defined")
	}
	for _, m := range h.MethodDescriptor {
		if err := m.Validate(); err != nil {
			return err
		}
	}
	return nil
}

type httpMethodDescriptor struct {
	Name            string                           `yaml:"name,required"`
	Path            string                           `yaml:"path"`
	Method          string                           `yaml:"method"`
	DefaultResponse string                           `yaml:"defaultResponse"`
	Conditions      ResponseConditionDescriptionList `yaml:"conditions" json:"conditions"`
}

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
	if d.DefaultResponse == "" {
		return errors.New("defaultResponse is required")
	}
	for _, condition := range d.Conditions {
		_ = condition
	}
	return nil
}
