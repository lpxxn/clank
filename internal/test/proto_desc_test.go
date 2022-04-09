package test

import (
	"testing"

	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/stretchr/testify/assert"
)

func TestDynamicProto(t *testing.T) {
	//fileDescriptors := []*desc.FileDescriptor{}
	parser := protoparse.Parser{
		ImportPaths: []string{"./"},
	}
	fileDescriptors, err := parser.ParseFiles("protos/common.proto",
		"protos/model/students.proto",
		"protos/api/student_api.proto")
	assert.Nil(t, err)
	t.Logf("fileDescriptors: %v", fileDescriptors)
}
