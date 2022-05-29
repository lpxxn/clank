package testdata

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"testing"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"github.com/lpxxn/clank/internal/clanklog"
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
	testDynamicProto(false)
	os.Exit(m.Run())
}

var fileDescriptors []*desc.FileDescriptor

// parse proto
func testDynamicProto(startServ bool) {
	//fileDescriptors := []*desc.FileDescriptor{}
	goPath, ok := os.LookupEnv("GOPATH")
	if !ok {
		log.Fatal("GOPATH not found")
		return
	}
	goPath += "/src"
	clanklog.Infof("GOPATH: %s", goPath)
	parser := &protoparse.Parser{
		ImportPaths: []string{"./", goPath}, //goPath,

	}
	var err error
	//t.Log(desc.ResolveImport("protos/common.proto"))
	fileDescriptors, err = parser.ParseFiles(
		//"protos/model/students.proto",
		"protos/api/student_api.proto",
		//"github.com/lpxxn/clank/internal/testdata/protos/common.proto",
		//"protos/common.proto",
	)
	if err, ok := err.(protoparse.ErrorWithPos); ok {
		log.Println(err.GetPosition())
		log.Println(err.Error())
		log.Fatal(err.Error())
	}
	//registerFiles, err := protodesc.NewFiles(desc.ToFileDescriptorSet(fileDescriptors...))
	if err != nil {
		log.Fatal(err)
	}
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

	for _, fileDesc := range fileDescriptors {
		log.Printf("===============\nfileDesc: %v", fileDesc)
		log.Printf("package: %s, gopackage: %s", fileDesc.GetPackage(), fileDesc.AsFileDescriptorProto().GetOptions().GetGoPackage())
		log.Printf("Metadata: %s", fileDesc.GetName()) // protos/api/student_api.proto
		for _, msgDesc := range fileDesc.GetMessageTypes() {
			log.Printf("msgDesc: %v", msgDesc)
			msgDesc.AsProto().ProtoMessage()
			b, err := proto.Marshal(msgDesc.AsDescriptorProto())
			log.Println(err)
			log.Println(string(b))
		}
		for _, servDesc := range fileDesc.GetServices() {
			log.Printf("service info: %v", servDesc)
			log.Printf("service name: %s", servDesc.GetName())
			for _, methodInfo := range servDesc.GetMethods() {
				log.Printf("methods %+v", methodInfo)
				input := methodInfo.GetInputType()
				log.Println(input)
				output := methodInfo.GetOutputType()
				log.Println(output)
			}
		}
		CreateServiceDesc(fileDesc, startServ)
	}

}

func CreateServiceDesc(fileDesc *desc.FileDescriptor, startServ bool) {
	for _, servDescriptor := range fileDesc.GetServices() {
		serviceDesc := grpc.ServiceDesc{
			ServiceName: servDescriptor.GetFullyQualifiedName(), // api.StudentSrv
			Metadata:    fileDesc.GetName(),
		}

		for _, methodDescriptor := range servDescriptor.GetMethods() {
			unaryMethodMap[methodDescriptor.GetFullyQualifiedName()] = methodDescriptor
			log.Println(methodDescriptor.GetFullyQualifiedName())
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
				clanklog.Infof("method: %s run handler", methodDescriptor.GetName())
				methodDesc := grpc.MethodDesc{
					MethodName: methodDescriptor.GetName(),
					Handler:    createUnaryServerHandler(serviceDesc, methodDescriptor),
				}
				serviceDesc.Methods = append(serviceDesc.Methods, methodDesc)
			}
		}
		if startServ {
			grpcServ := grpc.NewServer()
			grpcServ.RegisterService(&serviceDesc, nil)

			listener, err := net.Listen("tcp", testAddress)
			if err != nil {
				panic(err)
			}
			reflection.Register(grpcServ)
			go grpcServ.Serve(listener)
		}
	}
}

/*
api.StudentSrv.NewStudent
api.StudentSrv.StudentByID
api.StudentSrv.AllStudent
api.StudentSrv.StudentInfo
api.StudentSrv.QueryStudents
*/
func getMethodDesc(methodName string) *desc.MethodDescriptor {
	if methodDesc, ok := unaryMethodMap[methodName]; ok {
		return methodDesc
	}
	return nil
}

func TestDynamicClient(t *testing.T) {
	cc, err := grpc.Dial(fmt.Sprintf(":%d", testPort), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic(fmt.Sprintf("Failed to create client to %d: %s", testPort, err.Error()))
	}
	defer cc.Close()

	stub := grpcdynamic.NewStub(cc)
	t.Log(stub)
	newStudentDesc := getMethodDesc("api.StudentSrv.NewStudent")
	t.Log(newStudentDesc)
	msgFactory := dynamic.NewMessageFactoryWithDefaults()
	inputParam := msgFactory.NewMessage(newStudentDesc.GetInputType())

	dynamicInputParam, _ := dynamic.AsDynamicMessage(inputParam)
	dynamicInputParam.UnmarshalJSON([]byte(`{"id":222,"name":"abc"}`))
	resp, err := stub.InvokeRpc(context.Background(), newStudentDesc, dynamicInputParam)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp.String())

	//resp, err = stub.InvokeRpc(context.Background(), newStudentDesc, &model.Student{
	//	Id:   111,
	//	Name: "abc",
	//})
	//if err != nil {
	//	t.Fatal(err)
	//}
	//t.Log(resp.String())

	// api.StudentSrv.StudentByID
	studentByIDDesc := getMethodDesc("api.StudentSrv.StudentByID")
	t.Log(studentByIDDesc)
	inputParam = msgFactory.NewMessage(studentByIDDesc.GetInputType())

	dynamicInputParam, _ = dynamic.AsDynamicMessage(inputParam)
	dynamicInputParam.UnmarshalJSON([]byte(`{"id":11}`))
	resp, err = stub.InvokeRpc(context.Background(), studentByIDDesc, dynamicInputParam)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp.String())

	dynamicInputParam.UnmarshalJSON([]byte(`{"id":456}`))
	resp, err = stub.InvokeRpc(context.Background(), studentByIDDesc, dynamicInputParam)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp.String())

	// api.StudentSrv.AllStudent
	allStudentDesc := getMethodDesc("api.StudentSrv.AllStudent")
	t.Log(studentByIDDesc)
	inputParam = msgFactory.NewMessage(allStudentDesc.GetInputType())
	serverStream, err := stub.InvokeRpcServerStream(context.Background(), allStudentDesc, inputParam)
	if err != nil {
		t.Fatal(err)
	}
	allStuRev, err := serverStream.RecvMsg()
	if err == io.EOF {
		t.Log("EOF")
	}
	if err != nil && err != io.EOF {
		t.Fatal(err)
	}
	t.Logf("%+v", allStuRev.String())
}

func TestRpcClient(t *testing.T) {
	//conn, err := grpc.Dial(testAddress, grpc.WithInsecure())
	//if err != nil {
	//	t.Fatal(err)
	//}
	//client := api.NewStudentSrvClient(conn)
	//result, err := client.NewStudent(context.Background(), &model.Student{
	//	Id:   111,
	//	Name: "abc",
	//})
	//
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//t.Log(result)
	// 有pb.go和没有还不一样。
	//rev2, err := client.StudentByID(context.Background(), &api.QueryStudent{
	//	Id: 11,
	//})
	//if err != nil {
	//	t.Fatal(err)
	//}
	//t.Log(rev2)
	//
	//rev2, err = client.StudentByID(context.Background(), &api.QueryStudent{
	//	Id: 456,
	//})
	//if err != nil {
	//	t.Fatal(err)
	//}
	//t.Log(rev2)
	//out := map[string]interface{}{}
	//err = conn.Invoke(context.Background(), "/api.StudentSrv/NewStudent", map[string]string{"name": "test"}, out)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//t.Log(out)

	//revStream1, err := client.AllStudent(context.Background(), &empty.Empty{})
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//allStu, err := revStream1.Recv()
	//if err == io.EOF {
	//	t.Log("EOF")
	//}
	//if err != nil {
	//	t.Fatal(err)
	//}
	//t.Logf("%+v", allStu)
}

var unaryMethodMap map[string]*desc.MethodDescriptor = make(map[string]*desc.MethodDescriptor)

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
		inPut, _ := dynamic.AsDynamicMessage(inputParam)
		inPutBody, _ := inPut.MarshalJSON()
		log.Printf("input body: %s", inPutBody)
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
		stream.SendMsg(dynamicOutput)
		return nil
	}
}
