package main

import (
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestBasicJSONUnmarashel(t *testing.T) {
	val := PropValueT{}

	val.UnmarshalJSON([]byte(`["35%",1]`))

	assert.Equal(t, val.GetValue(), "35%")
	assert.Equal(t, val.GetValueType(), int64(1))
}

func TestBasicProtoRemashalling(t *testing.T) {
	v := `11\/315313`
	vt := int64(1)
	val := PropValueT{Value: v, ValueType: vt}

	payload, err := proto.Marshal(val)
	assert.Nil(t, err)
	assert.Len(t, payload, 10)

	// assert.Equal(t, val.GetValue(), "35%")
	// assert.Equal(t, val.GetValueType(), int64(1))
}
