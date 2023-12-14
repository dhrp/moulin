// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.25.1
// source: api.proto

package API

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Status int32

const (
	Status_SUCCESS Status = 0
	Status_FAILURE Status = 1
	Status_UNKNOWN Status = 2
)

// Enum value maps for Status.
var (
	Status_name = map[int32]string{
		0: "SUCCESS",
		1: "FAILURE",
		2: "UNKNOWN",
	}
	Status_value = map[string]int32{
		"SUCCESS": 0,
		"FAILURE": 1,
		"UNKNOWN": 2,
	}
)

func (x Status) Enum() *Status {
	p := new(Status)
	*p = x
	return p
}

func (x Status) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Status) Descriptor() protoreflect.EnumDescriptor {
	return file_api_proto_enumTypes[0].Descriptor()
}

func (Status) Type() protoreflect.EnumType {
	return &file_api_proto_enumTypes[0]
}

func (x Status) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Status.Descriptor instead.
func (Status) EnumDescriptor() ([]byte, []int) {
	return file_api_proto_rawDescGZIP(), []int{0}
}

// Task is the definition of a task
type Task struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// reserved 1; // we took queueID out
	// reserved "queueID";
	QueueID       string   `protobuf:"bytes,1,opt,name=queueID,proto3" json:"queueID,omitempty"` // needed to pass the queueID to push task to
	TaskID        string   `protobuf:"bytes,2,opt,name=taskID,proto3" json:"taskID,omitempty"`
	Body          string   `protobuf:"bytes,3,opt,name=body,proto3" json:"body,omitempty"` // perhaps should be called argument(s)
	Envs          []string `protobuf:"bytes,4,rep,name=envs,proto3" json:"envs,omitempty"`
	ExpirationSec int32    `protobuf:"varint,5,opt,name=expirationSec,proto3" json:"expirationSec,omitempty"` // this is used to heartbeat
}

func (x *Task) Reset() {
	*x = Task{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Task) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Task) ProtoMessage() {}

func (x *Task) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Task.ProtoReflect.Descriptor instead.
func (*Task) Descriptor() ([]byte, []int) {
	return file_api_proto_rawDescGZIP(), []int{0}
}

func (x *Task) GetQueueID() string {
	if x != nil {
		return x.QueueID
	}
	return ""
}

func (x *Task) GetTaskID() string {
	if x != nil {
		return x.TaskID
	}
	return ""
}

func (x *Task) GetBody() string {
	if x != nil {
		return x.Body
	}
	return ""
}

func (x *Task) GetEnvs() []string {
	if x != nil {
		return x.Envs
	}
	return nil
}

func (x *Task) GetExpirationSec() int32 {
	if x != nil {
		return x.ExpirationSec
	}
	return 0
}

type Meta struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	QueueID string `protobuf:"bytes,1,opt,name=queueID,proto3" json:"queueID,omitempty"`
}

func (x *Meta) Reset() {
	*x = Meta{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Meta) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Meta) ProtoMessage() {}

func (x *Meta) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Meta.ProtoReflect.Descriptor instead.
func (*Meta) Descriptor() ([]byte, []int) {
	return file_api_proto_rawDescGZIP(), []int{1}
}

func (x *Meta) GetQueueID() string {
	if x != nil {
		return x.QueueID
	}
	return ""
}

type RequestMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	QueueID       string `protobuf:"bytes,1,opt,name=queueID,proto3" json:"queueID,omitempty"`
	ExpirationSec int32  `protobuf:"varint,2,opt,name=expirationSec,proto3" json:"expirationSec,omitempty"`
	Phase         string `protobuf:"bytes,3,opt,name=phase,proto3" json:"phase,omitempty"`  // only valid for peek
	Limit         int32  `protobuf:"varint,4,opt,name=limit,proto3" json:"limit,omitempty"` // only valid for peek
}

func (x *RequestMessage) Reset() {
	*x = RequestMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestMessage) ProtoMessage() {}

func (x *RequestMessage) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestMessage.ProtoReflect.Descriptor instead.
func (*RequestMessage) Descriptor() ([]byte, []int) {
	return file_api_proto_rawDescGZIP(), []int{2}
}

func (x *RequestMessage) GetQueueID() string {
	if x != nil {
		return x.QueueID
	}
	return ""
}

func (x *RequestMessage) GetExpirationSec() int32 {
	if x != nil {
		return x.ExpirationSec
	}
	return 0
}

func (x *RequestMessage) GetPhase() string {
	if x != nil {
		return x.Phase
	}
	return ""
}

func (x *RequestMessage) GetLimit() int32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

type StatusMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status Status `protobuf:"varint,1,opt,name=status,proto3,enum=API.Status" json:"status,omitempty"`
	Detail string `protobuf:"bytes,2,opt,name=detail,proto3" json:"detail,omitempty"`
}

func (x *StatusMessage) Reset() {
	*x = StatusMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StatusMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StatusMessage) ProtoMessage() {}

func (x *StatusMessage) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StatusMessage.ProtoReflect.Descriptor instead.
func (*StatusMessage) Descriptor() ([]byte, []int) {
	return file_api_proto_rawDescGZIP(), []int{3}
}

func (x *StatusMessage) GetStatus() Status {
	if x != nil {
		return x.Status
	}
	return Status_SUCCESS
}

func (x *StatusMessage) GetDetail() string {
	if x != nil {
		return x.Detail
	}
	return ""
}

type QueueProgress struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IncomingCount  int32 `protobuf:"varint,1,opt,name=incomingCount,proto3" json:"incomingCount,omitempty"`
	ReceivedCount  int32 `protobuf:"varint,2,opt,name=receivedCount,proto3" json:"receivedCount,omitempty"`
	RunningCount   int32 `protobuf:"varint,3,opt,name=runningCount,proto3" json:"runningCount,omitempty"`
	ExpiredCount   int32 `protobuf:"varint,4,opt,name=expiredCount,proto3" json:"expiredCount,omitempty"`
	CompletedCount int32 `protobuf:"varint,5,opt,name=completedCount,proto3" json:"completedCount,omitempty"`
	FailedCount    int32 `protobuf:"varint,6,opt,name=failedCount,proto3" json:"failedCount,omitempty"`
}

func (x *QueueProgress) Reset() {
	*x = QueueProgress{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueueProgress) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueueProgress) ProtoMessage() {}

func (x *QueueProgress) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueueProgress.ProtoReflect.Descriptor instead.
func (*QueueProgress) Descriptor() ([]byte, []int) {
	return file_api_proto_rawDescGZIP(), []int{4}
}

func (x *QueueProgress) GetIncomingCount() int32 {
	if x != nil {
		return x.IncomingCount
	}
	return 0
}

func (x *QueueProgress) GetReceivedCount() int32 {
	if x != nil {
		return x.ReceivedCount
	}
	return 0
}

func (x *QueueProgress) GetRunningCount() int32 {
	if x != nil {
		return x.RunningCount
	}
	return 0
}

func (x *QueueProgress) GetExpiredCount() int32 {
	if x != nil {
		return x.ExpiredCount
	}
	return 0
}

func (x *QueueProgress) GetCompletedCount() int32 {
	if x != nil {
		return x.CompletedCount
	}
	return 0
}

func (x *QueueProgress) GetFailedCount() int32 {
	if x != nil {
		return x.FailedCount
	}
	return 0
}

type TaskList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TotalItems int32   `protobuf:"varint,1,opt,name=totalItems,proto3" json:"totalItems,omitempty"`
	Tasks      []*Task `protobuf:"bytes,2,rep,name=tasks,proto3" json:"tasks,omitempty"`
}

func (x *TaskList) Reset() {
	*x = TaskList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TaskList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TaskList) ProtoMessage() {}

func (x *TaskList) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TaskList.ProtoReflect.Descriptor instead.
func (*TaskList) Descriptor() ([]byte, []int) {
	return file_api_proto_rawDescGZIP(), []int{5}
}

func (x *TaskList) GetTotalItems() int32 {
	if x != nil {
		return x.TotalItems
	}
	return 0
}

func (x *TaskList) GetTasks() []*Task {
	if x != nil {
		return x.Tasks
	}
	return nil
}

type QueueMap struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TotalItems int32                     `protobuf:"varint,1,opt,name=totalItems,proto3" json:"totalItems,omitempty"`
	Queues     map[string]*QueueProgress `protobuf:"bytes,2,rep,name=queues,proto3" json:"queues,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *QueueMap) Reset() {
	*x = QueueMap{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueueMap) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueueMap) ProtoMessage() {}

func (x *QueueMap) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueueMap.ProtoReflect.Descriptor instead.
func (*QueueMap) Descriptor() ([]byte, []int) {
	return file_api_proto_rawDescGZIP(), []int{6}
}

func (x *QueueMap) GetTotalItems() int32 {
	if x != nil {
		return x.TotalItems
	}
	return 0
}

func (x *QueueMap) GetQueues() map[string]*QueueProgress {
	if x != nil {
		return x.Queues
	}
	return nil
}

var File_api_proto protoreflect.FileDescriptor

var file_api_proto_rawDesc = []byte{
	0x0a, 0x09, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x41, 0x50, 0x49,
	0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e,
	0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x86, 0x01, 0x0a, 0x04,
	0x54, 0x61, 0x73, 0x6b, 0x12, 0x18, 0x0a, 0x07, 0x71, 0x75, 0x65, 0x75, 0x65, 0x49, 0x44, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x71, 0x75, 0x65, 0x75, 0x65, 0x49, 0x44, 0x12, 0x16,
	0x0a, 0x06, 0x74, 0x61, 0x73, 0x6b, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x74, 0x61, 0x73, 0x6b, 0x49, 0x44, 0x12, 0x12, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x65, 0x6e,
	0x76, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x65, 0x6e, 0x76, 0x73, 0x12, 0x24,
	0x0a, 0x0d, 0x65, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x63, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0d, 0x65, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x53, 0x65, 0x63, 0x22, 0x20, 0x0a, 0x04, 0x4d, 0x65, 0x74, 0x61, 0x12, 0x18, 0x0a, 0x07,
	0x71, 0x75, 0x65, 0x75, 0x65, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x71,
	0x75, 0x65, 0x75, 0x65, 0x49, 0x44, 0x22, 0x7c, 0x0a, 0x0e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x71, 0x75, 0x65, 0x75,
	0x65, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x71, 0x75, 0x65, 0x75, 0x65,
	0x49, 0x44, 0x12, 0x24, 0x0a, 0x0d, 0x65, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x53, 0x65, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0d, 0x65, 0x78, 0x70, 0x69, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x63, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x68, 0x61, 0x73,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x70, 0x68, 0x61, 0x73, 0x65, 0x12, 0x14,
	0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x6c,
	0x69, 0x6d, 0x69, 0x74, 0x22, 0x4c, 0x0a, 0x0d, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x23, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0b, 0x2e, 0x41, 0x50, 0x49, 0x2e, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x64, 0x65,
	0x74, 0x61, 0x69, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x64, 0x65, 0x74, 0x61,
	0x69, 0x6c, 0x22, 0xed, 0x01, 0x0a, 0x0d, 0x51, 0x75, 0x65, 0x75, 0x65, 0x50, 0x72, 0x6f, 0x67,
	0x72, 0x65, 0x73, 0x73, 0x12, 0x24, 0x0a, 0x0d, 0x69, 0x6e, 0x63, 0x6f, 0x6d, 0x69, 0x6e, 0x67,
	0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0d, 0x69, 0x6e, 0x63,
	0x6f, 0x6d, 0x69, 0x6e, 0x67, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x24, 0x0a, 0x0d, 0x72, 0x65,
	0x63, 0x65, 0x69, 0x76, 0x65, 0x64, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x0d, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x64, 0x43, 0x6f, 0x75, 0x6e, 0x74,
	0x12, 0x22, 0x0a, 0x0c, 0x72, 0x75, 0x6e, 0x6e, 0x69, 0x6e, 0x67, 0x43, 0x6f, 0x75, 0x6e, 0x74,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0c, 0x72, 0x75, 0x6e, 0x6e, 0x69, 0x6e, 0x67, 0x43,
	0x6f, 0x75, 0x6e, 0x74, 0x12, 0x22, 0x0a, 0x0c, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x64, 0x43,
	0x6f, 0x75, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0c, 0x65, 0x78, 0x70, 0x69,
	0x72, 0x65, 0x64, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x26, 0x0a, 0x0e, 0x63, 0x6f, 0x6d, 0x70,
	0x6c, 0x65, 0x74, 0x65, 0x64, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x0e, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x43, 0x6f, 0x75, 0x6e, 0x74,
	0x12, 0x20, 0x0a, 0x0b, 0x66, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x66, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x43, 0x6f, 0x75,
	0x6e, 0x74, 0x22, 0x4b, 0x0a, 0x08, 0x54, 0x61, 0x73, 0x6b, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x1e,
	0x0a, 0x0a, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x0a, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x12, 0x1f,
	0x0a, 0x05, 0x74, 0x61, 0x73, 0x6b, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x09, 0x2e,
	0x41, 0x50, 0x49, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x05, 0x74, 0x61, 0x73, 0x6b, 0x73, 0x22,
	0xac, 0x01, 0x0a, 0x08, 0x51, 0x75, 0x65, 0x75, 0x65, 0x4d, 0x61, 0x70, 0x12, 0x1e, 0x0a, 0x0a,
	0x74, 0x6f, 0x74, 0x61, 0x6c, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x0a, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x12, 0x31, 0x0a, 0x06,
	0x71, 0x75, 0x65, 0x75, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x41,
	0x50, 0x49, 0x2e, 0x51, 0x75, 0x65, 0x75, 0x65, 0x4d, 0x61, 0x70, 0x2e, 0x51, 0x75, 0x65, 0x75,
	0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x71, 0x75, 0x65, 0x75, 0x65, 0x73, 0x1a,
	0x4d, 0x0a, 0x0b, 0x51, 0x75, 0x65, 0x75, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x12, 0x28, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x12, 0x2e, 0x41, 0x50, 0x49, 0x2e, 0x51, 0x75, 0x65, 0x75, 0x65, 0x50, 0x72, 0x6f, 0x67, 0x72,
	0x65, 0x73, 0x73, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x2a, 0x2f,
	0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x0b, 0x0a, 0x07, 0x53, 0x55, 0x43, 0x43,
	0x45, 0x53, 0x53, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x46, 0x41, 0x49, 0x4c, 0x55, 0x52, 0x45,
	0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x02, 0x32,
	0xb9, 0x05, 0x0a, 0x03, 0x41, 0x50, 0x49, 0x12, 0x47, 0x0a, 0x07, 0x48, 0x65, 0x61, 0x6c, 0x74,
	0x68, 0x7a, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x12, 0x2e, 0x41, 0x50, 0x49,
	0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x10,
	0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0a, 0x12, 0x08, 0x2f, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x7a,
	0x12, 0x49, 0x0a, 0x08, 0x50, 0x75, 0x73, 0x68, 0x54, 0x61, 0x73, 0x6b, 0x12, 0x09, 0x2e, 0x41,
	0x50, 0x49, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x1a, 0x12, 0x2e, 0x41, 0x50, 0x49, 0x2e, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x1e, 0x82, 0xd3, 0xe4,
	0x93, 0x02, 0x18, 0x3a, 0x01, 0x2a, 0x22, 0x13, 0x2f, 0x76, 0x31, 0x2f, 0x71, 0x75, 0x65, 0x75,
	0x65, 0x2f, 0x7b, 0x71, 0x75, 0x65, 0x75, 0x65, 0x49, 0x44, 0x7d, 0x12, 0x2c, 0x0a, 0x08, 0x4c,
	0x6f, 0x61, 0x64, 0x54, 0x61, 0x73, 0x6b, 0x12, 0x13, 0x2e, 0x41, 0x50, 0x49, 0x2e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x09, 0x2e, 0x41,
	0x50, 0x49, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x22, 0x00, 0x12, 0x54, 0x0a, 0x09, 0x48, 0x65, 0x61,
	0x72, 0x74, 0x42, 0x65, 0x61, 0x74, 0x12, 0x09, 0x2e, 0x41, 0x50, 0x49, 0x2e, 0x54, 0x61, 0x73,
	0x6b, 0x1a, 0x12, 0x2e, 0x41, 0x50, 0x49, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x28, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x22, 0x1a, 0x20, 0x2f,
	0x76, 0x31, 0x2f, 0x68, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x2f, 0x7b, 0x71, 0x75,
	0x65, 0x75, 0x65, 0x49, 0x44, 0x7d, 0x2f, 0x7b, 0x74, 0x61, 0x73, 0x6b, 0x49, 0x44, 0x7d, 0x12,
	0x52, 0x0a, 0x08, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x09, 0x2e, 0x41, 0x50,
	0x49, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x1a, 0x12, 0x2e, 0x41, 0x50, 0x49, 0x2e, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x27, 0x82, 0xd3, 0xe4, 0x93,
	0x02, 0x21, 0x1a, 0x1f, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65,
	0x2f, 0x7b, 0x71, 0x75, 0x65, 0x75, 0x65, 0x49, 0x44, 0x7d, 0x2f, 0x7b, 0x74, 0x61, 0x73, 0x6b,
	0x49, 0x44, 0x7d, 0x12, 0x4a, 0x0a, 0x04, 0x46, 0x61, 0x69, 0x6c, 0x12, 0x09, 0x2e, 0x41, 0x50,
	0x49, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x1a, 0x12, 0x2e, 0x41, 0x50, 0x49, 0x2e, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x23, 0x82, 0xd3, 0xe4, 0x93,
	0x02, 0x1d, 0x1a, 0x1b, 0x2f, 0x76, 0x31, 0x2f, 0x66, 0x61, 0x69, 0x6c, 0x2f, 0x7b, 0x71, 0x75,
	0x65, 0x75, 0x65, 0x49, 0x44, 0x7d, 0x2f, 0x7b, 0x74, 0x61, 0x73, 0x6b, 0x49, 0x44, 0x7d, 0x12,
	0x59, 0x0a, 0x08, 0x50, 0x72, 0x6f, 0x67, 0x72, 0x65, 0x73, 0x73, 0x12, 0x13, 0x2e, 0x41, 0x50,
	0x49, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x1a, 0x12, 0x2e, 0x41, 0x50, 0x49, 0x2e, 0x51, 0x75, 0x65, 0x75, 0x65, 0x50, 0x72, 0x6f, 0x67,
	0x72, 0x65, 0x73, 0x73, 0x22, 0x24, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1e, 0x12, 0x1c, 0x2f, 0x76,
	0x31, 0x2f, 0x71, 0x75, 0x65, 0x75, 0x65, 0x2f, 0x7b, 0x71, 0x75, 0x65, 0x75, 0x65, 0x49, 0x44,
	0x7d, 0x2f, 0x70, 0x72, 0x6f, 0x67, 0x72, 0x65, 0x73, 0x73, 0x12, 0x57, 0x0a, 0x04, 0x50, 0x65,
	0x65, 0x6b, 0x12, 0x13, 0x2e, 0x41, 0x50, 0x49, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x0d, 0x2e, 0x41, 0x50, 0x49, 0x2e, 0x54, 0x61,
	0x73, 0x6b, 0x4c, 0x69, 0x73, 0x74, 0x22, 0x2b, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x25, 0x12, 0x23,
	0x2f, 0x76, 0x31, 0x2f, 0x71, 0x75, 0x65, 0x75, 0x65, 0x2f, 0x7b, 0x71, 0x75, 0x65, 0x75, 0x65,
	0x49, 0x44, 0x7d, 0x2f, 0x7b, 0x70, 0x68, 0x61, 0x73, 0x65, 0x7d, 0x2f, 0x7b, 0x6c, 0x69, 0x6d,
	0x69, 0x74, 0x7d, 0x12, 0x46, 0x0a, 0x0a, 0x4c, 0x69, 0x73, 0x74, 0x51, 0x75, 0x65, 0x75, 0x65,
	0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x0d, 0x2e, 0x41, 0x50, 0x49, 0x2e,
	0x51, 0x75, 0x65, 0x75, 0x65, 0x4d, 0x61, 0x70, 0x22, 0x11, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0b,
	0x12, 0x09, 0x2f, 0x76, 0x31, 0x2f, 0x71, 0x75, 0x65, 0x75, 0x65, 0x42, 0x07, 0x5a, 0x05, 0x2e,
	0x2f, 0x41, 0x50, 0x49, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_proto_rawDescOnce sync.Once
	file_api_proto_rawDescData = file_api_proto_rawDesc
)

func file_api_proto_rawDescGZIP() []byte {
	file_api_proto_rawDescOnce.Do(func() {
		file_api_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_proto_rawDescData)
	})
	return file_api_proto_rawDescData
}

var file_api_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_api_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_api_proto_goTypes = []interface{}{
	(Status)(0),            // 0: API.Status
	(*Task)(nil),           // 1: API.Task
	(*Meta)(nil),           // 2: API.Meta
	(*RequestMessage)(nil), // 3: API.RequestMessage
	(*StatusMessage)(nil),  // 4: API.StatusMessage
	(*QueueProgress)(nil),  // 5: API.QueueProgress
	(*TaskList)(nil),       // 6: API.TaskList
	(*QueueMap)(nil),       // 7: API.QueueMap
	nil,                    // 8: API.QueueMap.QueuesEntry
	(*emptypb.Empty)(nil),  // 9: google.protobuf.Empty
}
var file_api_proto_depIdxs = []int32{
	0,  // 0: API.StatusMessage.status:type_name -> API.Status
	1,  // 1: API.TaskList.tasks:type_name -> API.Task
	8,  // 2: API.QueueMap.queues:type_name -> API.QueueMap.QueuesEntry
	5,  // 3: API.QueueMap.QueuesEntry.value:type_name -> API.QueueProgress
	9,  // 4: API.API.Healthz:input_type -> google.protobuf.Empty
	1,  // 5: API.API.PushTask:input_type -> API.Task
	3,  // 6: API.API.LoadTask:input_type -> API.RequestMessage
	1,  // 7: API.API.HeartBeat:input_type -> API.Task
	1,  // 8: API.API.Complete:input_type -> API.Task
	1,  // 9: API.API.Fail:input_type -> API.Task
	3,  // 10: API.API.Progress:input_type -> API.RequestMessage
	3,  // 11: API.API.Peek:input_type -> API.RequestMessage
	9,  // 12: API.API.ListQueues:input_type -> google.protobuf.Empty
	4,  // 13: API.API.Healthz:output_type -> API.StatusMessage
	4,  // 14: API.API.PushTask:output_type -> API.StatusMessage
	1,  // 15: API.API.LoadTask:output_type -> API.Task
	4,  // 16: API.API.HeartBeat:output_type -> API.StatusMessage
	4,  // 17: API.API.Complete:output_type -> API.StatusMessage
	4,  // 18: API.API.Fail:output_type -> API.StatusMessage
	5,  // 19: API.API.Progress:output_type -> API.QueueProgress
	6,  // 20: API.API.Peek:output_type -> API.TaskList
	7,  // 21: API.API.ListQueues:output_type -> API.QueueMap
	13, // [13:22] is the sub-list for method output_type
	4,  // [4:13] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_api_proto_init() }
func file_api_proto_init() {
	if File_api_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Task); i {
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
		file_api_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Meta); i {
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
		file_api_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestMessage); i {
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
		file_api_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StatusMessage); i {
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
		file_api_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QueueProgress); i {
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
		file_api_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TaskList); i {
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
		file_api_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QueueMap); i {
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
			RawDescriptor: file_api_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_proto_goTypes,
		DependencyIndexes: file_api_proto_depIdxs,
		EnumInfos:         file_api_proto_enumTypes,
		MessageInfos:      file_api_proto_msgTypes,
	}.Build()
	File_api_proto = out.File
	file_api_proto_rawDesc = nil
	file_api_proto_goTypes = nil
	file_api_proto_depIdxs = nil
}
