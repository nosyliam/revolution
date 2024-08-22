package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfigObject_SetPrimitive(t *testing.T) {
	type object struct {
		Val string `yaml:"val"`
	}

	var obj = Object[object]{}
	obj.Initialize("")
	assert.NoError(t, obj.SetPath("val", "test"))
	val, err := obj.GetPath("val")
	assert.NoError(t, err)
	assert.Equal(t, val, "test")
}

func TestConfigObject_SetNestedObjectPrimitive(t *testing.T) {
	type nested2 struct {
		Val2 int  `yaml:"val2"`
		Val3 bool `state:"val3"`
	}

	type nested struct {
		Val string `yaml:"val"`
		Obj *Object[nested2]
	}

	type object struct {
		Val string          `yaml:"val"`
		Obj *Object[nested] `yaml:"obj"`
	}

	var obj = Object[object]{}
	obj.Initialize("")
	assert.NoError(t, obj.SetPath("obj.val", "test"))
	val, err := obj.GetPath("obj.val")
	assert.NoError(t, err)
	assert.Equal(t, val, "test")
}
