package testdata

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"testing"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/lpxxn/clank/internal/testdata/protos/api"
	"github.com/lpxxn/clank/internal/testdata/protos/model"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/proto"
)

/*
https://github.com/jhump/protoreflect
*/

const testPort int = 54312

var (
	testAddress = fmt.Sprintf(":%d", testPort)
)

// parse proto
func TestDynamicProto(t *testing.T) {
	//fileDescriptors := []*desc.FileDescriptor{}
	goPath, ok := os.LookupEnv("GOPATH")
	if !ok {
		t.Error("GOPATH not found")
		return
	}
	goPath += "/src"
	parser := &protoparse.Parser{
		ImportPaths: []string{"./", goPath}, //goPath,

	}
	//t.Log(desc.ResolveImport("protos/common.proto"))
	fileDescriptors, err := parser.ParseFiles(
		//"protos/model/students.proto",
		"protos/api/student_api.proto",
		//"protos/common.proto"
	)
	if err, ok := err.(protoparse.ErrorWithPos); ok {
		t.Log(err.GetPosition())
		t.Log(err.Error())
	}
	t.Log(err)
	assert.Nil(t, err)
	t.Logf("fileDescriptors: %v", fileDescriptors)

	for _, fileDesc := range fileDescriptors {
		t.Logf("===============\nfileDesc: %v", fileDesc)
		t.Logf("package: %s, gopackage: %s", fileDesc.GetPackage(), fileDesc.AsFileDescriptorProto().GetOptions().GetGoPackage())
		t.Logf("Metadata: %s", fileDesc.GetName()) // protos/api/student_api.proto
		for _, msgDesc := range fileDesc.GetMessageTypes() {
			t.Logf("msgDesc: %v", msgDesc)
			msgDesc.AsProto().ProtoMessage()
			b, err := proto.Marshal(msgDesc.AsDescriptorProto())
			t.Log(err)
			t.Log(string(b))
		}
		for _, servDesc := range fileDesc.GetServices() {
			t.Logf("service info: %v", servDesc)
			t.Logf("service name: %s", servDesc.GetName())
			for _, methodInfo := range servDesc.GetMethods() {
				t.Logf("methods %+v", methodInfo)
				input := methodInfo.GetInputType()
				t.Log(input)
				output := methodInfo.GetOutputType()
				t.Log(output)
			}
		}
		CreateServiceDesc(fileDesc)
	}

}

func CreateServiceDesc(fileDesc *desc.FileDescriptor) {
	for _, servDescriptor := range fileDesc.GetServices() {
		serviceDesc := grpc.ServiceDesc{
			ServiceName: servDescriptor.GetFullyQualifiedName(),
			Metadata:    fileDesc.GetName(),
		}
		unaryMethodMap[serviceDesc.ServiceName] = make(map[string]grpc.MethodDesc)

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
				unaryMethodMap[serviceDesc.ServiceName][methodDescriptor.GetName()] = grpc.MethodDesc{
					MethodName: methodDescriptor.GetName(),
					Handler:    nil,
				}
				methodDesc := grpc.MethodDesc{
					MethodName: methodDescriptor.GetName(),
					Handler:    createUnaryServerHandler(serviceDesc, methodDescriptor),
				}
				serviceDesc.Methods = append(serviceDesc.Methods, methodDesc)
			}
		}
		grpcServ := grpc.NewServer()
		grpcServ.RegisterService(&serviceDesc, nil)

		listener, err := net.Listen("tcp", testAddress)
		if err != nil {
			panic(err)
		}
		reflection.Register(grpcServ)
		grpcServ.Serve(listener)
	}
}

func TestRpcClient(t *testing.T) {
	conn, err := grpc.Dial(testAddress, grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	client := api.NewStudentSrvClient(conn)
	result, err := client.NewStudent(context.Background(), &model.Student{
		Name: "test",
	})

	if err != nil {
		t.Fatal(err)
	}

	t.Log(result)
	out := map[string]interface{}{}
	err = conn.Invoke(context.Background(), "/api.StudentSrv/NewStudent", map[string]string{"name": "test"}, out)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(out)

}

var unaryMethodMap map[string]map[string]grpc.MethodDesc = make(map[string]map[string]grpc.MethodDesc)

func createUnaryServerHandler(serviceDesc grpc.ServiceDesc, methodDesc *desc.MethodDescriptor) func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {

	return func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
		fmt.Println(serviceDesc.ServiceName)
		fmt.Println(methodDesc.GetName())
		fmt.Println(srv)
		//inputParam := dynamic.NewMessage(methodDesc.GetInputType())
		msgFactory := dynamic.NewMessageFactoryWithDefaults()
		inputParam := msgFactory.NewMessage(methodDesc.GetInputType())
		if err := dec(inputParam); err != nil {
			return nil, err
		}

		outPut := msgFactory.NewMessage(methodDesc.GetOutputType())
		dynamicOutput, err := dynamic.AsDynamicMessage(outPut)
		if err != nil {
			return nil, err
		}
		if err := dynamicOutput.UnmarshalJSON([]byte(`{"code": "OK", "desc": "abcdef"}`)); err != nil {
			return nil, err
		}
		outPutJson, err := dynamicOutput.MarshalJSON()
		fmt.Println(outPutJson)

		outPutJson, err = json.Marshal(outPut)
		if err != nil {
			return nil, err
		}
		fmt.Println(outPutJson)
		dynamicOutput.SetFieldByName("desc", "hahahahah")
		return dynamicOutput, nil
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
