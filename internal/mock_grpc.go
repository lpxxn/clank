package internal

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"google.golang.org/grpc"
)

type gRpcServer struct {
	// map[serverName]map[methodName]methodDesc
	unaryMethodMap map[string]map[string]grpc.MethodDesc
	rpcServiceDesc *gRpcServiceDesc
}

type gRpcServiceDesc struct {
	*grpc.ServiceDesc
	isStream bool
}

func NewGRpcServer() *gRpcServer {
	return &gRpcServer{
		unaryMethodMap: make(map[string]map[string]grpc.MethodDesc),
	}
}

func (g *gRpcServer) ExtractMethods(fileDesc *desc.FileDescriptor) {
	for _, servDescriptor := range fileDesc.GetServices() {
		g.rpcServiceDesc = g.methodDesc(servDescriptor)
	}
}

func (g *gRpcServer) methodDesc(servDescriptor *desc.ServiceDescriptor) *gRpcServiceDesc {
	rev := &gRpcServiceDesc{
		ServiceDesc: &grpc.ServiceDesc{
			ServiceName: servDescriptor.GetFullyQualifiedName(),
			Metadata:    servDescriptor.GetFile().GetName(),
		},
		isStream: false,
	}
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
		} else {
			methodDesc := grpc.MethodDesc{
				MethodName: methodDescriptor.GetName(),
				Handler:    createUnaryServerHandler(*rev.ServiceDesc, methodDescriptor),
			}
			rev.Methods = append(rev.Methods, methodDesc)
		}
	}
	return rev

}

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
	}
}
