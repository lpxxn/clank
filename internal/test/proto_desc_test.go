package test

import (
	"testing"

	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/stretchr/testify/assert"
)

/*
https://github.com/jhump/protoreflect
*/

// parse proto
func TestDynamicProto(t *testing.T) {
	//fileDescriptors := []*desc.FileDescriptor{}
	parser := &protoparse.Parser{
		ImportPaths: []string{"./"},
	}
	//t.Log(desc.ResolveImport("protos/common.proto"))
	fileDescriptors, err := parser.ParseFiles(
		"protos/model/students.proto",
		"protos/api/student_api.proto",
		"protos/common.proto")
	assert.Nil(t, err)
	t.Logf("fileDescriptors: %v", fileDescriptors)

	for _, fileDesc := range fileDescriptors {
		t.Logf("===============\nfileDesc: %v", fileDesc)
		t.Logf("package: %s, gopackage: %s", fileDesc.GetPackage(), fileDesc.AsFileDescriptorProto().GetOptions().GetGoPackage())
		for _, msgDesc := range fileDesc.GetMessageTypes() {
			t.Logf("msgDesc: %v", msgDesc)
		}
		for _, servDesc := range fileDesc.GetServices() {
			t.Logf("service info: %v", servDesc)
			t.Logf("service name: %s", servDesc.GetName())
			for _, methodInfo := range servDesc.GetMethods() {
				t.Logf("methods %+v", methodInfo)
			}
		}
	}

}
