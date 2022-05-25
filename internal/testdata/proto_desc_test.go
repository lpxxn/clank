package testdata

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/lpxxn/clank/internal/clanklog"
	"github.com/lpxxn/clank/internal/testdata/protos"
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

func TestMain(m *testing.M) {
	clanklog.NewLogger()
	os.Exit(m.Run())
}

func TestProto(t *testing.T) {
	req := api.QueryStudent{Id: 1}
	b, _ := json.Marshal(req)
	t.Log(string(b))
	resp := api.QueryStudentResponse{StudentList: []*model.Student{
		{
			Name: "heihei",
			Age:  1,
		},
		{
			Name: "hahaha",
			Age:  9,
		},
	}}
	b, _ = json.Marshal(resp)
	t.Log(string(b))

	result := protos.Result{
		Code: "OK",
		Desc: "OK",
		Data: b,
	}

	b, _ = json.Marshal(result)
	t.Log(string(b))
}

// parse proto
func TestDynamicProto(t *testing.T) {
	//fileDescriptors := []*desc.FileDescriptor{}
	goPath, ok := os.LookupEnv("GOPATH")
	if !ok {
		t.Error("GOPATH not found")
		return
	}
	goPath += "/src"
	clanklog.Infof("GOPATH: %s", goPath)
	parser := &protoparse.Parser{
		ImportPaths: []string{"./", goPath}, //goPath,

	}
	//t.Log(desc.ResolveImport("protos/common.proto"))
	fileDescriptors, err := parser.ParseFiles(
		//"protos/model/students.proto",
		"protos/api/student_api.proto",
		//"github.com/lpxxn/clank/internal/testdata/protos/common.proto",
		//"protos/common.proto",
	)
	if err, ok := err.(protoparse.ErrorWithPos); ok {
		t.Log(err.GetPosition())
		t.Log(err.Error())
		t.Fatal(err.Error())
	}
	//registerFiles, err := protodesc.NewFiles(desc.ToFileDescriptorSet(fileDescriptors...))
	//if err != nil {
	//	t.Fatal(err)
	//}
	//registerFiles.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
	//	if ofd, _ := protoregistry.GlobalFiles.FindFileByPath(fd.Path()); ofd != nil {
	//		return true
	//	}
	//
	//	err = protoregistry.GlobalFiles.RegisterFile(fd)
	//	if err != nil {
	//		return false
	//	}
	//	return true
	//})
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
			ServiceName: servDescriptor.GetFullyQualifiedName(), // api.StudentSrv
			Metadata:    fileDesc.GetName(),
		}
		unaryMethodMap[serviceDesc.ServiceName] = make(map[string]grpc.MethodDesc)

		for _, methodDescriptor := range servDescriptor.GetMethods() {
			isServerStream := methodDescriptor.IsServerStreaming()
			isClientStream := methodDescriptor.IsClientStreaming()
			if isServerStream || isClientStream {
				streamDesc := grpc.StreamDesc{
					StreamName:    methodDescriptor.GetName(),
					Handler:       createStreamHandler(serviceDesc, methodDescriptor),
					ServerStreams: isServerStream,
					ClientStreams: isClientStream,
				}
				serviceDesc.Streams = append(serviceDesc.Streams, streamDesc)
			} else {
				unaryMethodMap[serviceDesc.ServiceName][methodDescriptor.GetName()] = grpc.MethodDesc{
					MethodName: methodDescriptor.GetName(),
					Handler:    nil,
				}
				clanklog.Infof("method: %s run handler", methodDescriptor.GetName())
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
		Id:   111,
		Name: "abc",
	})

	if err != nil {
		t.Fatal(err)
	}

	t.Log(result)
	// 有pb.go和没有还不一样。
	rev2, err := client.StudentByID(context.Background(), &api.QueryStudent{
		Id: 11,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(rev2)

	rev2, err = client.StudentByID(context.Background(), &api.QueryStudent{
		Id: 456,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(rev2)
	//out := map[string]interface{}{}
	//err = conn.Invoke(context.Background(), "/api.StudentSrv/NewStudent", map[string]string{"name": "test"}, out)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//t.Log(out)

	revStream1, err := client.AllStudent(context.Background(), &empty.Empty{})
	if err != nil {
		t.Fatal(err)
	}

	allStu, err := revStream1.Recv()
	if err == io.EOF {
		t.Log("EOF")
	}
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", allStu)
}

var unaryMethodMap map[string]map[string]grpc.MethodDesc = make(map[string]map[string]grpc.MethodDesc)

func createUnaryServerHandler(serviceDesc grpc.ServiceDesc, methodDesc *desc.MethodDescriptor) func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {

	return func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
		clanklog.Info("unary handler serviceName:", serviceDesc.ServiceName)
		clanklog.Info(methodDesc.GetName())
		clanklog.Info(srv)
		//inputParam := dynamic.NewMessage(methodDesc.GetInputType())
		msgFactory := dynamic.NewMessageFactoryWithDefaults()
		inputParam := msgFactory.NewMessage(methodDesc.GetInputType())
		if err := dec(inputParam); err != nil {
			return nil, err
		}
		clanklog.Info("input ", methodDesc.GetInputType())

		outPut := msgFactory.NewMessage(methodDesc.GetOutputType())
		dynamicOutput, err := dynamic.AsDynamicMessage(outPut)
		if err != nil {
			clanklog.Error(err)
			return nil, err
		}
		clanklog.Infof("output: %+v", dynamicOutput)
		if methodDesc.GetName() == "NewStudent" {
			if err := dynamicOutput.UnmarshalJSON([]byte(`{"code": "OK", "desc": "abcdef"}`)); err != nil {
				return nil, err
			}
		} else if methodDesc.GetName() == "StudentByID" {
			if err := dynamicOutput.UnmarshalJSON([]byte(`{"studentList": [{"id":111,"name":"abc","age":1298498081},{"id":222,"name":"def","age":2019727887}]}`)); err != nil {
				return nil, err
			}
		}

		outPutJson, err := dynamicOutput.MarshalJSON()
		clanklog.Info(outPutJson)

		outPutJson, err = json.Marshal(outPut)
		if err != nil {
			return nil, err
		}
		clanklog.Info(outPutJson)
		//dynamicOutput.SetFieldByName("desc", "hahahahah")
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
	}

}

func createStreamHandler(serviceDesc grpc.ServiceDesc, methodDesc *desc.MethodDescriptor) func(srv interface{}, stream grpc.ServerStream) error {
	return func(srv interface{}, stream grpc.ServerStream) error {
		isServerStream := methodDesc.IsServerStreaming()
		isClientStream := methodDesc.IsClientStreaming()
		clanklog.Info(isServerStream)
		clanklog.Info(isClientStream)
		msgFactory := dynamic.NewMessageFactoryWithDefaults()
		inputType := methodDesc.GetInputType()
		inputParam := msgFactory.NewMessage(inputType)
		if err := stream.RecvMsg(inputParam); err != nil {
			return err
		}
		clanklog.Info(inputParam.String())

		outPut := msgFactory.NewMessage(methodDesc.GetOutputType())
		dynamicOutput, err := dynamic.AsDynamicMessage(outPut)
		if err != nil {
			return err
		}
		if err := dynamicOutput.UnmarshalJSON([]byte(`{"studentList": [{"id":111,"name":"abc","age":1298498081},{"id":222,"name":"def","age":2019727887}]}`)); err != nil {
			return err
		}
		return nil
	}
}

func TestProtoEntity(t *testing.T) {
	StudentList := []*model.Student{
		{
			Id:   111,
			Name: "abc",
			Age:  rand.Int31(),
		},
		{
			Id:   222,
			Name: "def",
			Age:  rand.Int31(),
		},
	}
	b, _ := json.Marshal(StudentList)
	t.Log(string(b))
}
