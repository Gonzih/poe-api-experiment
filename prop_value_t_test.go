package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPropBasicJSONUnmarashel(t *testing.T) {
	val := PropValueT{}

	val.UnmarshalJSON([]byte(`["35%",1]`))

	assert.Equal(t, "35%", val.GetValue())
	assert.Equal(t, int64(1), val.GetValueType())
}

func TestPropBasicSize(t *testing.T) {
	v := `11\/315313`
	vt := int64(1)
	val := PropValueT{Value: v, ValueType: vt}

	assert.Equal(t, 14, val.Size())
}

func TestPropBasicMarshalTo(t *testing.T) {
	v := `11\/315313oe`
	vt := int64(1)
	val := PropValueT{Value: v, ValueType: vt}

	data := make([]byte, val.Size())
	size, err := val.MarshalTo(data)

	assert.Nil(t, err)
	assert.Len(t, data, 16)
	assert.Equal(t, 16, size)
}

func TestPropBasicProtoMarshal(t *testing.T) {
	v := `11\/315313e`
	vt := int64(1)
	val := PropValueT{Value: v, ValueType: vt}

	data, err := val.Marshal()
	assert.Nil(t, err)
	assert.Len(t, data, 15)
}

func TestPropBasicProtoRemarshal(t *testing.T) {
	v := `11\/315313eooeo`
	vt := int64(1)
	val := PropValueT{Value: v, ValueType: vt}

	data, err := val.Marshal()
	assert.Nil(t, err)
	assert.Len(t, data, 19)

	val2 := PropValueT{}

	err = val2.Unmarshal(data)
	assert.Nil(t, err)

	assert.Equal(t, v, val2.GetValue())
	assert.Equal(t, vt, val2.GetValueType())
}
