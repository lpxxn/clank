package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"github.com/golang/protobuf/proto"
	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type gRpcServer struct {
	// map[serverName]map[methodName]methodDesc
	unaryMethodMap      map[string]map[string]grpc.MethodDesc
	streamMethodMap     map[string]map[string]grpc.StreamDesc
	rpcServiceDescGroup []*gRpcServiceDesc
	serverNames         map[string]struct{}
}

type gRpcServiceDesc struct {
	*grpc.ServiceDesc
}

func ParseServerMethodsFromProto(importPath []string, filePath []string) (*gRpcServer, error) {
	fileDesc, err := ParseProtoFile(importPath, filePath)
	if err != nil {
		return nil, err
	}
	return ParseServerMethodsFromFileDescriptor(fileDesc...), nil
}

func ParseServerMethodsFromProtoset(filePath string) (*gRpcServer, error) {
	fileDesc, err := ParseProtoFileFromProtoset(filePath)
	if err != nil {
		return nil, err
	}
	return ParseServerMethodsFromFileDescriptor(fileDesc), nil
}

func ParseProtoFile(importPath []string, filePath []string) ([]*desc.FileDescriptor, error) {
	goPath, ok := os.LookupEnv("GOPATH")
	if ok {
		importPath = append(importPath, goPath+"/src/")
	}
	parser := &protoparse.Parser{
		ImportPaths: importPath,
	}
	return parser.ParseFiles(
		filePath...,
	)
}

func ParseProtoFileFromProtoset(protosetPath string) (*desc.FileDescriptor, error) {
	var fds dpb.FileDescriptorSet
	f, err := os.Open(protosetPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	bb, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	if err = proto.Unmarshal(bb, &fds); err != nil {
		return nil, err
	}
	return desc.CreateFileDescriptorFromSet(&fds)
}

func ParseServerMethodsFromFileDescriptor(fileDesc ...*desc.FileDescriptor) *gRpcServer {
	rev := &gRpcServer{
		serverNames:     make(map[string]struct{}),
		unaryMethodMap:  make(map[string]map[string]grpc.MethodDesc),
		streamMethodMap: map[string]map[string]grpc.StreamDesc{},
	}
	rev.extractServicesInfo(fileDesc...)
	return rev
}

func (g *gRpcServer) Start(port int) {
	grpcServ := grpc.NewServer()
	for _, item := range g.rpcServiceDescGroup {
		grpcServ.RegisterService(item.ServiceDesc, nil)
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}
	reflection.Register(grpcServ)
	grpcServ.Serve(listener)
}

func (g *gRpcServer) extractServicesInfo(fileDescList ...*desc.FileDescriptor) {
	for _, fileDesc := range fileDescList {
		for _, servDescriptor := range fileDesc.GetServices() {
			g.serverNames[servDescriptor.GetName()] = struct{}{}
			g.rpcServiceDescGroup = append(g.rpcServiceDescGroup, g.methodDesc(servDescriptor))
		}
	}
}

func (g *gRpcServer) methodDesc(servDescriptor *desc.ServiceDescriptor) *gRpcServiceDesc {
	rev := &gRpcServiceDesc{
		ServiceDesc: &grpc.ServiceDesc{
			ServiceName: servDescriptor.GetFullyQualifiedName(),
			Metadata:    servDescriptor.GetFile().GetName(),
		},
	}
	g.unaryMethodMap[rev.ServiceName] = make(map[string]grpc.MethodDesc)
	g.streamMethodMap[rev.ServiceName] = make(map[string]grpc.StreamDesc)
	for _, methodDescriptor := range servDescriptor.GetMethods() {
		isServerStream := methodDescriptor.IsServerStreaming()
		isClientStream := methodDescriptor.IsClientStreaming()
		if isServerStream || isClientStream {
			streamDesc := grpc.StreamDesc{
				StreamName: methodDescriptor.GetName(),
				// TODO: // wait a moment
				Handler:       nil,
				ServerStreams: isServerStream,
				ClientStreams: isClientStream,
			}
			rev.Streams = append(rev.Streams, streamDesc)
			g.streamMethodMap[rev.ServiceName][methodDescriptor.GetName()] = streamDesc
		} else {
			methodDesc := grpc.MethodDesc{
				MethodName: methodDescriptor.GetName(),
				Handler:    createUnaryServerHandler(*rev.ServiceDesc, methodDescriptor),
			}
			rev.Methods = append(rev.Methods, methodDesc)
			g.unaryMethodMap[rev.ServiceName][methodDesc.MethodName] = methodDesc
		}
	}
	return rev
}

func createUnaryServerHandler(serviceDesc grpc.ServiceDesc, methodDesc *desc.MethodDescriptor) func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	return func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
		fmt.Println(serviceDesc.ServiceName)
		fmt.Println(methodDesc.GetName())
		fmt.Println(srv)
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
	}
}
