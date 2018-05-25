package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValueBasicJSONUnmarashel(t *testing.T) {
	val := BoolStringT{}

	val.UnmarshalJSON([]byte(`true`))

	assert.Equal(t, "true", val.GetValue())

	val.UnmarshalJSON([]byte(`"S"`))

	assert.Equal(t, "S", val.GetValue())
}

func TestValueBasicSize(t *testing.T) {
	v := `"S"`
	val := BoolStringT{Value: v}

	assert.Equal(t, 5, val.Size())
}

func TestValueBasicMarshalTo(t *testing.T) {
	v := `"S"`
	val := BoolStringT{Value: v}

	data := make([]byte, val.Size())
	size, err := val.MarshalTo(data)

	assert.Nil(t, err)
	assert.Len(t, data, 5)
	assert.Equal(t, 5, size)
}

func TestValueBasicProtoMarshal(t *testing.T) {
	v := `"S"`
	val := BoolStringT{Value: v}

	data, err := val.Marshal()
	assert.Nil(t, err)
	assert.Len(t, data, 5)
}

func TestValueBasicProtoRemarshal(t *testing.T) {
	v := `"S"`
	val := BoolStringT{Value: v}

	data, err := val.Marshal()
	assert.Nil(t, err)
	assert.Len(t, data, 5)

	val2 := BoolStringT{}

	err = val2.Unmarshal(data)
	assert.Nil(t, err)

	assert.Equal(t, v, val2.GetValue())
}
