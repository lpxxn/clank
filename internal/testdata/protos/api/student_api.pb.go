// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.22.0
// 	protoc        v3.13.0
// source: protos/api/student_api.proto

package api

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
	protos "github.com/lpxxn/clank/internal/testdata/protos"
	model "github.com/lpxxn/clank/internal/testdata/protos/model"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type QueryStudent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *QueryStudent) Reset() {
	*x = QueryStudent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_api_student_api_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueryStudent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueryStudent) ProtoMessage() {}

func (x *QueryStudent) ProtoReflect() protoreflect.Message {
	mi := &file_protos_api_student_api_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueryStudent.ProtoReflect.Descriptor instead.
func (*QueryStudent) Descriptor() ([]byte, []int) {
	return file_protos_api_student_api_proto_rawDescGZIP(), []int{0}
}

func (x *QueryStudent) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type QueryStudentResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StudentList []*model.Student `protobuf:"bytes,1,rep,name=studentList,proto3" json:"studentList,omitempty"`
}

func (x *QueryStudentResponse) Reset() {
	*x = QueryStudentResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_api_student_api_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueryStudentResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueryStudentResponse) ProtoMessage() {}

func (x *QueryStudentResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protos_api_student_api_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueryStudentResponse.ProtoReflect.Descriptor instead.
func (*QueryStudentResponse) Descriptor() ([]byte, []int) {
	return file_protos_api_student_api_proto_rawDescGZIP(), []int{1}
}

func (x *QueryStudentResponse) GetStudentList() []*model.Student {
	if x != nil {
		return x.StudentList
	}
	return nil
}

var File_protos_api_student_api_proto protoreflect.FileDescriptor

var file_protos_api_student_api_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x73, 0x74, 0x75,
	0x64, 0x65, 0x6e, 0x74, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03,
	0x61, 0x70, 0x69, 0x1a, 0x3c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x6c, 0x70, 0x78, 0x78, 0x6e, 0x2f, 0x63, 0x6c, 0x61, 0x6e, 0x6b, 0x2f, 0x69, 0x6e, 0x74, 0x65,
	0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x1b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2f,
	0x73, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x1e, 0x0a, 0x0c, 0x51,
	0x75, 0x65, 0x72, 0x79, 0x53, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x22, 0x48, 0x0a, 0x14, 0x51,
	0x75, 0x65, 0x72, 0x79, 0x53, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x30, 0x0a, 0x0b, 0x73, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x4c, 0x69,
	0x73, 0x74, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c,
	0x2e, 0x53, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x52, 0x0b, 0x73, 0x74, 0x75, 0x64, 0x65, 0x6e,
	0x74, 0x4c, 0x69, 0x73, 0x74, 0x32, 0xfb, 0x01, 0x0a, 0x0a, 0x53, 0x74, 0x75, 0x64, 0x65, 0x6e,
	0x74, 0x53, 0x72, 0x76, 0x12, 0x2c, 0x0a, 0x0a, 0x4e, 0x65, 0x77, 0x53, 0x74, 0x75, 0x64, 0x65,
	0x6e, 0x74, 0x12, 0x0e, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x53, 0x74, 0x75, 0x64, 0x65,
	0x6e, 0x74, 0x1a, 0x0e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x52, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x12, 0x3b, 0x0a, 0x0b, 0x53, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x42, 0x79, 0x49,
	0x44, 0x12, 0x11, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x53, 0x74, 0x75,
	0x64, 0x65, 0x6e, 0x74, 0x1a, 0x19, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79,
	0x53, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x41, 0x0a, 0x0a, 0x41, 0x6c, 0x6c, 0x53, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x12, 0x16, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x19, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x51, 0x75, 0x65, 0x72,
	0x79, 0x53, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x30, 0x01, 0x12, 0x3f, 0x0a, 0x0b, 0x53, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x49, 0x6e, 0x66,
	0x6f, 0x12, 0x11, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x53, 0x74, 0x75,
	0x64, 0x65, 0x6e, 0x74, 0x1a, 0x19, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79,
	0x53, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x28,
	0x01, 0x30, 0x01, 0x42, 0x44, 0x0a, 0x0d, 0x63, 0x6f, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x61, 0x70, 0x69, 0x5a, 0x33, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x6c, 0x70, 0x78, 0x78, 0x6e, 0x2f, 0x63, 0x6c, 0x61, 0x6e, 0x6b, 0x2f, 0x69, 0x6e, 0x74,
	0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_protos_api_student_api_proto_rawDescOnce sync.Once
	file_protos_api_student_api_proto_rawDescData = file_protos_api_student_api_proto_rawDesc
)

func file_protos_api_student_api_proto_rawDescGZIP() []byte {
	file_protos_api_student_api_proto_rawDescOnce.Do(func() {
		file_protos_api_student_api_proto_rawDescData = protoimpl.X.CompressGZIP(file_protos_api_student_api_proto_rawDescData)
	})
	return file_protos_api_student_api_proto_rawDescData
}

var file_protos_api_student_api_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_protos_api_student_api_proto_goTypes = []interface{}{
	(*QueryStudent)(nil),         // 0: api.QueryStudent
	(*QueryStudentResponse)(nil), // 1: api.QueryStudentResponse
	(*model.Student)(nil),        // 2: model.Student
	(*empty.Empty)(nil),          // 3: google.protobuf.Empty
	(*protos.Result)(nil),        // 4: protos.Result
}
var file_protos_api_student_api_proto_depIdxs = []int32{
	2, // 0: api.QueryStudentResponse.studentList:type_name -> model.Student
	2, // 1: api.StudentSrv.NewStudent:input_type -> model.Student
	0, // 2: api.StudentSrv.StudentByID:input_type -> api.QueryStudent
	3, // 3: api.StudentSrv.AllStudent:input_type -> google.protobuf.Empty
	0, // 4: api.StudentSrv.StudentInfo:input_type -> api.QueryStudent
	4, // 5: api.StudentSrv.NewStudent:output_type -> protos.Result
	1, // 6: api.StudentSrv.StudentByID:output_type -> api.QueryStudentResponse
	1, // 7: api.StudentSrv.AllStudent:output_type -> api.QueryStudentResponse
	1, // 8: api.StudentSrv.StudentInfo:output_type -> api.QueryStudentResponse
	5, // [5:9] is the sub-list for method output_type
	1, // [1:5] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_protos_api_student_api_proto_init() }
func file_protos_api_student_api_proto_init() {
	if File_protos_api_student_api_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protos_api_student_api_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QueryStudent); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_protos_api_student_api_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QueryStudentResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_protos_api_student_api_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_protos_api_student_api_proto_goTypes,
		DependencyIndexes: file_protos_api_student_api_proto_depIdxs,
		MessageInfos:      file_protos_api_student_api_proto_msgTypes,
	}.Build()
	File_protos_api_student_api_proto = out.File
	file_protos_api_student_api_proto_rawDesc = nil
	file_protos_api_student_api_proto_goTypes = nil
	file_protos_api_student_api_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// StudentSrvClient is the client API for StudentSrv service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type StudentSrvClient interface {
	NewStudent(ctx context.Context, in *model.Student, opts ...grpc.CallOption) (*protos.Result, error)
	StudentByID(ctx context.Context, in *QueryStudent, opts ...grpc.CallOption) (*QueryStudentResponse, error)
	AllStudent(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (StudentSrv_AllStudentClient, error)
	StudentInfo(ctx context.Context, opts ...grpc.CallOption) (StudentSrv_StudentInfoClient, error)
}

type studentSrvClient struct {
	cc grpc.ClientConnInterface
}

func NewStudentSrvClient(cc grpc.ClientConnInterface) StudentSrvClient {
	return &studentSrvClient{cc}
}

func (c *studentSrvClient) NewStudent(ctx context.Context, in *model.Student, opts ...grpc.CallOption) (*protos.Result, error) {
	out := new(protos.Result)
	err := c.cc.Invoke(ctx, "/api.StudentSrv/NewStudent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *studentSrvClient) StudentByID(ctx context.Context, in *QueryStudent, opts ...grpc.CallOption) (*QueryStudentResponse, error) {
	out := new(QueryStudentResponse)
	err := c.cc.Invoke(ctx, "/api.StudentSrv/StudentByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *studentSrvClient) AllStudent(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (StudentSrv_AllStudentClient, error) {
	stream, err := c.cc.NewStream(ctx, &_StudentSrv_serviceDesc.Streams[0], "/api.StudentSrv/AllStudent", opts...)
	if err != nil {
		return nil, err
	}
	x := &studentSrvAllStudentClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type StudentSrv_AllStudentClient interface {
	Recv() (*QueryStudentResponse, error)
	grpc.ClientStream
}

type studentSrvAllStudentClient struct {
	grpc.ClientStream
}

func (x *studentSrvAllStudentClient) Recv() (*QueryStudentResponse, error) {
	m := new(QueryStudentResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *studentSrvClient) StudentInfo(ctx context.Context, opts ...grpc.CallOption) (StudentSrv_StudentInfoClient, error) {
	stream, err := c.cc.NewStream(ctx, &_StudentSrv_serviceDesc.Streams[1], "/api.StudentSrv/StudentInfo", opts...)
	if err != nil {
		return nil, err
	}
	x := &studentSrvStudentInfoClient{stream}
	return x, nil
}

type StudentSrv_StudentInfoClient interface {
	Send(*QueryStudent) error
	Recv() (*QueryStudentResponse, error)
	grpc.ClientStream
}

type studentSrvStudentInfoClient struct {
	grpc.ClientStream
}

func (x *studentSrvStudentInfoClient) Send(m *QueryStudent) error {
	return x.ClientStream.SendMsg(m)
}

func (x *studentSrvStudentInfoClient) Recv() (*QueryStudentResponse, error) {
	m := new(QueryStudentResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// StudentSrvServer is the server API for StudentSrv service.
type StudentSrvServer interface {
	NewStudent(context.Context, *model.Student) (*protos.Result, error)
	StudentByID(context.Context, *QueryStudent) (*QueryStudentResponse, error)
	AllStudent(*empty.Empty, StudentSrv_AllStudentServer) error
	StudentInfo(StudentSrv_StudentInfoServer) error
}

// UnimplementedStudentSrvServer can be embedded to have forward compatible implementations.
type UnimplementedStudentSrvServer struct {
}

func (*UnimplementedStudentSrvServer) NewStudent(context.Context, *model.Student) (*protos.Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewStudent not implemented")
}
func (*UnimplementedStudentSrvServer) StudentByID(context.Context, *QueryStudent) (*QueryStudentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StudentByID not implemented")
}
func (*UnimplementedStudentSrvServer) AllStudent(*empty.Empty, StudentSrv_AllStudentServer) error {
	return status.Errorf(codes.Unimplemented, "method AllStudent not implemented")
}
func (*UnimplementedStudentSrvServer) StudentInfo(StudentSrv_StudentInfoServer) error {
	return status.Errorf(codes.Unimplemented, "method StudentInfo not implemented")
}

func RegisterStudentSrvServer(s *grpc.Server, srv StudentSrvServer) {
	s.RegisterService(&_StudentSrv_serviceDesc, srv)
}

func _StudentSrv_NewStudent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(model.Student)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StudentSrvServer).NewStudent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.StudentSrv/NewStudent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StudentSrvServer).NewStudent(ctx, req.(*model.Student))
	}
	return interceptor(ctx, in, info, handler)
}

func _StudentSrv_StudentByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryStudent)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StudentSrvServer).StudentByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.StudentSrv/StudentByID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StudentSrvServer).StudentByID(ctx, req.(*QueryStudent))
	}
	return interceptor(ctx, in, info, handler)
}

func _StudentSrv_AllStudent_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(empty.Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StudentSrvServer).AllStudent(m, &studentSrvAllStudentServer{stream})
}

type StudentSrv_AllStudentServer interface {
	Send(*QueryStudentResponse) error
	grpc.ServerStream
}

type studentSrvAllStudentServer struct {
	grpc.ServerStream
}

func (x *studentSrvAllStudentServer) Send(m *QueryStudentResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _StudentSrv_StudentInfo_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(StudentSrvServer).StudentInfo(&studentSrvStudentInfoServer{stream})
}

type StudentSrv_StudentInfoServer interface {
	Send(*QueryStudentResponse) error
	Recv() (*QueryStudent, error)
	grpc.ServerStream
}

type studentSrvStudentInfoServer struct {
	grpc.ServerStream
}

func (x *studentSrvStudentInfoServer) Send(m *QueryStudentResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *studentSrvStudentInfoServer) Recv() (*QueryStudent, error) {
	m := new(QueryStudent)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _StudentSrv_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.StudentSrv",
	HandlerType: (*StudentSrvServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "NewStudent",
			Handler:    _StudentSrv_NewStudent_Handler,
		},
		{
			MethodName: "StudentByID",
			Handler:    _StudentSrv_StudentByID_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "AllStudent",
			Handler:       _StudentSrv_AllStudent_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "StudentInfo",
			Handler:       _StudentSrv_StudentInfo_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "protos/api/student_api.proto",
}
