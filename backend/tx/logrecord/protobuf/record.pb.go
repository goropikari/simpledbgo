// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.20.0
// source: backend/tx/logrecord/protofile/record.proto

package protobuf

import (
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

type SetInt32Record struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Filename    string `protobuf:"bytes,1,opt,name=filename,proto3" json:"filename,omitempty"`
	Txnum       int32  `protobuf:"varint,2,opt,name=txnum,proto3" json:"txnum,omitempty"`
	BlockNumber int32  `protobuf:"varint,3,opt,name=block_number,json=blockNumber,proto3" json:"block_number,omitempty"`
	Offset      int64  `protobuf:"varint,4,opt,name=offset,proto3" json:"offset,omitempty"`
	Val         int32  `protobuf:"varint,5,opt,name=val,proto3" json:"val,omitempty"`
}

func (x *SetInt32Record) Reset() {
	*x = SetInt32Record{}
	if protoimpl.UnsafeEnabled {
		mi := &file_backend_tx_logrecord_protofile_record_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetInt32Record) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetInt32Record) ProtoMessage() {}

func (x *SetInt32Record) ProtoReflect() protoreflect.Message {
	mi := &file_backend_tx_logrecord_protofile_record_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetInt32Record.ProtoReflect.Descriptor instead.
func (*SetInt32Record) Descriptor() ([]byte, []int) {
	return file_backend_tx_logrecord_protofile_record_proto_rawDescGZIP(), []int{0}
}

func (x *SetInt32Record) GetFilename() string {
	if x != nil {
		return x.Filename
	}
	return ""
}

func (x *SetInt32Record) GetTxnum() int32 {
	if x != nil {
		return x.Txnum
	}
	return 0
}

func (x *SetInt32Record) GetBlockNumber() int32 {
	if x != nil {
		return x.BlockNumber
	}
	return 0
}

func (x *SetInt32Record) GetOffset() int64 {
	if x != nil {
		return x.Offset
	}
	return 0
}

func (x *SetInt32Record) GetVal() int32 {
	if x != nil {
		return x.Val
	}
	return 0
}

var File_backend_tx_logrecord_protofile_record_proto protoreflect.FileDescriptor

var file_backend_tx_logrecord_protofile_record_proto_rawDesc = []byte{
	0x0a, 0x2b, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2f, 0x74, 0x78, 0x2f, 0x6c, 0x6f, 0x67,
	0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x66, 0x69, 0x6c, 0x65,
	0x2f, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x22, 0x8f, 0x01, 0x0a, 0x0e, 0x53, 0x65, 0x74, 0x49,
	0x6e, 0x74, 0x33, 0x32, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x69,
	0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69,
	0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x78, 0x6e, 0x75, 0x6d, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x74, 0x78, 0x6e, 0x75, 0x6d, 0x12, 0x21, 0x0a, 0x0c,
	0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x0b, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12,
	0x16, 0x0a, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x76, 0x61, 0x6c, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x76, 0x61, 0x6c, 0x42, 0x0c, 0x5a, 0x0a, 0x2e, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_backend_tx_logrecord_protofile_record_proto_rawDescOnce sync.Once
	file_backend_tx_logrecord_protofile_record_proto_rawDescData = file_backend_tx_logrecord_protofile_record_proto_rawDesc
)

func file_backend_tx_logrecord_protofile_record_proto_rawDescGZIP() []byte {
	file_backend_tx_logrecord_protofile_record_proto_rawDescOnce.Do(func() {
		file_backend_tx_logrecord_protofile_record_proto_rawDescData = protoimpl.X.CompressGZIP(file_backend_tx_logrecord_protofile_record_proto_rawDescData)
	})
	return file_backend_tx_logrecord_protofile_record_proto_rawDescData
}

var file_backend_tx_logrecord_protofile_record_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_backend_tx_logrecord_protofile_record_proto_goTypes = []interface{}{
	(*SetInt32Record)(nil), // 0: protobuf.SetInt32Record
}
var file_backend_tx_logrecord_protofile_record_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_backend_tx_logrecord_protofile_record_proto_init() }
func file_backend_tx_logrecord_protofile_record_proto_init() {
	if File_backend_tx_logrecord_protofile_record_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_backend_tx_logrecord_protofile_record_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SetInt32Record); i {
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
			RawDescriptor: file_backend_tx_logrecord_protofile_record_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_backend_tx_logrecord_protofile_record_proto_goTypes,
		DependencyIndexes: file_backend_tx_logrecord_protofile_record_proto_depIdxs,
		MessageInfos:      file_backend_tx_logrecord_protofile_record_proto_msgTypes,
	}.Build()
	File_backend_tx_logrecord_protofile_record_proto = out.File
	file_backend_tx_logrecord_protofile_record_proto_rawDesc = nil
	file_backend_tx_logrecord_protofile_record_proto_goTypes = nil
	file_backend_tx_logrecord_protofile_record_proto_depIdxs = nil
}
