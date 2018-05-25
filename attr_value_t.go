package main

import (
	"encoding/json"
	fmt "fmt"
	io "io"
	"log"

	"github.com/gogo/protobuf/proto"
)

type AttrAttrT struct {
	Attr string `protobuf:"bytes,1,opt,name=attr,proto3" json:"attr,omitempty"`
}

func (m *AttrAttrT) Reset()                    { *m = AttrAttrT{} }
func (m *AttrAttrT) String() string            { return proto.CompactTextString(m) }
func (*AttrAttrT) ProtoMessage()               {}
func (*AttrAttrT) Descriptor() ([]byte, []int) { return fileDescriptorResponse, []int{5} }

func (m *AttrAttrT) GetAttr() string {
	if m != nil {
		return m.Attr
	}
	return ""
}

func (m *AttrAttrT) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AttrAttrT) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Attr) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintResponse(dAtA, i, uint64(len(m.Attr)))
		i += copy(dAtA[i:], m.Attr)
	}
	return i, nil
}

func (m *AttrAttrT) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowResponse
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: AttrAttrT: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AttrAttrT: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Attr", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowResponse
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthResponse
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Attr = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipResponse(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthResponse
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

func (m *AttrAttrT) Size() (n int) {
	var l int
	_ = l
	l = len(m.Attr)
	if l > 0 {
		n += 1 + l + sovResponse(uint64(l))
	}
	return n
}

// ======== Manually implemented ============== //

func (t AttrAttrT) MarshalJSON() ([]byte, error) {
	return json.Marshal(&t)
}

func (t *AttrAttrT) UnmarshalJSON(data []byte) error {
	var v interface{}
	var s string
	err := json.Unmarshal(data, &v)

	if err != nil {
		return err
	}

	switch t := v.(type) {
	case bool:
		s = fmt.Sprintf("%t", v.(bool))
	case string:
		s = v.(string)
	default:
		log.Fatalf("Unknown type %v", t)
	}

	// log.Printf("%s -> %#v ", string(data), s)
	t.Attr = s

	return nil
}
