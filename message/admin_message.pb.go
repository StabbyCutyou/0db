// Code generated by protoc-gen-go.
// source: admin_message.proto
// DO NOT EDIT!

/*
Package message is a generated protocol buffer package.

It is generated from these files:
	admin_message.proto

It has these top-level messages:
	AdminMessage
*/
package message

import proto "github.com/golang/protobuf/proto"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

type AdminMessage struct {
	Command          *string `protobuf:"bytes,1,req,name=command" json:"command,omitempty"`
	Message          *string `protobuf:"bytes,2,req,name=message" json:"message,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *AdminMessage) Reset()         { *m = AdminMessage{} }
func (m *AdminMessage) String() string { return proto.CompactTextString(m) }
func (*AdminMessage) ProtoMessage()    {}

func (m *AdminMessage) GetCommand() string {
	if m != nil && m.Command != nil {
		return *m.Command
	}
	return ""
}

func (m *AdminMessage) GetMessage() string {
	if m != nil && m.Message != nil {
		return *m.Message
	}
	return ""
}

func init() {
}
