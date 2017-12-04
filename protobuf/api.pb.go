// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api.proto

/*
Package API is a generated protocol buffer package.

It is generated from these files:
	api.proto

It has these top-level messages:
	Task
	Meta
	RequestMessage
	StatusMessage
	QueueProgress
*/
package API

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "google.golang.org/genproto/googleapis/api/annotations"
import google_protobuf1 "github.com/golang/protobuf/ptypes/empty"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Status int32

const (
	Status_SUCCESS Status = 0
	Status_FAILURE Status = 1
	Status_UNKNOWN Status = 2
)

var Status_name = map[int32]string{
	0: "SUCCESS",
	1: "FAILURE",
	2: "UNKNOWN",
}
var Status_value = map[string]int32{
	"SUCCESS": 0,
	"FAILURE": 1,
	"UNKNOWN": 2,
}

func (x Status) String() string {
	return proto.EnumName(Status_name, int32(x))
}
func (Status) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

// Task is the definition of a task
type Task struct {
	// reserved 1; // we took queueID out
	// reserved "queueID";
	QueueID       string   `protobuf:"bytes,1,opt,name=queueID" json:"queueID,omitempty"`
	TaskID        string   `protobuf:"bytes,2,opt,name=taskID" json:"taskID,omitempty"`
	Body          string   `protobuf:"bytes,3,opt,name=body" json:"body,omitempty"`
	Envs          []string `protobuf:"bytes,4,rep,name=envs" json:"envs,omitempty"`
	ExpirationSec int32    `protobuf:"varint,5,opt,name=expirationSec" json:"expirationSec,omitempty"`
}

func (m *Task) Reset()                    { *m = Task{} }
func (m *Task) String() string            { return proto.CompactTextString(m) }
func (*Task) ProtoMessage()               {}
func (*Task) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Task) GetQueueID() string {
	if m != nil {
		return m.QueueID
	}
	return ""
}

func (m *Task) GetTaskID() string {
	if m != nil {
		return m.TaskID
	}
	return ""
}

func (m *Task) GetBody() string {
	if m != nil {
		return m.Body
	}
	return ""
}

func (m *Task) GetEnvs() []string {
	if m != nil {
		return m.Envs
	}
	return nil
}

func (m *Task) GetExpirationSec() int32 {
	if m != nil {
		return m.ExpirationSec
	}
	return 0
}

type Meta struct {
	QueueID string `protobuf:"bytes,1,opt,name=queueID" json:"queueID,omitempty"`
}

func (m *Meta) Reset()                    { *m = Meta{} }
func (m *Meta) String() string            { return proto.CompactTextString(m) }
func (*Meta) ProtoMessage()               {}
func (*Meta) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Meta) GetQueueID() string {
	if m != nil {
		return m.QueueID
	}
	return ""
}

type RequestMessage struct {
	QueueID       string `protobuf:"bytes,1,opt,name=queueID" json:"queueID,omitempty"`
	ExpirationSec int32  `protobuf:"varint,2,opt,name=expirationSec" json:"expirationSec,omitempty"`
}

func (m *RequestMessage) Reset()                    { *m = RequestMessage{} }
func (m *RequestMessage) String() string            { return proto.CompactTextString(m) }
func (*RequestMessage) ProtoMessage()               {}
func (*RequestMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *RequestMessage) GetQueueID() string {
	if m != nil {
		return m.QueueID
	}
	return ""
}

func (m *RequestMessage) GetExpirationSec() int32 {
	if m != nil {
		return m.ExpirationSec
	}
	return 0
}

type StatusMessage struct {
	Status Status `protobuf:"varint,1,opt,name=status,enum=API.Status" json:"status,omitempty"`
	Detail string `protobuf:"bytes,2,opt,name=detail" json:"detail,omitempty"`
}

func (m *StatusMessage) Reset()                    { *m = StatusMessage{} }
func (m *StatusMessage) String() string            { return proto.CompactTextString(m) }
func (*StatusMessage) ProtoMessage()               {}
func (*StatusMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *StatusMessage) GetStatus() Status {
	if m != nil {
		return m.Status
	}
	return Status_SUCCESS
}

func (m *StatusMessage) GetDetail() string {
	if m != nil {
		return m.Detail
	}
	return ""
}

type QueueProgress struct {
	IncomingListLength int32   `protobuf:"varint,1,opt,name=incomingListLength" json:"incomingListLength,omitempty"`
	ReceivedListLength int32   `protobuf:"varint,2,opt,name=receivedListLength" json:"receivedListLength,omitempty"`
	NonExpiredCount    int32   `protobuf:"varint,3,opt,name=nonExpiredCount" json:"nonExpiredCount,omitempty"`
	ExpiredCount       int32   `protobuf:"varint,4,opt,name=expiredCount" json:"expiredCount,omitempty"`
	CompletedCount     int32   `protobuf:"varint,5,opt,name=completedCount" json:"completedCount,omitempty"`
	FailedCount        int32   `protobuf:"varint,6,opt,name=failedCount" json:"failedCount,omitempty"`
	RunningTasks       []*Task `protobuf:"bytes,7,rep,name=runningTasks" json:"runningTasks,omitempty"`
	ExpiredTasks       []*Task `protobuf:"bytes,8,rep,name=expiredTasks" json:"expiredTasks,omitempty"`
}

func (m *QueueProgress) Reset()                    { *m = QueueProgress{} }
func (m *QueueProgress) String() string            { return proto.CompactTextString(m) }
func (*QueueProgress) ProtoMessage()               {}
func (*QueueProgress) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *QueueProgress) GetIncomingListLength() int32 {
	if m != nil {
		return m.IncomingListLength
	}
	return 0
}

func (m *QueueProgress) GetReceivedListLength() int32 {
	if m != nil {
		return m.ReceivedListLength
	}
	return 0
}

func (m *QueueProgress) GetNonExpiredCount() int32 {
	if m != nil {
		return m.NonExpiredCount
	}
	return 0
}

func (m *QueueProgress) GetExpiredCount() int32 {
	if m != nil {
		return m.ExpiredCount
	}
	return 0
}

func (m *QueueProgress) GetCompletedCount() int32 {
	if m != nil {
		return m.CompletedCount
	}
	return 0
}

func (m *QueueProgress) GetFailedCount() int32 {
	if m != nil {
		return m.FailedCount
	}
	return 0
}

func (m *QueueProgress) GetRunningTasks() []*Task {
	if m != nil {
		return m.RunningTasks
	}
	return nil
}

func (m *QueueProgress) GetExpiredTasks() []*Task {
	if m != nil {
		return m.ExpiredTasks
	}
	return nil
}

func init() {
	proto.RegisterType((*Task)(nil), "API.Task")
	proto.RegisterType((*Meta)(nil), "API.Meta")
	proto.RegisterType((*RequestMessage)(nil), "API.RequestMessage")
	proto.RegisterType((*StatusMessage)(nil), "API.StatusMessage")
	proto.RegisterType((*QueueProgress)(nil), "API.QueueProgress")
	proto.RegisterEnum("API.Status", Status_name, Status_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for API service

type APIClient interface {
	Healthz(ctx context.Context, in *google_protobuf1.Empty, opts ...grpc.CallOption) (*StatusMessage, error)
	PushTask(ctx context.Context, in *Task, opts ...grpc.CallOption) (*StatusMessage, error)
	LoadTask(ctx context.Context, in *RequestMessage, opts ...grpc.CallOption) (*Task, error)
	HeartBeat(ctx context.Context, in *Task, opts ...grpc.CallOption) (*StatusMessage, error)
	Complete(ctx context.Context, in *Task, opts ...grpc.CallOption) (*StatusMessage, error)
	Progress(ctx context.Context, in *RequestMessage, opts ...grpc.CallOption) (*QueueProgress, error)
}

type aPIClient struct {
	cc *grpc.ClientConn
}

func NewAPIClient(cc *grpc.ClientConn) APIClient {
	return &aPIClient{cc}
}

func (c *aPIClient) Healthz(ctx context.Context, in *google_protobuf1.Empty, opts ...grpc.CallOption) (*StatusMessage, error) {
	out := new(StatusMessage)
	err := grpc.Invoke(ctx, "/API.API/Healthz", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aPIClient) PushTask(ctx context.Context, in *Task, opts ...grpc.CallOption) (*StatusMessage, error) {
	out := new(StatusMessage)
	err := grpc.Invoke(ctx, "/API.API/PushTask", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aPIClient) LoadTask(ctx context.Context, in *RequestMessage, opts ...grpc.CallOption) (*Task, error) {
	out := new(Task)
	err := grpc.Invoke(ctx, "/API.API/LoadTask", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aPIClient) HeartBeat(ctx context.Context, in *Task, opts ...grpc.CallOption) (*StatusMessage, error) {
	out := new(StatusMessage)
	err := grpc.Invoke(ctx, "/API.API/HeartBeat", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aPIClient) Complete(ctx context.Context, in *Task, opts ...grpc.CallOption) (*StatusMessage, error) {
	out := new(StatusMessage)
	err := grpc.Invoke(ctx, "/API.API/Complete", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aPIClient) Progress(ctx context.Context, in *RequestMessage, opts ...grpc.CallOption) (*QueueProgress, error) {
	out := new(QueueProgress)
	err := grpc.Invoke(ctx, "/API.API/Progress", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for API service

type APIServer interface {
	Healthz(context.Context, *google_protobuf1.Empty) (*StatusMessage, error)
	PushTask(context.Context, *Task) (*StatusMessage, error)
	LoadTask(context.Context, *RequestMessage) (*Task, error)
	HeartBeat(context.Context, *Task) (*StatusMessage, error)
	Complete(context.Context, *Task) (*StatusMessage, error)
	Progress(context.Context, *RequestMessage) (*QueueProgress, error)
}

func RegisterAPIServer(s *grpc.Server, srv APIServer) {
	s.RegisterService(&_API_serviceDesc, srv)
}

func _API_Healthz_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(google_protobuf1.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APIServer).Healthz(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/API.API/Healthz",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APIServer).Healthz(ctx, req.(*google_protobuf1.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _API_PushTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Task)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APIServer).PushTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/API.API/PushTask",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APIServer).PushTask(ctx, req.(*Task))
	}
	return interceptor(ctx, in, info, handler)
}

func _API_LoadTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APIServer).LoadTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/API.API/LoadTask",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APIServer).LoadTask(ctx, req.(*RequestMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _API_HeartBeat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Task)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APIServer).HeartBeat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/API.API/HeartBeat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APIServer).HeartBeat(ctx, req.(*Task))
	}
	return interceptor(ctx, in, info, handler)
}

func _API_Complete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Task)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APIServer).Complete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/API.API/Complete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APIServer).Complete(ctx, req.(*Task))
	}
	return interceptor(ctx, in, info, handler)
}

func _API_Progress_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APIServer).Progress(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/API.API/Progress",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APIServer).Progress(ctx, req.(*RequestMessage))
	}
	return interceptor(ctx, in, info, handler)
}

var _API_serviceDesc = grpc.ServiceDesc{
	ServiceName: "API.API",
	HandlerType: (*APIServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Healthz",
			Handler:    _API_Healthz_Handler,
		},
		{
			MethodName: "PushTask",
			Handler:    _API_PushTask_Handler,
		},
		{
			MethodName: "LoadTask",
			Handler:    _API_LoadTask_Handler,
		},
		{
			MethodName: "HeartBeat",
			Handler:    _API_HeartBeat_Handler,
		},
		{
			MethodName: "Complete",
			Handler:    _API_Complete_Handler,
		},
		{
			MethodName: "Progress",
			Handler:    _API_Progress_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api.proto",
}

func init() { proto.RegisterFile("api.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 599 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x53, 0xcf, 0x4e, 0xdb, 0x4e,
	0x10, 0xfe, 0xe5, 0x7f, 0x32, 0x01, 0x7e, 0xd1, 0xa2, 0x22, 0x2b, 0x20, 0x64, 0xb9, 0xa8, 0x8a,
	0x90, 0x6a, 0xab, 0xf4, 0xd6, 0x1b, 0x0d, 0x29, 0x58, 0x0d, 0xe0, 0x3a, 0x45, 0x3d, 0x2f, 0xc9,
	0xe0, 0x58, 0x84, 0x5d, 0xe3, 0x5d, 0xa3, 0x52, 0xc4, 0xa5, 0x87, 0xbe, 0x40, 0x1f, 0xad, 0x87,
	0xbe, 0x40, 0xfb, 0x1e, 0xd5, 0xae, 0xed, 0x92, 0x80, 0x4b, 0x2f, 0xbd, 0xed, 0x7c, 0xf3, 0xed,
	0xb7, 0x33, 0x3b, 0xf3, 0x41, 0x8b, 0x46, 0xa1, 0x1d, 0xc5, 0x5c, 0x72, 0x52, 0xd9, 0xf5, 0xdc,
	0xee, 0x46, 0xc0, 0x79, 0x30, 0x43, 0x87, 0x46, 0xa1, 0x43, 0x19, 0xe3, 0x92, 0xca, 0x90, 0x33,
	0x91, 0x52, 0xba, 0xeb, 0x59, 0x56, 0x47, 0xa7, 0xc9, 0x99, 0x83, 0x17, 0x91, 0xbc, 0x4e, 0x93,
	0xd6, 0x97, 0x12, 0x54, 0xdf, 0x53, 0x71, 0x4e, 0x0c, 0x68, 0x5c, 0x26, 0x98, 0xa0, 0xbb, 0x67,
	0x94, 0xcc, 0x52, 0xaf, 0xe5, 0xe7, 0x21, 0x59, 0x83, 0xba, 0xa4, 0xe2, 0xdc, 0xdd, 0x33, 0xca,
	0x3a, 0x91, 0x45, 0x84, 0x40, 0xf5, 0x94, 0x4f, 0xae, 0x8d, 0x8a, 0x46, 0xf5, 0x59, 0x61, 0xc8,
	0xae, 0x84, 0x51, 0x35, 0x2b, 0x0a, 0x53, 0x67, 0xb2, 0x05, 0xcb, 0xf8, 0x31, 0x0a, 0x63, 0x5d,
	0xd4, 0x08, 0xc7, 0x46, 0xcd, 0x2c, 0xf5, 0x6a, 0xfe, 0x22, 0x68, 0x99, 0x50, 0x3d, 0x44, 0x49,
	0xff, 0x5c, 0x87, 0xe5, 0xc1, 0x8a, 0x8f, 0x97, 0x09, 0x0a, 0x79, 0x88, 0x42, 0xd0, 0x00, 0x1f,
	0xa9, 0xf9, 0xc1, 0x9b, 0xe5, 0xa2, 0x37, 0x87, 0xb0, 0x3c, 0x92, 0x54, 0x26, 0x22, 0x17, 0x7c,
	0x0a, 0x75, 0xa1, 0x01, 0xad, 0xb7, 0xb2, 0xd3, 0xb6, 0x77, 0x3d, 0xd7, 0x4e, 0x39, 0x7e, 0x96,
	0x52, 0xff, 0x31, 0x41, 0x49, 0xc3, 0x59, 0xfe, 0x1f, 0x69, 0x64, 0xfd, 0x2c, 0xc3, 0xf2, 0x3b,
	0xf5, 0xbe, 0x17, 0xf3, 0x20, 0x46, 0x21, 0x88, 0x0d, 0x24, 0x64, 0x63, 0x7e, 0x11, 0xb2, 0x60,
	0x18, 0x0a, 0x39, 0x44, 0x16, 0xc8, 0xa9, 0x96, 0xae, 0xf9, 0x05, 0x19, 0xc5, 0x8f, 0x71, 0x8c,
	0xe1, 0x15, 0x4e, 0xe6, 0xf8, 0x69, 0xe9, 0x05, 0x19, 0xd2, 0x83, 0xff, 0x19, 0x67, 0x03, 0xd5,
	0x13, 0x4e, 0xfa, 0x3c, 0x61, 0x52, 0x0f, 0xa3, 0xe6, 0xdf, 0x87, 0x89, 0x05, 0x4b, 0x38, 0x4f,
	0xab, 0x6a, 0xda, 0x02, 0x46, 0x9e, 0xc1, 0xca, 0x98, 0x5f, 0x44, 0x33, 0x94, 0x39, 0x2b, 0x1d,
	0xd4, 0x3d, 0x94, 0x98, 0xd0, 0x3e, 0xa3, 0xe1, 0x2c, 0x27, 0xd5, 0x35, 0x69, 0x1e, 0x22, 0xcf,
	0x61, 0x29, 0x4e, 0x18, 0x0b, 0x59, 0xa0, 0x56, 0x4b, 0x18, 0x0d, 0xb3, 0xd2, 0x6b, 0xef, 0xb4,
	0xf4, 0x67, 0x2a, 0xc4, 0x5f, 0x48, 0x2b, 0x7a, 0x56, 0x48, 0x4a, 0x6f, 0x3e, 0xa0, 0xcf, 0xa7,
	0xb7, 0x1d, 0xa8, 0xa7, 0x13, 0x21, 0x6d, 0x68, 0x8c, 0x4e, 0xfa, 0xfd, 0xc1, 0x68, 0xd4, 0xf9,
	0x4f, 0x05, 0x6f, 0x76, 0xdd, 0xe1, 0x89, 0x3f, 0xe8, 0x94, 0x54, 0x70, 0x72, 0xf4, 0xf6, 0xe8,
	0xf8, 0xc3, 0x51, 0xa7, 0xbc, 0xf3, 0xbd, 0x02, 0xca, 0x26, 0x64, 0x1f, 0x1a, 0x07, 0x48, 0x67,
	0x72, 0xfa, 0x89, 0xac, 0xd9, 0xa9, 0x29, 0xec, 0xdc, 0x14, 0xf6, 0x40, 0x99, 0xa2, 0x4b, 0xe6,
	0x06, 0x9e, 0x2d, 0x85, 0xd5, 0xf9, 0xfc, 0xed, 0xc7, 0xd7, 0x32, 0x90, 0xa6, 0x33, 0xcd, 0x6e,
	0xbb, 0xd0, 0xf4, 0x12, 0x31, 0xd5, 0xbe, 0xb9, 0x2b, 0xb3, 0xf0, 0xf2, 0xa6, 0xbe, 0x6c, 0x58,
	0xab, 0xce, 0xd5, 0x0b, 0x47, 0x6f, 0xa7, 0x73, 0x93, 0x2d, 0xe9, 0xed, 0xab, 0xd2, 0x36, 0xd9,
	0x87, 0xe6, 0x90, 0x53, 0xdd, 0x19, 0x59, 0xd5, 0xf7, 0x17, 0x77, 0xbc, 0x7b, 0xa7, 0x6f, 0xad,
	0x6b, 0xad, 0x27, 0xa4, 0x48, 0x8b, 0x78, 0xd0, 0x3a, 0x40, 0x1a, 0xcb, 0xd7, 0x48, 0xe5, 0xdf,
	0x8a, 0xda, 0xd2, 0x42, 0x9b, 0xdd, 0x8d, 0x02, 0x21, 0xe7, 0x26, 0xb5, 0xf7, 0x2d, 0x39, 0x86,
	0x66, 0x3f, 0x9b, 0xfc, 0xbf, 0x11, 0xf4, 0xa0, 0xf9, 0xdb, 0x1a, 0x85, 0xbd, 0xa6, 0xd2, 0x0b,
	0x1e, 0x7a, 0xb4, 0xe9, 0xd3, 0xba, 0x1e, 0xdf, 0xcb, 0x5f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x3b,
	0xbb, 0xa0, 0x87, 0x11, 0x05, 0x00, 0x00,
}
