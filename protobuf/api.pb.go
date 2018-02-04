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
	TaskList
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
	Phase         string `protobuf:"bytes,3,opt,name=phase" json:"phase,omitempty"`
	Limit         int32  `protobuf:"varint,4,opt,name=limit" json:"limit,omitempty"`
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

func (m *RequestMessage) GetPhase() string {
	if m != nil {
		return m.Phase
	}
	return ""
}

func (m *RequestMessage) GetLimit() int32 {
	if m != nil {
		return m.Limit
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
	IncomingCount  int32 `protobuf:"varint,1,opt,name=incomingCount" json:"incomingCount,omitempty"`
	ReceivedCount  int32 `protobuf:"varint,2,opt,name=receivedCount" json:"receivedCount,omitempty"`
	RunningCount   int32 `protobuf:"varint,3,opt,name=runningCount" json:"runningCount,omitempty"`
	ExpiredCount   int32 `protobuf:"varint,4,opt,name=expiredCount" json:"expiredCount,omitempty"`
	CompletedCount int32 `protobuf:"varint,5,opt,name=completedCount" json:"completedCount,omitempty"`
	FailedCount    int32 `protobuf:"varint,6,opt,name=failedCount" json:"failedCount,omitempty"`
}

func (m *QueueProgress) Reset()                    { *m = QueueProgress{} }
func (m *QueueProgress) String() string            { return proto.CompactTextString(m) }
func (*QueueProgress) ProtoMessage()               {}
func (*QueueProgress) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *QueueProgress) GetIncomingCount() int32 {
	if m != nil {
		return m.IncomingCount
	}
	return 0
}

func (m *QueueProgress) GetReceivedCount() int32 {
	if m != nil {
		return m.ReceivedCount
	}
	return 0
}

func (m *QueueProgress) GetRunningCount() int32 {
	if m != nil {
		return m.RunningCount
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

type TaskList struct {
	TotalItems int32   `protobuf:"varint,1,opt,name=totalItems" json:"totalItems,omitempty"`
	Tasks      []*Task `protobuf:"bytes,2,rep,name=tasks" json:"tasks,omitempty"`
}

func (m *TaskList) Reset()                    { *m = TaskList{} }
func (m *TaskList) String() string            { return proto.CompactTextString(m) }
func (*TaskList) ProtoMessage()               {}
func (*TaskList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *TaskList) GetTotalItems() int32 {
	if m != nil {
		return m.TotalItems
	}
	return 0
}

func (m *TaskList) GetTasks() []*Task {
	if m != nil {
		return m.Tasks
	}
	return nil
}

func init() {
	proto.RegisterType((*Task)(nil), "API.Task")
	proto.RegisterType((*Meta)(nil), "API.Meta")
	proto.RegisterType((*RequestMessage)(nil), "API.RequestMessage")
	proto.RegisterType((*StatusMessage)(nil), "API.StatusMessage")
	proto.RegisterType((*QueueProgress)(nil), "API.QueueProgress")
	proto.RegisterType((*TaskList)(nil), "API.TaskList")
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
	Peek(ctx context.Context, in *RequestMessage, opts ...grpc.CallOption) (*TaskList, error)
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

func (c *aPIClient) Peek(ctx context.Context, in *RequestMessage, opts ...grpc.CallOption) (*TaskList, error) {
	out := new(TaskList)
	err := grpc.Invoke(ctx, "/API.API/Peek", in, out, c.cc, opts...)
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
	Peek(context.Context, *RequestMessage) (*TaskList, error)
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

func _API_Peek_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APIServer).Peek(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/API.API/Peek",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APIServer).Peek(ctx, req.(*RequestMessage))
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
		{
			MethodName: "Peek",
			Handler:    _API_Peek_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api.proto",
}

func init() { proto.RegisterFile("api.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 645 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x53, 0xc1, 0x4e, 0xdb, 0x4c,
	0x10, 0xc6, 0x89, 0x13, 0x92, 0xc9, 0x1f, 0x14, 0x2d, 0xbf, 0x90, 0x95, 0x22, 0x1a, 0x19, 0x5a,
	0x21, 0x5a, 0xc5, 0x2a, 0xbd, 0xf5, 0x46, 0x03, 0x2d, 0x16, 0x01, 0x52, 0xa7, 0x08, 0xf5, 0xb8,
	0x24, 0x43, 0x62, 0xe1, 0x78, 0x8d, 0x77, 0x8d, 0x4a, 0x69, 0x2e, 0x3d, 0xf4, 0x05, 0x7a, 0xef,
	0x4b, 0xf5, 0x15, 0x7a, 0xed, 0x3b, 0x54, 0xbb, 0x6b, 0x43, 0xdc, 0x86, 0xf6, 0xd2, 0xdb, 0xce,
	0xb7, 0xdf, 0x7c, 0x33, 0x3b, 0xfb, 0x0d, 0x54, 0x69, 0xe4, 0xb7, 0xa3, 0x98, 0x09, 0x46, 0x8a,
	0x3b, 0x3d, 0xb7, 0xb9, 0x3a, 0x62, 0x6c, 0x14, 0xa0, 0x43, 0x23, 0xdf, 0xa1, 0x61, 0xc8, 0x04,
	0x15, 0x3e, 0x0b, 0xb9, 0xa6, 0x34, 0x1f, 0xa4, 0xb7, 0x2a, 0x3a, 0x4b, 0xce, 0x1d, 0x9c, 0x44,
	0xe2, 0x5a, 0x5f, 0xda, 0x9f, 0x0d, 0x30, 0xdf, 0x52, 0x7e, 0x41, 0x2c, 0x58, 0xbc, 0x4c, 0x30,
	0x41, 0x77, 0xd7, 0x32, 0x5a, 0xc6, 0x66, 0xd5, 0xcb, 0x42, 0xb2, 0x02, 0x65, 0x41, 0xf9, 0x85,
	0xbb, 0x6b, 0x15, 0xd4, 0x45, 0x1a, 0x11, 0x02, 0xe6, 0x19, 0x1b, 0x5e, 0x5b, 0x45, 0x85, 0xaa,
	0xb3, 0xc4, 0x30, 0xbc, 0xe2, 0x96, 0xd9, 0x2a, 0x4a, 0x4c, 0x9e, 0xc9, 0x06, 0xd4, 0xf1, 0x7d,
	0xe4, 0xc7, 0xaa, 0xa9, 0x3e, 0x0e, 0xac, 0x52, 0xcb, 0xd8, 0x2c, 0x79, 0x79, 0xd0, 0x6e, 0x81,
	0x79, 0x88, 0x82, 0xde, 0xdf, 0x87, 0xfd, 0x11, 0x96, 0x3c, 0xbc, 0x4c, 0x90, 0x8b, 0x43, 0xe4,
	0x9c, 0x8e, 0xf0, 0x0f, 0x3d, 0xff, 0x56, 0xb3, 0x30, 0xa7, 0x26, 0xf9, 0x1f, 0x4a, 0xd1, 0x98,
	0x72, 0x4c, 0x9f, 0xa0, 0x03, 0x89, 0x06, 0xfe, 0xc4, 0x17, 0x96, 0xa9, 0x72, 0x74, 0x60, 0x77,
	0xa1, 0xde, 0x17, 0x54, 0x24, 0x3c, 0x2b, 0xbe, 0x0e, 0x65, 0xae, 0x00, 0x55, 0x7b, 0x69, 0xbb,
	0xd6, 0xde, 0xe9, 0xb9, 0x6d, 0xcd, 0xf1, 0xd2, 0x2b, 0x39, 0xbb, 0x21, 0x0a, 0xea, 0x07, 0xd9,
	0xec, 0x74, 0x64, 0xff, 0x30, 0xa0, 0xfe, 0x46, 0xf6, 0xda, 0x8b, 0xd9, 0x28, 0x46, 0xae, 0xa6,
	0xe4, 0x87, 0x03, 0x36, 0xf1, 0xc3, 0x51, 0x87, 0x25, 0xa1, 0x50, 0xaa, 0x25, 0x2f, 0x0f, 0x4a,
	0x56, 0x8c, 0x03, 0xf4, 0xaf, 0x70, 0xa8, 0x59, 0xe9, 0xbb, 0x72, 0x20, 0xb1, 0xe1, 0xbf, 0x38,
	0x09, 0xc3, 0x5b, 0xa9, 0xa2, 0x22, 0xe5, 0x30, 0xc9, 0x51, 0xc3, 0xc8, 0x84, 0xf4, 0x63, 0x73,
	0x18, 0x79, 0x0c, 0x4b, 0x03, 0x36, 0x89, 0x02, 0x14, 0x19, 0x4b, 0x7f, 0xdd, 0x2f, 0x28, 0x69,
	0x41, 0xed, 0x9c, 0xfa, 0x41, 0x46, 0x2a, 0x2b, 0xd2, 0x2c, 0x64, 0x1f, 0x40, 0x45, 0xba, 0xac,
	0xeb, 0x73, 0x41, 0xd6, 0x00, 0x04, 0x13, 0x34, 0x70, 0x05, 0x4e, 0x78, 0xfa, 0xcc, 0x19, 0x84,
	0x3c, 0x84, 0x92, 0x74, 0x18, 0xb7, 0x0a, 0xad, 0xe2, 0x66, 0x6d, 0xbb, 0xaa, 0xe6, 0x2a, 0xb3,
	0x3d, 0x8d, 0x6f, 0x39, 0x50, 0xd6, 0x63, 0x26, 0x35, 0x58, 0xec, 0x9f, 0x74, 0x3a, 0x7b, 0xfd,
	0x7e, 0x63, 0x41, 0x06, 0xaf, 0x76, 0xdc, 0xee, 0x89, 0xb7, 0xd7, 0x30, 0x64, 0x70, 0x72, 0x74,
	0x70, 0x74, 0x7c, 0x7a, 0xd4, 0x28, 0x6c, 0x7f, 0x35, 0x41, 0xee, 0x09, 0x79, 0x0d, 0x8b, 0xfb,
	0x48, 0x03, 0x31, 0xfe, 0x40, 0x56, 0xda, 0x7a, 0x2b, 0xda, 0xd9, 0x56, 0xb4, 0xf7, 0xe4, 0x56,
	0x34, 0xc9, 0xcc, 0x2f, 0xa6, 0x3f, 0x6d, 0x37, 0x3e, 0x7d, 0xfb, 0xfe, 0xa5, 0x00, 0xa4, 0xe2,
	0x8c, 0xd3, 0x6c, 0x17, 0x2a, 0xbd, 0x84, 0x8f, 0xd5, 0xe2, 0xdc, 0xf5, 0x37, 0x37, 0x79, 0x4d,
	0x25, 0x5b, 0xf6, 0xb2, 0x73, 0xf5, 0xcc, 0x51, 0xf6, 0x74, 0x6e, 0x52, 0x97, 0x4e, 0x5f, 0x18,
	0x5b, 0xe4, 0x29, 0x54, 0xba, 0x8c, 0x0e, 0x95, 0xd4, 0xb2, 0xca, 0xcf, 0x9b, 0xbc, 0x79, 0xa7,
	0x6f, 0x2f, 0x90, 0x1e, 0x54, 0xf7, 0x91, 0xc6, 0xe2, 0x25, 0x52, 0xf1, 0xb7, 0xca, 0x1b, 0xaa,
	0xf2, 0x5a, 0x73, 0x75, 0x4e, 0x65, 0xe7, 0x46, 0x2f, 0xf1, 0x94, 0x1c, 0x43, 0xa5, 0x93, 0xfe,
	0xe6, 0xbf, 0x11, 0x7c, 0x07, 0x95, 0x5b, 0x53, 0xcf, 0x7d, 0x90, 0x96, 0xce, 0xb9, 0x3f, 0x93,
	0x26, 0x73, 0xa5, 0xa3, 0x4c, 0xee, 0x14, 0xcc, 0x1e, 0xe2, 0x3d, 0x73, 0xaa, 0xdf, 0x36, 0x2f,
	0x5d, 0x66, 0x3f, 0x51, 0x8a, 0x8f, 0xc8, 0xfa, 0xdc, 0x66, 0xd5, 0xa6, 0x4f, 0x9d, 0x1b, 0xb5,
	0xdb, 0xd3, 0xb3, 0xb2, 0x72, 0xc1, 0xf3, 0x9f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x59, 0xb0, 0xca,
	0x38, 0x59, 0x05, 0x00, 0x00,
}
