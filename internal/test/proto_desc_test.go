package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
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
		t.Logf("Metadata: %s", fileDesc.GetName()) // protos/api/student_api.proto
		for _, msgDesc := range fileDesc.GetMessageTypes() {
			t.Logf("msgDesc: %v", msgDesc)
			msgDesc.AsProto().ProtoMessage()
			proto.Marshal(msgDesc.AsDescriptorProto())
		}
		for _, servDesc := range fileDesc.GetServices() {
			t.Logf("service info: %v", servDesc)
			t.Logf("service name: %s", servDesc.GetName())
			for _, methodInfo := range servDesc.GetMethods() {
				t.Logf("methods %+v", methodInfo)
			}
		}
		CreateServiceDesc(fileDesc)
	}

}

func CreateServiceDesc(fileDesc *desc.FileDescriptor) {
	for _, servDescriptor := range fileDesc.GetServices() {
		serviceDesc := grpc.ServiceDesc{
			ServiceName: servDescriptor.GetName(),
			Metadata:    fileDesc.GetName(),
		}
		for _, methodDescriptor := range servDescriptor.GetMethods() {
			isServerStream := methodDescriptor.IsServerStreaming()
			isClientStream := methodDescriptor.IsClientStreaming()
			if isServerStream || isClientStream {
				streamDesc := grpc.StreamDesc{
					StreamName:    methodDescriptor.GetName(),
					Handler:       nil,
					ServerStreams: isServerStream,
					ClientStreams: isClientStream,
				}
				serviceDesc.Streams = append(serviceDesc.Streams, streamDesc)
			} else {
				methodDesc := grpc.MethodDesc{
					MethodName: methodDescriptor.GetName(),
					Handler:    createUnaryServerHandler(serviceDesc, methodDescriptor.GetName()),
				}
				serviceDesc.Methods = append(serviceDesc.Methods, methodDesc)
			}
		}

		//Methods: []grpc.MethodDesc{
		//	{
		//		MethodName: "NewStudent",
		//		Handler:    _StudentSrv_NewStudent_Handler,
		//	},
		//	{
		//		MethodName: "StudentByID",
		//		Handler:    _StudentSrv_StudentByID_Handler,
		//	},
		//},
		//Streams: []grpc.StreamDesc{
		//	{
		//		StreamName:    "AllStudent",
		//		Handler:       _StudentSrv_AllStudent_Handler,
		//		ServerStreams: true,
		//	},
		//	{
		//		StreamName:    "StudentInfo",
		//		Handler:       _StudentSrv_StudentInfo_Handler,
		//		ServerStreams: true,
		//		ClientStreams: true,
		//	},
		//},
	}
}

func createUnaryServerHandler(serviceDesc grpc.ServiceDesc, methodName string) func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	return func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
		fmt.Println(serviceDesc.ServiceName)
		fmt.Println(methodName)
		fmt.Println(srv)
		//in := new(QueryStudent)
		//if err := dec(in); err != nil {
		//	return nil, err
		//}
		//if interceptor == nil {
		//	return srv.(StudentSrvServer).StudentByID(ctx, in)
		//}
		//info := &grpc.UnaryServerInfo{
		//	Server:     srv,
		//	FullMethod: "/api.StudentSrv/StudentByID",
		//}
		//handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		//	return srv.(StudentSrvServer).StudentByID(ctx, req.(*QueryStudent))
		//}
		//return interceptor(ctx, in, info, handler)
		return nil, nil
	}

}
