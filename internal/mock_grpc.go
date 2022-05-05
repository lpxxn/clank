package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"

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
	unaryMethodMap      map[string]map[string]gRpcMethodDesc
	streamMethodMap     map[string]map[string]gRpcStreamDesc
	rpcServiceDescGroup []*gRpcServiceDesc
	serverNames         map[string]struct{}
}

type MockGrpcResponse struct {
	serverName string
	methodName string
	respBody   map[string]struct{}
	conditions string
}

type gRpcServiceDesc struct {
	*grpc.ServiceDesc
}

type gRpcMethodDesc struct {
	*grpc.MethodDesc
	methodDescriptor *desc.MethodDescriptor
}

type gRpcStreamDesc struct {
	*grpc.StreamDesc
	methodDescriptor *desc.MethodDescriptor
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
		unaryMethodMap:  make(map[string]map[string]gRpcMethodDesc),
		streamMethodMap: map[string]map[string]gRpcStreamDesc{},
	}
	rev.extractServicesInfo(fileDesc...)
	return rev
}

func (g *gRpcServer) Start(port int) error {
	grpcServ := grpc.NewServer()
	for _, item := range g.rpcServiceDescGroup {
		grpcServ.RegisterService(item.ServiceDesc, nil)
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	reflection.Register(grpcServ)
	return grpcServ.Serve(listener)
}

func (g *gRpcServer) extractServicesInfo(fileDescList ...*desc.FileDescriptor) {
	for _, fileDesc := range fileDescList {
		for _, servDescriptor := range fileDesc.GetServices() {
			g.serverNames[servDescriptor.GetFullyQualifiedName()] = struct{}{}
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
	g.unaryMethodMap[rev.ServiceName] = make(map[string]gRpcMethodDesc)
	g.streamMethodMap[rev.ServiceName] = make(map[string]gRpcStreamDesc)
	for _, methodDescriptor := range servDescriptor.GetMethods() {
		isServerStream := methodDescriptor.IsServerStreaming()
		isClientStream := methodDescriptor.IsClientStreaming()
		if isServerStream || isClientStream {
			streamDesc := gRpcStreamDesc{StreamDesc: &grpc.StreamDesc{
				StreamName: methodDescriptor.GetName(),
				// TODO: // wait a moment
				Handler:       nil,
				ServerStreams: isServerStream,
				ClientStreams: isClientStream,
			}, methodDescriptor: methodDescriptor}
			rev.Streams = append(rev.Streams, *streamDesc.StreamDesc)
			g.streamMethodMap[rev.ServiceName][methodDescriptor.GetName()] = streamDesc
		} else {
			methodDesc := gRpcMethodDesc{MethodDesc: &grpc.MethodDesc{
				MethodName: methodDescriptor.GetName(),
				Handler:    g.createUnaryServerHandler(*rev.ServiceDesc, methodDescriptor),
			},
				methodDescriptor: methodDescriptor,
			}
			rev.Methods = append(rev.Methods, *methodDesc.MethodDesc)
			g.unaryMethodMap[rev.ServiceName][methodDesc.MethodName] = methodDesc
		}
	}
	return rev
}

func (g *gRpcServer) createUnaryServerHandler(serviceDesc grpc.ServiceDesc, methodDesc *desc.MethodDescriptor) func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
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

func (g *gRpcServer) ValidateSchemaMethod(serverSchema *GrpcServerDescription) error {
	if _, ok := g.serverNames[serverSchema.Name]; !ok {
		return fmt.Errorf("invalid server name %s", serverSchema.Name)
	}
	for _, item := range serverSchema.Methods {
		unaryMethod, ok := g.unaryMethodMap[serverSchema.Name][item.Name]
		conditionParameters := item.Parameters
		if ok {
			_ = unaryMethod
			msgFactory := dynamic.NewMessageFactoryWithDefaults()
			inputParam := msgFactory.NewMessage(unaryMethod.methodDescriptor.GetInputType())

			fmt.Println(json.Marshal(inputParam)) // {}
			dynamicMsg, err := dynamic.AsDynamicMessage(inputParam)
			if err != nil {
				return err
			}
			jsonBody, err := dynamicMsg.MarshalJSON()
			if err != nil {
				return err
			}
			fmt.Println(string(jsonBody))

			for _, v := range conditionParameters {
				if !strings.Contains(v, ".") {
					if _, err := dynamicMsg.TryGetFieldByName(v); err != nil {
						return err
					}
				}
			}
			outPut := msgFactory.NewMessage(unaryMethod.methodDescriptor.GetOutputType())
			dynamicOutput, err := dynamic.AsDynamicMessage(outPut)
			if err != nil {
				return err
			}

			if item.DefaultResponse != "" {
				errors.New("default response is not supported")

				if err := dynamicOutput.UnmarshalJSON([]byte(item.DefaultResponse)); err != nil {
					return fmt.Errorf("server: %s method: %s, invalid default response %s, err: [%w]", serverSchema.Name, item.Name, item.DefaultResponse, err)
				}
			}

			for _, v := range item.Conditions {
				if err := dynamicOutput.UnmarshalJSON([]byte(v.Response)); err != nil {
					return fmt.Errorf("server: %s method: %s, invalid condition response %s, err: [%w]", serverSchema.Name, item.Name, v.Response, err)
				}
			}
			continue
		}

		streamMethod, ok := g.streamMethodMap[serverSchema.Name][item.Name]
		if ok {
			_ = streamMethod
			continue
		}
		return fmt.Errorf("invalid method name %s", item.Name)
	}
	return nil
}

func ValidateServiceInputAndOutput(schemaList ServerList, gRpcServ *gRpcServer) error {
	grpcServersSchema := ServerDescriptionList{}
	for _, server := range schemaList {
		if s, ok := server.(*GrpcServerDescription); ok {
			grpcServersSchema = append(grpcServersSchema, s)
		} else {
			return fmt.Errorf("invalid server type %T, need *GrpcServerDescription type", server)
		}
	}
	for _, item := range grpcServersSchema {
		if err := gRpcServ.ValidateSchemaMethod(item); err != nil {
			return err
		}
	}
	return nil
}
