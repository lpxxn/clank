package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
	_ "github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"github.com/lpxxn/clank/internal/clanklog"
	"github.com/tidwall/sjson"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

type gRpcServer struct {
	// map[serverName]map[methodName]methodDesc
	unaryMethodMap      map[string]map[string]gRpcMethodDesc
	streamMethodMap     map[string]map[string]gRpcStreamDesc
	rpcServiceDescGroup []*gRpcServiceDesc
	serverNames         map[string]struct{}

	GetOutputJson    func(serviceDesc grpc.ServiceDesc, methodDesc *desc.MethodDescriptor, jBody string) ([]byte, error)
	makeHttpCallback func(serviceDesc grpc.ServiceDesc, methodDesc *desc.MethodDescriptor, jBody string)
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

func (g *gRpcServer) StartWithPort(port int) error {
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
				StreamName:    methodDescriptor.GetName(),
				Handler:       g.createStreamHandler(*rev.ServiceDesc, methodDescriptor),
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
		clanklog.Info(serviceDesc.ServiceName)
		clanklog.Info(methodDesc.GetName())
		clanklog.Info(srv)
		msgFactory := dynamic.NewMessageFactoryWithDefaults()
		inputParam := msgFactory.NewMessage(methodDesc.GetInputType())
		if err := dec(inputParam); err != nil {
			return nil, err
		}
		clanklog.Info(inputParam.String())

		outPut := msgFactory.NewMessage(methodDesc.GetOutputType())
		dynamicOutput, err := dynamic.AsDynamicMessage(outPut)
		if err != nil {
			return nil, err
		}

		ctx = g.setRequestJBody(ctx, inputParam)
		jBody := g.getJBody(ctx)
		outputJson, err := g.GetOutputJson(serviceDesc, methodDesc, jBody)
		clanklog.Infof("output: %s", string(outputJson))
		if err != nil {
			return nil, err
		}
		if err := dynamicOutput.UnmarshalJSON(outputJson); err != nil {
			return nil, err
		}

		jBody, err = sjson.SetRaw(jBody, "response", string(outputJson))
		if err != nil {
			clanklog.Errorf("commonHandler sjson.Set error: %s", err.Error())
		}
		defer g.makeHttpCallback(serviceDesc, methodDesc, jBody)
		return dynamicOutput, nil
	}
}

func (g *gRpcServer) getRequestHeaders(md metadata.MD) map[string]string {
	headers := make(map[string]string)
	for k, v := range md {
		headers[k] = v[0]
	}
	return headers
}

func (g *gRpcServer) setRequestJBody(ctx context.Context, body proto.Message) (rev context.Context) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}
	jBody := ``
	var err error
	headers := g.getRequestHeaders(md)
	if len(headers) > 0 {
		jBody, err = sjson.Set(jBody, `header`, headers)
		if err != nil {
			clanklog.Errorf("set request header jBody error: %v", err)
		}
	}

	jBody, err = sjson.Set(jBody, grpcRequestParam, body)
	if err != nil {
		clanklog.Errorf("setRequestParamJBody error: %v", err)
	}
	// set metadata
	md["customer-body"] = []string{jBody}
	return metadata.NewIncomingContext(ctx, md)
}

func (g *gRpcServer) getJBody(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ``
	}
	return md["customer-body"][0]
}

func (g *gRpcServer) createStreamHandler(serviceDesc grpc.ServiceDesc, methodDesc *desc.MethodDescriptor) func(srv interface{}, stream grpc.ServerStream) error {
	return func(srv interface{}, stream grpc.ServerStream) error {
		isServerStream := methodDesc.IsServerStreaming()
		isClientStream := methodDesc.IsClientStreaming()
		clanklog.Info(isServerStream)
		clanklog.Info(isClientStream)
		ctx := stream.Context()
		msgFactory := dynamic.NewMessageFactoryWithDefaults()
		inputType := methodDesc.GetInputType()
		inputParam := msgFactory.NewMessage(inputType)
		if err := stream.RecvMsg(inputParam); err != nil {
			return err
		}
		clanklog.Info(inputParam.String())

		ctx = g.setRequestJBody(ctx, inputParam)

		outPut := msgFactory.NewMessage(methodDesc.GetOutputType())
		dynamicOutput, err := dynamic.AsDynamicMessage(outPut)
		if err != nil {
			return err
		}
		jBody := g.getJBody(ctx)
		outputJson, err := g.GetOutputJson(serviceDesc, methodDesc, jBody)
		clanklog.Info(string(outputJson))
		if err := dynamicOutput.UnmarshalJSON([]byte(outputJson)); err != nil {
			return err
		}
		jBody, err = sjson.SetRaw(jBody, "response", string(outputJson))
		if err != nil {
			clanklog.Errorf("commonHandler sjson.Set error: %s", err.Error())
		}
		defer g.makeHttpCallback(serviceDesc, methodDesc, jBody)
		return stream.SendMsg(dynamicOutput)
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

			clanklog.Info(json.Marshal(inputParam)) // {}
			dynamicMsg, err := dynamic.AsDynamicMessage(inputParam)
			if err != nil {
				return err
			}
			jsonBody, err := dynamicMsg.MarshalJSON()
			if err != nil {
				return err
			}
			clanklog.Info(string(jsonBody))

			for _, v := range conditionParameters {
				if !strings.Contains(v, ".") {
					if _, err := dynamicMsg.TryGetFieldByName(v); err != nil {
						return err
					}
				}
			}
		}

		streamMethod, ok := g.streamMethodMap[serverSchema.Name][item.Name]
		if ok {
			_ = streamMethod
			continue
		}
		//return fmt.Errorf("invalid method name %s", item.Name)
	}
	return nil
}

func ValidateGrpcServiceInputAndOutput(schemaList GrpcServerDescriptionList, gRpcServ *gRpcServer) error {
	for _, item := range schemaList {
		if err := gRpcServ.ValidateSchemaMethod(item); err != nil {
			return err
		}
	}
	return nil
}

func SetOutputFunc(schemaList GrpcServerDescriptionList, gRpcServ *gRpcServer) error {
	gRpcServ.GetOutputJson = func(serviceDesc grpc.ServiceDesc, methodDesc *desc.MethodDescriptor, jBody string) ([]byte, error) {
		methodSchema, err := schemaList.GetMethod(serviceDesc.ServiceName, methodDesc.GetName())
		if err != nil {
			return nil, fmt.Errorf("server: %s method: %s, err: [%w]", serviceDesc.ServiceName, methodDesc.GetName(), err)
		}

		if len(methodSchema.Conditions) == 0 {
			return GenerateDefaultTemplate(methodSchema.DefaultResponse)
		}
		for _, condition := range methodSchema.Conditions {
			conditionStr := condition.Condition
			paramValue, err := ParamValue(condition.Parameters, jBody)
			if err != nil {
				clanklog.Errorf("get condition param value error: %s", err)
				continue
			}
			if len(paramValue) != len(condition.Parameters) {
				continue
			}
			for k, v := range paramValue {
				conditionStr = strings.ReplaceAll(conditionStr, "$"+k, fmt.Sprintf("%v", v))
			}

			clanklog.Info("conditionStr", conditionStr)
			result, err := ValuableBoolExpression(conditionStr)
			if err != nil {
				return nil, err
			}

			if result == true {
				return GenerateDefaultTemplate(condition.Response)
			}
		}
		return GenerateDefaultTemplate(methodSchema.DefaultResponse)
	}

	gRpcServ.makeHttpCallback = func(serviceDesc grpc.ServiceDesc, methodDesc *desc.MethodDescriptor, jBody string) {
		methodSchema, err := schemaList.GetMethod(serviceDesc.ServiceName, methodDesc.GetName())
		if err != nil {
			clanklog.Errorf("server: %s method: %s, err: [%+v]", serviceDesc.ServiceName, methodDesc.GetName(), err)
			return
		}
		if len(methodSchema.HttpCallback) == 0 {
			return
		}
		for _, callback := range methodSchema.HttpCallback {
			delayTime := callback.DelayTime
			if delayTime <= 0 {
				delayTime = 1
			}
			time.AfterFunc(time.Duration(delayTime)*time.Second, func() {
				if err := callback.makeRequest(context.Background(), jBody); err != nil {
					clanklog.Errorf("callback err: %+v", err)
				}
			})
		}
	}
	return nil
}

func keysInterfaceSlice(k string) []interface{} {
	keys := strings.Split(k, ".")
	var keyList []interface{}
	for _, v := range keys {
		if i, err := strconv.Atoi(v); err != nil {
			keyList = append(keyList, v)
		} else {
			keyList = append(keyList, i)
		}
	}
	clanklog.Info("keyList", keyList)
	return keyList
}
