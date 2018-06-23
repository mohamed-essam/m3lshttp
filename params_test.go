package m3lshttp

import (
	"testing"

	"github.com/mohamed-essam/m3lsh"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewParams(t *testing.T) {
	p := newParams(map[string]interface{}{"a": "b"})
	require.IsType(t, map[string]interface{}{}, p.Object())
	obj, _ := p.Object().(map[string]interface{})
	assert.NotNil(t, obj["_data"])
	assert.IsType(t, map[string]interface{}{}, obj["_data"])
	data, _ := obj["_data"].(map[string]interface{})
	assert.NotNil(t, data["a"])
	assert.IsType(t, "", data["a"])
}

func TestGetObject(t *testing.T) {
	p := newParams(map[string]interface{}{"a": "b"})
	obj := p.GetObject("_data").GetObject("a").StringValue()
	assert.NotNil(t, obj)
	assert.Equal(t, "b", obj)
}

func TestGetObjectWrongValue(t *testing.T) {
	p := newParams([]interface{}{"a", "b"})
	ex := m3lsh.Try(func() {
		p.GetObject("_data").GetObject("a")
		t.Error("Not panicked")
	})
	assert.NotNil(t, ex)
}

func TestGetArray(t *testing.T) {
	p := newParams([]interface{}{"a", "b"})
	obj := p.GetObject("_data").GetArray(0).StringValue()
	assert.NotNil(t, obj)
	assert.Equal(t, "a", obj)
}

func TestGetArrayWrongValue(t *testing.T) {
	p := newParams(map[string]interface{}{"a": "b"})
	ex := m3lsh.Try(func() {
		p.GetObject("_data").GetArray(0)
		t.Error("Not panicked")
	})
	assert.NotNil(t, ex)
}

func TestStringValue(t *testing.T) {
	p := newParams("abc")
	obj := p.GetObject("_data").StringValue()
	assert.NotNil(t, obj)
	assert.Equal(t, "abc", obj)
}

func TestStringValueWrong(t *testing.T) {
	p := newParams(5)
	ex := m3lsh.Try(func() {
		p.GetObject("_data").StringValue()
		t.Error("Not panicked")
	})
	assert.NotNil(t, ex)
}

func TestInteger(t *testing.T) {
	p := newParams(5)
	obj := p.GetObject("_data").Integer()
	assert.NotNil(t, obj)
	assert.Equal(t, 5, obj)
}

func TestIntegerWrong(t *testing.T) {
	p := newParams("abc")
	ex := m3lsh.Try(func() {
		p.GetObject("_data").Integer()
		t.Error("Not panicked")
	})
	assert.NotNil(t, ex)
}

func TestFloat(t *testing.T) {
	p := newParams(float32(5.0))
	obj := p.GetObject("_data").Float()
	assert.NotNil(t, obj)
	assert.Equal(t, float32(5.0), obj)
}

func TestFloatWrong(t *testing.T) {
	p := newParams("abc")
	ex := m3lsh.Try(func() {
		p.GetObject("_data").Float()
		t.Error("Not panicked")
	})
	assert.NotNil(t, ex)
}

func TestDouble(t *testing.T) {
	p := newParams(float64(5.0))
	obj := p.GetObject("_data").Double()
	assert.NotNil(t, obj)
	assert.Equal(t, float64(5.0), obj)
}

func TestDoubleWrong(t *testing.T) {
	p := newParams("abc")
	ex := m3lsh.Try(func() {
		p.GetObject("_data").Double()
		t.Error("Not panicked")
	})
	assert.NotNil(t, ex)
}

func TestLong(t *testing.T) {
	p := newParams(int64(5))
	obj := p.GetObject("_data").Long()
	assert.NotNil(t, obj)
	assert.Equal(t, int64(5), obj)
}

func TestLongWrong(t *testing.T) {
	p := newParams("abc")
	ex := m3lsh.Try(func() {
		p.GetObject("_data").Long()
		t.Error("Not panicked")
	})
	assert.NotNil(t, ex)
}

func TestAsString(t *testing.T) {
	p := newParams(5)
	obj := p.GetObject("_data").AsString()
	assert.NotNil(t, obj)
	assert.Equal(t, "5", obj)
}

func TestAsIntegerFromString(t *testing.T) {
	p := newParams("5")
	obj := p.GetObject("_data").AsInteger()
	assert.NotNil(t, obj)
	assert.Equal(t, 5, obj)
}

func TestAsIntegerFromStringWrong(t *testing.T) {
	p := newParams("abc")
	ex := m3lsh.Try(func() {
		p.GetObject("_data").AsInteger()
		t.Error("Not panicked")
	})
	assert.NotNil(t, ex)
}

func TestAsIntegerFromInteger(t *testing.T) {
	p := newParams(5)
	obj := p.GetObject("_data").AsInteger()
	assert.NotNil(t, obj)
	assert.Equal(t, 5, obj)
}

func TestAsIntegerFromLong(t *testing.T) {
	p := newParams(int64(5))
	obj := p.GetObject("_data").AsInteger()
	assert.NotNil(t, obj)
	assert.Equal(t, 5, obj)
}

func TestAsIntegerFromFloat(t *testing.T) {
	p := newParams(float32(5.3))
	obj := p.GetObject("_data").AsInteger()
	assert.NotNil(t, obj)
	assert.Equal(t, 5, obj)
}

func TestAsIntegerFromDouble(t *testing.T) {
	p := newParams(float64(5.3))
	obj := p.GetObject("_data").AsInteger()
	assert.NotNil(t, obj)
	assert.Equal(t, 5, obj)
}

func TestAsLongFromString(t *testing.T) {
	p := newParams("5")
	obj := p.GetObject("_data").AsLong()
	assert.NotNil(t, obj)
	assert.Equal(t, int64(5), obj)
}

func TestAsLongFromStringWrong(t *testing.T) {
	p := newParams("abc")
	ex := m3lsh.Try(func() {
		p.GetObject("_data").AsLong()
		t.Error("Not panicked")
	})
	assert.NotNil(t, ex)
}

func TestAsLongFromInteger(t *testing.T) {
	p := newParams(5)
	obj := p.GetObject("_data").AsLong()
	assert.NotNil(t, obj)
	assert.Equal(t, int64(5), obj)
}

func TestAsLongFromLong(t *testing.T) {
	p := newParams(int64(5))
	obj := p.GetObject("_data").AsLong()
	assert.NotNil(t, obj)
	assert.Equal(t, int64(5), obj)
}

func TestAsLongFromFloat(t *testing.T) {
	p := newParams(float32(5.3))
	obj := p.GetObject("_data").AsLong()
	assert.NotNil(t, obj)
	assert.Equal(t, int64(5), obj)
}

func TestAsLongFromDouble(t *testing.T) {
	p := newParams(float64(5.3))
	obj := p.GetObject("_data").AsLong()
	assert.NotNil(t, obj)
	assert.Equal(t, int64(5), obj)
}

func TestAsFloatFromString(t *testing.T) {
	p := newParams("5")
	obj := p.GetObject("_data").AsFloat()
	assert.NotNil(t, obj)
	assert.Equal(t, float32(5), obj)
}

func TestAsFloatFromStringWrong(t *testing.T) {
	p := newParams("abc")
	ex := m3lsh.Try(func() {
		p.GetObject("_data").AsFloat()
		t.Error("Not panicked")
	})
	assert.NotNil(t, ex)
}

func TestAsFloatFromInteger(t *testing.T) {
	p := newParams(5)
	obj := p.GetObject("_data").AsFloat()
	assert.NotNil(t, obj)
	assert.Equal(t, float32(5), obj)
}

func TestAsFloatFromLong(t *testing.T) {
	p := newParams(int64(5))
	obj := p.GetObject("_data").AsFloat()
	assert.NotNil(t, obj)
	assert.Equal(t, float32(5), obj)
}

func TestAsFloatFromFloat(t *testing.T) {
	p := newParams(float32(5.3))
	obj := p.GetObject("_data").AsFloat()
	assert.NotNil(t, obj)
	assert.Equal(t, float32(5.3), obj)
}

func TestAsFloatFromDouble(t *testing.T) {
	p := newParams(float64(5.3))
	obj := p.GetObject("_data").AsFloat()
	assert.NotNil(t, obj)
	assert.Equal(t, float32(5.3), obj)
}

func TestAsDoubleFromString(t *testing.T) {
	p := newParams("5")
	obj := p.GetObject("_data").AsDouble()
	assert.NotNil(t, obj)
	assert.Equal(t, float64(5), obj)
}

func TestAsDoubleFromStringWrong(t *testing.T) {
	p := newParams("abc")
	ex := m3lsh.Try(func() {
		p.GetObject("_data").AsDouble()
		t.Error("Not panicked")
	})
	assert.NotNil(t, ex)
}

func TestAsDoubleFromInteger(t *testing.T) {
	p := newParams(5)
	obj := p.GetObject("_data").AsDouble()
	assert.NotNil(t, obj)
	assert.Equal(t, float64(5), obj)
}

func TestAsDoubleFromLong(t *testing.T) {
	p := newParams(int64(5))
	obj := p.GetObject("_data").AsDouble()
	assert.NotNil(t, obj)
	assert.Equal(t, float64(5), obj)
}

func TestAsDoubleFromFloat(t *testing.T) {
	p := newParams(float32(5.3))
	obj := p.GetObject("_data").AsDouble()
	assert.NotNil(t, obj)
	assert.Equal(t, float64(float32(5.3)), obj)
}

func TestAsDoubleFromDouble(t *testing.T) {
	p := newParams(float64(5.3))
	obj := p.GetObject("_data").AsDouble()
	assert.NotNil(t, obj)
	assert.Equal(t, float64(5.3), obj)
}
