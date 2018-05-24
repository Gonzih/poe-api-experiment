package main

import (
	"encoding/json"
	fmt "fmt"
	io "io"

	"github.com/gogo/protobuf/proto"
)

type PropValueT struct {
	Value     string `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	ValueType int64  `protobuf:"varint,2,opt,name=valueType,proto3" json:"valueType,omitempty"`
}

func (m PropValueT) Reset()                     { m = PropValueT{} }
func (m PropValueT) String() string             { return proto.CompactTextString(m) }
func (PropValueT) ProtoMessage()                {}
func (*PropValueT) Descriptor() ([]byte, []int) { return fileDescriptorResponse, []int{5} }

func (m *PropValueT) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func (m *PropValueT) GetValueType() int64 {
	if m != nil {
		return m.ValueType
	}
	return 0
}

func (m *PropValueT) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PropValueT) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Value) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintResponse(dAtA, i, uint64(len(m.Value)))
		i += copy(dAtA[i:], m.Value)
	}
	if m.ValueType != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintResponse(dAtA, i, uint64(m.ValueType))
	}
	return i, nil
}

func (m PropValueT) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: PropValueT: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PropValueT: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Value", wireType)
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
			m.Value = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ValueType", wireType)
			}
			m.ValueType = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowResponse
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ValueType |= (int64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
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

func (m *PropValueT) Size() (n int) {
	var l int
	_ = l
	l = len(m.Value)
	if l > 0 {
		n += 1 + l + sovResponse(uint64(l))
	}
	if m.ValueType != 0 {
		n += 1 + sovResponse(uint64(m.ValueType))
	}
	return n
}

// ======== Manually implemented ============== //

func (t PropValueT) MarshalJSON() ([]byte, error) {
	return json.Marshal(&t)
}
func (t *PropValueT) UnmarshalJSON(data []byte) error {
	arr := make([]interface{}, 2)
	err := json.Unmarshal(data, &arr)

	if err != nil {
		return err
	}

	// log.Printf("%s -> %#v ", string(data), arr)
	t.Value = arr[0].(string)
	t.ValueType = int64(arr[1].(float64))

	return nil
}
