package xreflect_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/todennus/x/xreflect"
)

type Custom struct {
	A int    `custom:"a"`
	B string `custom:"b"`
}

func Test_ToMap_Value(t *testing.T) {
	c := Custom{A: 1, B: "2"}
	m := xreflect.ToMap(c, "custom")
	assert.Equal(t, 1, m["a"])
	assert.Equal(t, "2", m["b"])
}

func Test_ToMap_Pointer(t *testing.T) {
	c := Custom{A: 1, B: "2"}
	m := xreflect.ToMap(&c, "custom")
	assert.Equal(t, 1, m["a"])
	assert.Equal(t, "2", m["b"])
}
