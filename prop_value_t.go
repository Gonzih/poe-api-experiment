package main

import (
	"encoding/json"

	"github.com/gogo/protobuf/proto"
)

type PropValueT struct {
	Value     string `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	ValueType int64  `protobuf:"varint,3,opt,name=valueType,proto3" json:"valueType,omitempty"`
}

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

func (m PropValueT) Reset()         { m = PropValueT{} }
func (m PropValueT) String() string { return proto.CompactTextString(m) }
func (PropValueT) ProtoMessage()    {}

func (t PropValueT) Marshal() ([]byte, error) {
	return proto.Marshal(t)
}

func (t *PropValueT) MarshalTo(data []byte) (n int, err error) {
	return 0, nil
}

func (t *PropValueT) Unmarshal(data []byte) error {
	return proto.Unmarshal(data, t)
}

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

func (t *PropValueT) Size() int {
	return proto.Size(t)
}
