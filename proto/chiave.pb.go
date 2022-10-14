// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.7
// source: proto/chiave.proto

package chiave

import (
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

type Key struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *Key) Reset() {
	*x = Key{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_chiave_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Key) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Key) ProtoMessage() {}

func (x *Key) ProtoReflect() protoreflect.Message {
	mi := &file_proto_chiave_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Key.ProtoReflect.Descriptor instead.
func (*Key) Descriptor() ([]byte, []int) {
	return file_proto_chiave_proto_rawDescGZIP(), []int{0}
}

func (x *Key) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type ValueResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value  int64 `protobuf:"varint,1,opt,name=value,proto3" json:"value,omitempty"`
	Exists bool  `protobuf:"varint,2,opt,name=exists,proto3" json:"exists,omitempty"`
}

func (x *ValueResponse) Reset() {
	*x = ValueResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_chiave_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ValueResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ValueResponse) ProtoMessage() {}

func (x *ValueResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_chiave_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ValueResponse.ProtoReflect.Descriptor instead.
func (*ValueResponse) Descriptor() ([]byte, []int) {
	return file_proto_chiave_proto_rawDescGZIP(), []int{1}
}

func (x *ValueResponse) GetValue() int64 {
	if x != nil {
		return x.Value
	}
	return 0
}

func (x *ValueResponse) GetExists() bool {
	if x != nil {
		return x.Exists
	}
	return false
}

var File_proto_chiave_proto protoreflect.FileDescriptor

var file_proto_chiave_proto_rawDesc = []byte{
	0x0a, 0x12, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x68, 0x69, 0x61, 0x76, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x63, 0x68, 0x69, 0x61, 0x76, 0x65, 0x1a, 0x1b, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d,
	0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x15, 0x0a, 0x03, 0x4b, 0x65, 0x79,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64,
	0x22, 0x3d, 0x0a, 0x0d, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x65, 0x78, 0x69, 0x73, 0x74,
	0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x65, 0x78, 0x69, 0x73, 0x74, 0x73, 0x32,
	0x9f, 0x01, 0x0a, 0x06, 0x43, 0x68, 0x69, 0x61, 0x76, 0x65, 0x12, 0x2d, 0x0a, 0x05, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x12, 0x0b, 0x2e, 0x63, 0x68, 0x69, 0x61, 0x76, 0x65, 0x2e, 0x4b, 0x65, 0x79,
	0x1a, 0x15, 0x2e, 0x63, 0x68, 0x69, 0x61, 0x76, 0x65, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x32, 0x0a, 0x09, 0x49, 0x6e, 0x63,
	0x72, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x0b, 0x2e, 0x63, 0x68, 0x69, 0x61, 0x76, 0x65, 0x2e,
	0x4b, 0x65, 0x79, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x32, 0x0a,
	0x09, 0x44, 0x65, 0x63, 0x72, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x0b, 0x2e, 0x63, 0x68, 0x69,
	0x61, 0x76, 0x65, 0x2e, 0x4b, 0x65, 0x79, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22,
	0x00, 0x42, 0x12, 0x5a, 0x10, 0x6b, 0x76, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63,
	0x68, 0x69, 0x61, 0x76, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_chiave_proto_rawDescOnce sync.Once
	file_proto_chiave_proto_rawDescData = file_proto_chiave_proto_rawDesc
)

func file_proto_chiave_proto_rawDescGZIP() []byte {
	file_proto_chiave_proto_rawDescOnce.Do(func() {
		file_proto_chiave_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_chiave_proto_rawDescData)
	})
	return file_proto_chiave_proto_rawDescData
}

var file_proto_chiave_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_chiave_proto_goTypes = []interface{}{
	(*Key)(nil),           // 0: chiave.Key
	(*ValueResponse)(nil), // 1: chiave.ValueResponse
	(*emptypb.Empty)(nil), // 2: google.protobuf.Empty
}
var file_proto_chiave_proto_depIdxs = []int32{
	0, // 0: chiave.Chiave.Value:input_type -> chiave.Key
	0, // 1: chiave.Chiave.Increment:input_type -> chiave.Key
	0, // 2: chiave.Chiave.Decrement:input_type -> chiave.Key
	1, // 3: chiave.Chiave.Value:output_type -> chiave.ValueResponse
	2, // 4: chiave.Chiave.Increment:output_type -> google.protobuf.Empty
	2, // 5: chiave.Chiave.Decrement:output_type -> google.protobuf.Empty
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_chiave_proto_init() }
func file_proto_chiave_proto_init() {
	if File_proto_chiave_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_chiave_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Key); i {
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
		file_proto_chiave_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ValueResponse); i {
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
			RawDescriptor: file_proto_chiave_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_chiave_proto_goTypes,
		DependencyIndexes: file_proto_chiave_proto_depIdxs,
		MessageInfos:      file_proto_chiave_proto_msgTypes,
	}.Build()
	File_proto_chiave_proto = out.File
	file_proto_chiave_proto_rawDesc = nil
	file_proto_chiave_proto_goTypes = nil
	file_proto_chiave_proto_depIdxs = nil
}
