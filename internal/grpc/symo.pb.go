// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.14.0
// source: symo.proto

package grpc

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
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

type LoadAvg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Load1  float64 `protobuf:"fixed64,1,opt,name=Load1,proto3" json:"Load1,omitempty"`
	Load5  float64 `protobuf:"fixed64,2,opt,name=Load5,proto3" json:"Load5,omitempty"`
	Load15 float64 `protobuf:"fixed64,3,opt,name=Load15,proto3" json:"Load15,omitempty"`
}

func (x *LoadAvg) Reset() {
	*x = LoadAvg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_symo_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoadAvg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoadAvg) ProtoMessage() {}

func (x *LoadAvg) ProtoReflect() protoreflect.Message {
	mi := &file_symo_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoadAvg.ProtoReflect.Descriptor instead.
func (*LoadAvg) Descriptor() ([]byte, []int) {
	return file_symo_proto_rawDescGZIP(), []int{0}
}

func (x *LoadAvg) GetLoad1() float64 {
	if x != nil {
		return x.Load1
	}
	return 0
}

func (x *LoadAvg) GetLoad5() float64 {
	if x != nil {
		return x.Load5
	}
	return 0
}

func (x *LoadAvg) GetLoad15() float64 {
	if x != nil {
		return x.Load15
	}
	return 0
}

type CPU struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	User   float64 `protobuf:"fixed64,1,opt,name=User,proto3" json:"User,omitempty"`
	System float64 `protobuf:"fixed64,2,opt,name=System,proto3" json:"System,omitempty"`
	Idle   float64 `protobuf:"fixed64,3,opt,name=Idle,proto3" json:"Idle,omitempty"`
}

func (x *CPU) Reset() {
	*x = CPU{}
	if protoimpl.UnsafeEnabled {
		mi := &file_symo_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CPU) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CPU) ProtoMessage() {}

func (x *CPU) ProtoReflect() protoreflect.Message {
	mi := &file_symo_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CPU.ProtoReflect.Descriptor instead.
func (*CPU) Descriptor() ([]byte, []int) {
	return file_symo_proto_rawDescGZIP(), []int{1}
}

func (x *CPU) GetUser() float64 {
	if x != nil {
		return x.User
	}
	return 0
}

func (x *CPU) GetSystem() float64 {
	if x != nil {
		return x.System
	}
	return 0
}

func (x *CPU) GetIdle() float64 {
	if x != nil {
		return x.Idle
	}
	return 0
}

type LoadDisk struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name    string  `protobuf:"bytes,1,opt,name=Name,proto3" json:"Name,omitempty"`
	Tps     float64 `protobuf:"fixed64,2,opt,name=Tps,proto3" json:"Tps,omitempty"`
	KBRead  float64 `protobuf:"fixed64,3,opt,name=KBRead,proto3" json:"KBRead,omitempty"`
	KBWrite float64 `protobuf:"fixed64,4,opt,name=KBWrite,proto3" json:"KBWrite,omitempty"`
}

func (x *LoadDisk) Reset() {
	*x = LoadDisk{}
	if protoimpl.UnsafeEnabled {
		mi := &file_symo_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoadDisk) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoadDisk) ProtoMessage() {}

func (x *LoadDisk) ProtoReflect() protoreflect.Message {
	mi := &file_symo_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoadDisk.ProtoReflect.Descriptor instead.
func (*LoadDisk) Descriptor() ([]byte, []int) {
	return file_symo_proto_rawDescGZIP(), []int{2}
}

func (x *LoadDisk) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *LoadDisk) GetTps() float64 {
	if x != nil {
		return x.Tps
	}
	return 0
}

func (x *LoadDisk) GetKBRead() float64 {
	if x != nil {
		return x.KBRead
	}
	return 0
}

func (x *LoadDisk) GetKBWrite() float64 {
	if x != nil {
		return x.KBWrite
	}
	return 0
}

type UsedFS struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Path      string  `protobuf:"bytes,1,opt,name=Path,proto3" json:"Path,omitempty"`
	UsedSpace float64 `protobuf:"fixed64,2,opt,name=UsedSpace,proto3" json:"UsedSpace,omitempty"`
	UsedInode float64 `protobuf:"fixed64,3,opt,name=UsedInode,proto3" json:"UsedInode,omitempty"`
}

func (x *UsedFS) Reset() {
	*x = UsedFS{}
	if protoimpl.UnsafeEnabled {
		mi := &file_symo_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UsedFS) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UsedFS) ProtoMessage() {}

func (x *UsedFS) ProtoReflect() protoreflect.Message {
	mi := &file_symo_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UsedFS.ProtoReflect.Descriptor instead.
func (*UsedFS) Descriptor() ([]byte, []int) {
	return file_symo_proto_rawDescGZIP(), []int{3}
}

func (x *UsedFS) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *UsedFS) GetUsedSpace() float64 {
	if x != nil {
		return x.UsedSpace
	}
	return 0
}

func (x *UsedFS) GetUsedInode() float64 {
	if x != nil {
		return x.UsedInode
	}
	return 0
}

type Stats struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Time      *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=time,proto3" json:"time,omitempty"`
	LoadAvg   *LoadAvg               `protobuf:"bytes,2,opt,name=load_avg,json=loadAvg,proto3" json:"load_avg,omitempty"`
	Cpu       *CPU                   `protobuf:"bytes,3,opt,name=cpu,proto3" json:"cpu,omitempty"`
	LoadDisks []*LoadDisk            `protobuf:"bytes,4,rep,name=load_disks,json=loadDisks,proto3" json:"load_disks,omitempty"`
	UsedFs    []*UsedFS              `protobuf:"bytes,5,rep,name=used_fs,json=usedFs,proto3" json:"used_fs,omitempty"`
}

func (x *Stats) Reset() {
	*x = Stats{}
	if protoimpl.UnsafeEnabled {
		mi := &file_symo_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Stats) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Stats) ProtoMessage() {}

func (x *Stats) ProtoReflect() protoreflect.Message {
	mi := &file_symo_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Stats.ProtoReflect.Descriptor instead.
func (*Stats) Descriptor() ([]byte, []int) {
	return file_symo_proto_rawDescGZIP(), []int{4}
}

func (x *Stats) GetTime() *timestamppb.Timestamp {
	if x != nil {
		return x.Time
	}
	return nil
}

func (x *Stats) GetLoadAvg() *LoadAvg {
	if x != nil {
		return x.LoadAvg
	}
	return nil
}

func (x *Stats) GetCpu() *CPU {
	if x != nil {
		return x.Cpu
	}
	return nil
}

func (x *Stats) GetLoadDisks() []*LoadDisk {
	if x != nil {
		return x.LoadDisks
	}
	return nil
}

func (x *Stats) GetUsedFs() []*UsedFS {
	if x != nil {
		return x.UsedFs
	}
	return nil
}

type StatsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	N int32 `protobuf:"varint,1,opt,name=N,proto3" json:"N,omitempty"`
	M int32 `protobuf:"varint,2,opt,name=M,proto3" json:"M,omitempty"`
}

func (x *StatsRequest) Reset() {
	*x = StatsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_symo_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StatsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StatsRequest) ProtoMessage() {}

func (x *StatsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_symo_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StatsRequest.ProtoReflect.Descriptor instead.
func (*StatsRequest) Descriptor() ([]byte, []int) {
	return file_symo_proto_rawDescGZIP(), []int{5}
}

func (x *StatsRequest) GetN() int32 {
	if x != nil {
		return x.N
	}
	return 0
}

func (x *StatsRequest) GetM() int32 {
	if x != nil {
		return x.M
	}
	return 0
}

var File_symo_proto protoreflect.FileDescriptor

var file_symo_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x73, 0x79, 0x6d, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x73, 0x74,
	0x61, 0x74, 0x73, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x4d, 0x0a, 0x07, 0x4c, 0x6f, 0x61, 0x64, 0x41, 0x76, 0x67, 0x12,
	0x14, 0x0a, 0x05, 0x4c, 0x6f, 0x61, 0x64, 0x31, 0x18, 0x01, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05,
	0x4c, 0x6f, 0x61, 0x64, 0x31, 0x12, 0x14, 0x0a, 0x05, 0x4c, 0x6f, 0x61, 0x64, 0x35, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x4c, 0x6f, 0x61, 0x64, 0x35, 0x12, 0x16, 0x0a, 0x06, 0x4c,
	0x6f, 0x61, 0x64, 0x31, 0x35, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x06, 0x4c, 0x6f, 0x61,
	0x64, 0x31, 0x35, 0x22, 0x45, 0x0a, 0x03, 0x43, 0x50, 0x55, 0x12, 0x12, 0x0a, 0x04, 0x55, 0x73,
	0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x01, 0x52, 0x04, 0x55, 0x73, 0x65, 0x72, 0x12, 0x16,
	0x0a, 0x06, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x06,
	0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x12, 0x12, 0x0a, 0x04, 0x49, 0x64, 0x6c, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x01, 0x52, 0x04, 0x49, 0x64, 0x6c, 0x65, 0x22, 0x62, 0x0a, 0x08, 0x4c, 0x6f,
	0x61, 0x64, 0x44, 0x69, 0x73, 0x6b, 0x12, 0x12, 0x0a, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x54, 0x70,
	0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x03, 0x54, 0x70, 0x73, 0x12, 0x16, 0x0a, 0x06,
	0x4b, 0x42, 0x52, 0x65, 0x61, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x06, 0x4b, 0x42,
	0x52, 0x65, 0x61, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x4b, 0x42, 0x57, 0x72, 0x69, 0x74, 0x65, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x01, 0x52, 0x07, 0x4b, 0x42, 0x57, 0x72, 0x69, 0x74, 0x65, 0x22, 0x58,
	0x0a, 0x06, 0x55, 0x73, 0x65, 0x64, 0x46, 0x53, 0x12, 0x12, 0x0a, 0x04, 0x50, 0x61, 0x74, 0x68,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x50, 0x61, 0x74, 0x68, 0x12, 0x1c, 0x0a, 0x09,
	0x55, 0x73, 0x65, 0x64, 0x53, 0x70, 0x61, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52,
	0x09, 0x55, 0x73, 0x65, 0x64, 0x53, 0x70, 0x61, 0x63, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x55, 0x73,
	0x65, 0x64, 0x49, 0x6e, 0x6f, 0x64, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x09, 0x55,
	0x73, 0x65, 0x64, 0x49, 0x6e, 0x6f, 0x64, 0x65, 0x22, 0xd8, 0x01, 0x0a, 0x05, 0x53, 0x74, 0x61,
	0x74, 0x73, 0x12, 0x2e, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x04, 0x74, 0x69,
	0x6d, 0x65, 0x12, 0x29, 0x0a, 0x08, 0x6c, 0x6f, 0x61, 0x64, 0x5f, 0x61, 0x76, 0x67, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x73, 0x2e, 0x4c, 0x6f, 0x61,
	0x64, 0x41, 0x76, 0x67, 0x52, 0x07, 0x6c, 0x6f, 0x61, 0x64, 0x41, 0x76, 0x67, 0x12, 0x1c, 0x0a,
	0x03, 0x63, 0x70, 0x75, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x73, 0x74, 0x61,
	0x74, 0x73, 0x2e, 0x43, 0x50, 0x55, 0x52, 0x03, 0x63, 0x70, 0x75, 0x12, 0x2e, 0x0a, 0x0a, 0x6c,
	0x6f, 0x61, 0x64, 0x5f, 0x64, 0x69, 0x73, 0x6b, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x0f, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x73, 0x2e, 0x4c, 0x6f, 0x61, 0x64, 0x44, 0x69, 0x73, 0x6b,
	0x52, 0x09, 0x6c, 0x6f, 0x61, 0x64, 0x44, 0x69, 0x73, 0x6b, 0x73, 0x12, 0x26, 0x0a, 0x07, 0x75,
	0x73, 0x65, 0x64, 0x5f, 0x66, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x73,
	0x74, 0x61, 0x74, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x64, 0x46, 0x53, 0x52, 0x06, 0x75, 0x73, 0x65,
	0x64, 0x46, 0x73, 0x22, 0x2a, 0x0a, 0x0c, 0x53, 0x74, 0x61, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x0c, 0x0a, 0x01, 0x4e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x01,
	0x4e, 0x12, 0x0c, 0x0a, 0x01, 0x4d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x01, 0x4d, 0x32,
	0x39, 0x0a, 0x04, 0x53, 0x79, 0x6d, 0x6f, 0x12, 0x31, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x53, 0x74,
	0x61, 0x74, 0x73, 0x12, 0x13, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x73, 0x2e, 0x53, 0x74, 0x61, 0x74,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0c, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x73,
	0x2e, 0x53, 0x74, 0x61, 0x74, 0x73, 0x22, 0x00, 0x30, 0x01, 0x42, 0x08, 0x5a, 0x06, 0x2e, 0x3b,
	0x67, 0x72, 0x70, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_symo_proto_rawDescOnce sync.Once
	file_symo_proto_rawDescData = file_symo_proto_rawDesc
)

func file_symo_proto_rawDescGZIP() []byte {
	file_symo_proto_rawDescOnce.Do(func() {
		file_symo_proto_rawDescData = protoimpl.X.CompressGZIP(file_symo_proto_rawDescData)
	})
	return file_symo_proto_rawDescData
}

var file_symo_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_symo_proto_goTypes = []interface{}{
	(*LoadAvg)(nil),               // 0: stats.LoadAvg
	(*CPU)(nil),                   // 1: stats.CPU
	(*LoadDisk)(nil),              // 2: stats.LoadDisk
	(*UsedFS)(nil),                // 3: stats.UsedFS
	(*Stats)(nil),                 // 4: stats.Stats
	(*StatsRequest)(nil),          // 5: stats.StatsRequest
	(*timestamppb.Timestamp)(nil), // 6: google.protobuf.Timestamp
}
var file_symo_proto_depIdxs = []int32{
	6, // 0: stats.Stats.time:type_name -> google.protobuf.Timestamp
	0, // 1: stats.Stats.load_avg:type_name -> stats.LoadAvg
	1, // 2: stats.Stats.cpu:type_name -> stats.CPU
	2, // 3: stats.Stats.load_disks:type_name -> stats.LoadDisk
	3, // 4: stats.Stats.used_fs:type_name -> stats.UsedFS
	5, // 5: stats.Symo.GetStats:input_type -> stats.StatsRequest
	4, // 6: stats.Symo.GetStats:output_type -> stats.Stats
	6, // [6:7] is the sub-list for method output_type
	5, // [5:6] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_symo_proto_init() }
func file_symo_proto_init() {
	if File_symo_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_symo_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LoadAvg); i {
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
		file_symo_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CPU); i {
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
		file_symo_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LoadDisk); i {
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
		file_symo_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UsedFS); i {
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
		file_symo_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Stats); i {
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
		file_symo_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StatsRequest); i {
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
			RawDescriptor: file_symo_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_symo_proto_goTypes,
		DependencyIndexes: file_symo_proto_depIdxs,
		MessageInfos:      file_symo_proto_msgTypes,
	}.Build()
	File_symo_proto = out.File
	file_symo_proto_rawDesc = nil
	file_symo_proto_goTypes = nil
	file_symo_proto_depIdxs = nil
}
