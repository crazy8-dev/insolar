// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: ledger/object/lifeline.proto

package object

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	github_com_insolar_insolar_insolar "github.com/insolar/insolar/insolar"
	github_com_insolar_insolar_insolar_record "github.com/insolar/insolar/insolar/record"
	io "io"
	math "math"
	reflect "reflect"
	strings "strings"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type Lifeline struct {
	XPolymorph          int32                                             `protobuf:"varint,16,opt,name=__polymorph,json=Polymorph,proto3" json:"__polymorph,omitempty"`
	LatestState         *github_com_insolar_insolar_insolar.ID            `protobuf:"bytes,20,opt,name=LatestState,proto3,customtype=github.com/insolar/insolar/insolar.ID" json:"LatestState,omitempty"`
	LatestStateApproved *github_com_insolar_insolar_insolar.ID            `protobuf:"bytes,21,opt,name=LatestStateApproved,proto3,customtype=github.com/insolar/insolar/insolar.ID" json:"LatestStateApproved,omitempty"`
	ChildPointer        *github_com_insolar_insolar_insolar.ID            `protobuf:"bytes,22,opt,name=ChildPointer,proto3,customtype=github.com/insolar/insolar/insolar.ID" json:"ChildPointer,omitempty"`
	Parent              github_com_insolar_insolar_insolar.Reference      `protobuf:"bytes,23,opt,name=Parent,proto3,customtype=github.com/insolar/insolar/insolar.Reference" json:"Parent"`
	Delegates           []LifelineDelegate                                `protobuf:"bytes,24,rep,name=Delegates,proto3" json:"Delegates"`
	StateID             github_com_insolar_insolar_insolar_record.StateID `protobuf:"varint,25,opt,name=StateID,proto3,customtype=github.com/insolar/insolar/insolar/record.StateID" json:"StateID"`
	LatestUpdate        github_com_insolar_insolar_insolar.PulseNumber    `protobuf:"varint,26,opt,name=LatestUpdate,proto3,customtype=github.com/insolar/insolar/insolar.PulseNumber" json:"LatestUpdate"`
	JetID               github_com_insolar_insolar_insolar.JetID          `protobuf:"bytes,27,opt,name=JetID,proto3,customtype=github.com/insolar/insolar/insolar.JetID" json:"JetID"`
	LatestRequest       *github_com_insolar_insolar_insolar.ID            `protobuf:"bytes,28,opt,name=LatestRequest,proto3,customtype=github.com/insolar/insolar/insolar.ID" json:"LatestRequest,omitempty"`
}

func (m *Lifeline) Reset()      { *m = Lifeline{} }
func (*Lifeline) ProtoMessage() {}
func (*Lifeline) Descriptor() ([]byte, []int) {
	return fileDescriptor_4ca0c51acb0e6740, []int{0}
}
func (m *Lifeline) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Lifeline) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Lifeline.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Lifeline) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Lifeline.Merge(m, src)
}
func (m *Lifeline) XXX_Size() int {
	return m.Size()
}
func (m *Lifeline) XXX_DiscardUnknown() {
	xxx_messageInfo_Lifeline.DiscardUnknown(m)
}

var xxx_messageInfo_Lifeline proto.InternalMessageInfo

func (m *Lifeline) GetXPolymorph() int32 {
	if m != nil {
		return m.XPolymorph
	}
	return 0
}

func (m *Lifeline) GetDelegates() []LifelineDelegate {
	if m != nil {
		return m.Delegates
	}
	return nil
}

type LifelineDelegate struct {
	XPolymorph int32                                        `protobuf:"varint,16,opt,name=__polymorph,json=Polymorph,proto3" json:"__polymorph,omitempty"`
	Key        github_com_insolar_insolar_insolar.Reference `protobuf:"bytes,20,opt,name=Key,proto3,customtype=github.com/insolar/insolar/insolar.Reference" json:"Key"`
	Value      github_com_insolar_insolar_insolar.Reference `protobuf:"bytes,21,opt,name=Value,proto3,customtype=github.com/insolar/insolar/insolar.Reference" json:"Value"`
}

func (m *LifelineDelegate) Reset()      { *m = LifelineDelegate{} }
func (*LifelineDelegate) ProtoMessage() {}
func (*LifelineDelegate) Descriptor() ([]byte, []int) {
	return fileDescriptor_4ca0c51acb0e6740, []int{1}
}
func (m *LifelineDelegate) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *LifelineDelegate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_LifelineDelegate.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *LifelineDelegate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LifelineDelegate.Merge(m, src)
}
func (m *LifelineDelegate) XXX_Size() int {
	return m.Size()
}
func (m *LifelineDelegate) XXX_DiscardUnknown() {
	xxx_messageInfo_LifelineDelegate.DiscardUnknown(m)
}

var xxx_messageInfo_LifelineDelegate proto.InternalMessageInfo

func (m *LifelineDelegate) GetXPolymorph() int32 {
	if m != nil {
		return m.XPolymorph
	}
	return 0
}

func init() {
	proto.RegisterType((*Lifeline)(nil), "object.Lifeline")
	proto.RegisterType((*LifelineDelegate)(nil), "object.LifelineDelegate")
}

func init() { proto.RegisterFile("ledger/object/lifeline.proto", fileDescriptor_4ca0c51acb0e6740) }

var fileDescriptor_4ca0c51acb0e6740 = []byte{
	// 500 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0xd4, 0x4f, 0x6f, 0xd3, 0x30,
	0x14, 0x00, 0xf0, 0x58, 0x5b, 0x0b, 0x73, 0x3b, 0x69, 0x32, 0xff, 0xcc, 0x98, 0xdc, 0x6a, 0x12,
	0x52, 0x0e, 0x2c, 0xe1, 0xcf, 0x84, 0x04, 0xe2, 0x42, 0xa9, 0x26, 0x75, 0x54, 0x50, 0x52, 0xc1,
	0x81, 0x4b, 0x95, 0x34, 0xaf, 0x69, 0x50, 0x1a, 0x07, 0xd7, 0x41, 0xda, 0x8d, 0x6f, 0x00, 0x1f,
	0x81, 0x23, 0x9f, 0x04, 0xf5, 0xd8, 0xe3, 0xb4, 0x43, 0x45, 0xd3, 0xcb, 0x8e, 0xfb, 0x08, 0xa8,
	0x4e, 0x2a, 0xda, 0x69, 0xd2, 0xaa, 0x9c, 0x12, 0x3b, 0x79, 0x3f, 0xdb, 0xef, 0xe5, 0x05, 0xef,
	0x05, 0xe0, 0x7a, 0x20, 0x4c, 0xee, 0x7c, 0x81, 0xae, 0x34, 0x03, 0xbf, 0x07, 0x81, 0x1f, 0x82,
	0x11, 0x09, 0x2e, 0x39, 0x29, 0xa6, 0xd3, 0xbb, 0x07, 0x9e, 0x2f, 0xfb, 0xb1, 0x63, 0x74, 0xf9,
	0xc0, 0xf4, 0xb8, 0xc7, 0x4d, 0xf5, 0xd8, 0x89, 0x7b, 0x6a, 0xa4, 0x06, 0xea, 0x2e, 0x0d, 0xdb,
	0xff, 0x51, 0xc4, 0x37, 0x9b, 0x99, 0x44, 0x18, 0x2e, 0x75, 0x3a, 0x11, 0x0f, 0x4e, 0x06, 0x5c,
	0x44, 0x7d, 0xba, 0x53, 0x45, 0x7a, 0xc1, 0xda, 0x6a, 0x2d, 0x26, 0xc8, 0x7b, 0x5c, 0x6a, 0xda,
	0x12, 0x86, 0xb2, 0x2d, 0x6d, 0x09, 0xf4, 0x76, 0x15, 0xe9, 0xe5, 0xda, 0xc1, 0x68, 0x52, 0x41,
	0x67, 0x93, 0xca, 0xc3, 0xa5, 0x85, 0xfd, 0x70, 0xc8, 0x03, 0x5b, 0x5c, 0xbe, 0x1a, 0x8d, 0xba,
	0xb5, 0x2c, 0x90, 0x0e, 0xbe, 0xb5, 0x34, 0x7c, 0x1d, 0x45, 0x82, 0x7f, 0x03, 0x97, 0xde, 0xc9,
	0x03, 0x5f, 0x25, 0x91, 0x0f, 0xb8, 0xfc, 0xa6, 0xef, 0x07, 0x6e, 0x8b, 0xfb, 0xa1, 0x04, 0x41,
	0xef, 0xe6, 0x91, 0x57, 0x08, 0xd2, 0xc4, 0xc5, 0x96, 0x2d, 0x20, 0x94, 0xf4, 0x9e, 0xc2, 0x0e,
	0x47, 0x93, 0x8a, 0x76, 0x36, 0xa9, 0x3c, 0x5a, 0x03, 0xb3, 0xa0, 0x07, 0x02, 0xc2, 0x2e, 0x58,
	0x99, 0x41, 0x5e, 0xe1, 0xad, 0x3a, 0x04, 0xe0, 0xcd, 0xf7, 0x4e, 0x69, 0x75, 0x43, 0x2f, 0x3d,
	0xa5, 0x46, 0x5a, 0x4a, 0x63, 0x51, 0x97, 0xc5, 0x0b, 0xb5, 0xcd, 0xf9, 0x52, 0xd6, 0xff, 0x00,
	0xd2, 0xc6, 0x37, 0xd4, 0x79, 0x1b, 0x75, 0x7a, 0xbf, 0x8a, 0xf4, 0xed, 0xda, 0x8b, 0x6c, 0x33,
	0x4f, 0xae, 0xdf, 0x8c, 0x29, 0xa0, 0xcb, 0x85, 0x6b, 0x64, 0x80, 0xb5, 0x90, 0xc8, 0x67, 0x5c,
	0x4e, 0x53, 0xf9, 0x31, 0x72, 0xe7, 0x65, 0xde, 0x55, 0xf2, 0xf3, 0x4c, 0x36, 0xd6, 0x38, 0x66,
	0x2b, 0x0e, 0x86, 0xf0, 0x2e, 0x1e, 0x38, 0x20, 0xac, 0x15, 0x8b, 0x1c, 0xe1, 0xc2, 0x31, 0xc8,
	0x46, 0x9d, 0x3e, 0x50, 0xb9, 0x7b, 0x9c, 0xa1, 0xfa, 0x1a, 0xa8, 0x8a, 0xb3, 0xd2, 0x70, 0xd2,
	0xc6, 0xdb, 0xa9, 0x6b, 0xc1, 0xd7, 0x18, 0x86, 0x92, 0xee, 0xe5, 0x29, 0xec, 0xaa, 0xf1, 0x72,
	0xf3, 0xfc, 0x57, 0x45, 0xdb, 0xff, 0x83, 0xf0, 0xce, 0xe5, 0xcc, 0x5f, 0xdb, 0x19, 0x47, 0x78,
	0xe3, 0x2d, 0x9c, 0x64, 0x1d, 0x91, 0xef, 0x8b, 0x98, 0x03, 0xe4, 0x18, 0x17, 0x3e, 0xd9, 0x41,
	0x0c, 0x59, 0x0b, 0xe4, 0x93, 0x52, 0xa2, 0x76, 0x38, 0x9e, 0x32, 0xed, 0x74, 0xca, 0xb4, 0x8b,
	0x29, 0x43, 0xdf, 0x13, 0x86, 0x7e, 0x27, 0x0c, 0x8d, 0x12, 0x86, 0xc6, 0x09, 0x43, 0x7f, 0x13,
	0x86, 0xce, 0x13, 0xa6, 0x5d, 0x24, 0x0c, 0xfd, 0x9c, 0x31, 0x6d, 0x3c, 0x63, 0xda, 0xe9, 0x8c,
	0x69, 0x4e, 0x51, 0xfd, 0x17, 0x9e, 0xfd, 0x0b, 0x00, 0x00, 0xff, 0xff, 0xcb, 0x11, 0xa7, 0x96,
	0x6e, 0x04, 0x00, 0x00,
}

func (this *LifelineDelegate) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*LifelineDelegate)
	if !ok {
		that2, ok := that.(LifelineDelegate)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.XPolymorph != that1.XPolymorph {
		return false
	}
	if !this.Key.Equal(that1.Key) {
		return false
	}
	if !this.Value.Equal(that1.Value) {
		return false
	}
	return true
}
func (this *Lifeline) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 14)
	s = append(s, "&object.Lifeline{")
	s = append(s, "XPolymorph: "+fmt.Sprintf("%#v", this.XPolymorph)+",\n")
	s = append(s, "LatestState: "+fmt.Sprintf("%#v", this.LatestState)+",\n")
	s = append(s, "LatestStateApproved: "+fmt.Sprintf("%#v", this.LatestStateApproved)+",\n")
	s = append(s, "ChildPointer: "+fmt.Sprintf("%#v", this.ChildPointer)+",\n")
	s = append(s, "Parent: "+fmt.Sprintf("%#v", this.Parent)+",\n")
	if this.Delegates != nil {
		vs := make([]*LifelineDelegate, len(this.Delegates))
		for i := range vs {
			vs[i] = &this.Delegates[i]
		}
		s = append(s, "Delegates: "+fmt.Sprintf("%#v", vs)+",\n")
	}
	s = append(s, "StateID: "+fmt.Sprintf("%#v", this.StateID)+",\n")
	s = append(s, "LatestUpdate: "+fmt.Sprintf("%#v", this.LatestUpdate)+",\n")
	s = append(s, "JetID: "+fmt.Sprintf("%#v", this.JetID)+",\n")
	s = append(s, "LatestRequest: "+fmt.Sprintf("%#v", this.LatestRequest)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *LifelineDelegate) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 7)
	s = append(s, "&object.LifelineDelegate{")
	s = append(s, "XPolymorph: "+fmt.Sprintf("%#v", this.XPolymorph)+",\n")
	s = append(s, "Key: "+fmt.Sprintf("%#v", this.Key)+",\n")
	s = append(s, "Value: "+fmt.Sprintf("%#v", this.Value)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func valueToGoStringLifeline(v interface{}, typ string) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)
}
func (m *Lifeline) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Lifeline) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.XPolymorph != 0 {
		dAtA[i] = 0x80
		i++
		dAtA[i] = 0x1
		i++
		i = encodeVarintLifeline(dAtA, i, uint64(m.XPolymorph))
	}
	if m.LatestState != nil {
		dAtA[i] = 0xa2
		i++
		dAtA[i] = 0x1
		i++
		i = encodeVarintLifeline(dAtA, i, uint64(m.LatestState.Size()))
		n1, err := m.LatestState.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	if m.LatestStateApproved != nil {
		dAtA[i] = 0xaa
		i++
		dAtA[i] = 0x1
		i++
		i = encodeVarintLifeline(dAtA, i, uint64(m.LatestStateApproved.Size()))
		n2, err := m.LatestStateApproved.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n2
	}
	if m.ChildPointer != nil {
		dAtA[i] = 0xb2
		i++
		dAtA[i] = 0x1
		i++
		i = encodeVarintLifeline(dAtA, i, uint64(m.ChildPointer.Size()))
		n3, err := m.ChildPointer.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n3
	}
	dAtA[i] = 0xba
	i++
	dAtA[i] = 0x1
	i++
	i = encodeVarintLifeline(dAtA, i, uint64(m.Parent.Size()))
	n4, err := m.Parent.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n4
	if len(m.Delegates) > 0 {
		for _, msg := range m.Delegates {
			dAtA[i] = 0xc2
			i++
			dAtA[i] = 0x1
			i++
			i = encodeVarintLifeline(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if m.StateID != 0 {
		dAtA[i] = 0xc8
		i++
		dAtA[i] = 0x1
		i++
		i = encodeVarintLifeline(dAtA, i, uint64(m.StateID))
	}
	if m.LatestUpdate != 0 {
		dAtA[i] = 0xd0
		i++
		dAtA[i] = 0x1
		i++
		i = encodeVarintLifeline(dAtA, i, uint64(m.LatestUpdate))
	}
	dAtA[i] = 0xda
	i++
	dAtA[i] = 0x1
	i++
	i = encodeVarintLifeline(dAtA, i, uint64(m.JetID.Size()))
	n5, err := m.JetID.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n5
	if m.LatestRequest != nil {
		dAtA[i] = 0xe2
		i++
		dAtA[i] = 0x1
		i++
		i = encodeVarintLifeline(dAtA, i, uint64(m.LatestRequest.Size()))
		n6, err := m.LatestRequest.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n6
	}
	return i, nil
}

func (m *LifelineDelegate) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *LifelineDelegate) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.XPolymorph != 0 {
		dAtA[i] = 0x80
		i++
		dAtA[i] = 0x1
		i++
		i = encodeVarintLifeline(dAtA, i, uint64(m.XPolymorph))
	}
	dAtA[i] = 0xa2
	i++
	dAtA[i] = 0x1
	i++
	i = encodeVarintLifeline(dAtA, i, uint64(m.Key.Size()))
	n7, err := m.Key.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n7
	dAtA[i] = 0xaa
	i++
	dAtA[i] = 0x1
	i++
	i = encodeVarintLifeline(dAtA, i, uint64(m.Value.Size()))
	n8, err := m.Value.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n8
	return i, nil
}

func encodeVarintLifeline(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *Lifeline) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.XPolymorph != 0 {
		n += 2 + sovLifeline(uint64(m.XPolymorph))
	}
	if m.LatestState != nil {
		l = m.LatestState.Size()
		n += 2 + l + sovLifeline(uint64(l))
	}
	if m.LatestStateApproved != nil {
		l = m.LatestStateApproved.Size()
		n += 2 + l + sovLifeline(uint64(l))
	}
	if m.ChildPointer != nil {
		l = m.ChildPointer.Size()
		n += 2 + l + sovLifeline(uint64(l))
	}
	l = m.Parent.Size()
	n += 2 + l + sovLifeline(uint64(l))
	if len(m.Delegates) > 0 {
		for _, e := range m.Delegates {
			l = e.Size()
			n += 2 + l + sovLifeline(uint64(l))
		}
	}
	if m.StateID != 0 {
		n += 2 + sovLifeline(uint64(m.StateID))
	}
	if m.LatestUpdate != 0 {
		n += 2 + sovLifeline(uint64(m.LatestUpdate))
	}
	l = m.JetID.Size()
	n += 2 + l + sovLifeline(uint64(l))
	if m.LatestRequest != nil {
		l = m.LatestRequest.Size()
		n += 2 + l + sovLifeline(uint64(l))
	}
	return n
}

func (m *LifelineDelegate) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.XPolymorph != 0 {
		n += 2 + sovLifeline(uint64(m.XPolymorph))
	}
	l = m.Key.Size()
	n += 2 + l + sovLifeline(uint64(l))
	l = m.Value.Size()
	n += 2 + l + sovLifeline(uint64(l))
	return n
}

func sovLifeline(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozLifeline(x uint64) (n int) {
	return sovLifeline(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *Lifeline) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Lifeline{`,
		`XPolymorph:` + fmt.Sprintf("%v", this.XPolymorph) + `,`,
		`LatestState:` + fmt.Sprintf("%v", this.LatestState) + `,`,
		`LatestStateApproved:` + fmt.Sprintf("%v", this.LatestStateApproved) + `,`,
		`ChildPointer:` + fmt.Sprintf("%v", this.ChildPointer) + `,`,
		`Parent:` + fmt.Sprintf("%v", this.Parent) + `,`,
		`Delegates:` + strings.Replace(strings.Replace(fmt.Sprintf("%v", this.Delegates), "LifelineDelegate", "LifelineDelegate", 1), `&`, ``, 1) + `,`,
		`StateID:` + fmt.Sprintf("%v", this.StateID) + `,`,
		`LatestUpdate:` + fmt.Sprintf("%v", this.LatestUpdate) + `,`,
		`JetID:` + fmt.Sprintf("%v", this.JetID) + `,`,
		`LatestRequest:` + fmt.Sprintf("%v", this.LatestRequest) + `,`,
		`}`,
	}, "")
	return s
}
func (this *LifelineDelegate) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&LifelineDelegate{`,
		`XPolymorph:` + fmt.Sprintf("%v", this.XPolymorph) + `,`,
		`Key:` + fmt.Sprintf("%v", this.Key) + `,`,
		`Value:` + fmt.Sprintf("%v", this.Value) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringLifeline(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *Lifeline) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowLifeline
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Lifeline: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Lifeline: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 16:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field XPolymorph", wireType)
			}
			m.XPolymorph = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLifeline
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.XPolymorph |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 20:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LatestState", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLifeline
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthLifeline
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthLifeline
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			var v github_com_insolar_insolar_insolar.ID
			m.LatestState = &v
			if err := m.LatestState.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 21:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LatestStateApproved", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLifeline
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthLifeline
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthLifeline
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			var v github_com_insolar_insolar_insolar.ID
			m.LatestStateApproved = &v
			if err := m.LatestStateApproved.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 22:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChildPointer", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLifeline
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthLifeline
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthLifeline
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			var v github_com_insolar_insolar_insolar.ID
			m.ChildPointer = &v
			if err := m.ChildPointer.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 23:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Parent", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLifeline
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthLifeline
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthLifeline
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Parent.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 24:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Delegates", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLifeline
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthLifeline
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthLifeline
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Delegates = append(m.Delegates, LifelineDelegate{})
			if err := m.Delegates[len(m.Delegates)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 25:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field StateID", wireType)
			}
			m.StateID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLifeline
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.StateID |= github_com_insolar_insolar_insolar_record.StateID(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 26:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field LatestUpdate", wireType)
			}
			m.LatestUpdate = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLifeline
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.LatestUpdate |= github_com_insolar_insolar_insolar.PulseNumber(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 27:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field JetID", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLifeline
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthLifeline
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthLifeline
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.JetID.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 28:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LatestRequest", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLifeline
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthLifeline
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthLifeline
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			var v github_com_insolar_insolar_insolar.ID
			m.LatestRequest = &v
			if err := m.LatestRequest.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipLifeline(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthLifeline
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthLifeline
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *LifelineDelegate) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowLifeline
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: LifelineDelegate: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: LifelineDelegate: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 16:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field XPolymorph", wireType)
			}
			m.XPolymorph = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLifeline
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.XPolymorph |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 20:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Key", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLifeline
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthLifeline
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthLifeline
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Key.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 21:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Value", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLifeline
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthLifeline
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthLifeline
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Value.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipLifeline(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthLifeline
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthLifeline
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipLifeline(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowLifeline
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowLifeline
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowLifeline
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthLifeline
			}
			iNdEx += length
			if iNdEx < 0 {
				return 0, ErrInvalidLengthLifeline
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowLifeline
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipLifeline(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
				if iNdEx < 0 {
					return 0, ErrInvalidLengthLifeline
				}
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthLifeline = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowLifeline   = fmt.Errorf("proto: integer overflow")
)
