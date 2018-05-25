package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAttrBasicJSONUnmarashel(t *testing.T) {
	val := AttrAttrT{}

	val.UnmarshalJSON([]byte(`true`))

	assert.Equal(t, "true", val.GetAttr())

	val.UnmarshalJSON([]byte(`"S"`))

	assert.Equal(t, "S", val.GetAttr())
}

func TestAttrBasicSize(t *testing.T) {
	v := `"S"`
	val := AttrAttrT{Attr: v}

	assert.Equal(t, 5, val.Size())
}

func TestAttrBasicMarshalTo(t *testing.T) {
	v := `"S"`
	val := AttrAttrT{Attr: v}

	data := make([]byte, val.Size())
	size, err := val.MarshalTo(data)

	assert.Nil(t, err)
	assert.Len(t, data, 5)
	assert.Equal(t, 5, size)
}

func TestAttrBasicProtoMarshal(t *testing.T) {
	v := `"S"`
	val := AttrAttrT{Attr: v}

	data, err := val.Marshal()
	assert.Nil(t, err)
	assert.Len(t, data, 5)
}

func TestAttrBasicProtoRemarshal(t *testing.T) {
	v := `"S"`
	val := AttrAttrT{Attr: v}

	data, err := val.Marshal()
	assert.Nil(t, err)
	assert.Len(t, data, 5)

	val2 := AttrAttrT{}

	err = val2.Unmarshal(data)
	assert.Nil(t, err)

	assert.Equal(t, v, val2.GetAttr())
}
