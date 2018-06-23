package m3lshttp

import (
	"fmt"
	"strconv"

	"github.com/mohamed-essam/m3lsh"
)

type Params interface {
	Object() interface{}
	GetObject(name string) Params
	GetArray(idx int) Params
	StringValue() string
	AsString() string
	Integer() int
	Long() int64
	Float() float32
	Double() float64
	AsInteger() int
	AsLong() int64
	AsFloat() float32
	AsDouble() float64
	addObject(key, value string)
}

type ParamsWrapper struct {
	object interface{}
}

type InvalidTypeException struct {
	m3lsh.BaseException
	Object interface{}
}

func newParams(object interface{}) Params {
	return &ParamsWrapper{object: cleanObject(object)}
}

func cleanObject(object interface{}) interface{} {
	return map[string]interface{}{"_data": object}
}

func (p ParamsWrapper) Object() interface{} {
	return p.object
}

func (p ParamsWrapper) GetObject(name string) Params {
	mp, ok := p.object.(map[string]interface{})
	if !ok {
		m3lsh.Throw(&InvalidTypeException{Object: p.object}, "Object is not a map")
	}
	return &ParamsWrapper{object: mp[name]}
}

func (p ParamsWrapper) GetArray(idx int) Params {
	ar, ok := p.object.([]interface{})
	if !ok {
		m3lsh.Throw(&InvalidTypeException{Object: p.object}, "Object is not an array")
	}
	return &ParamsWrapper{object: ar[idx]}
}

func (p ParamsWrapper) StringValue() string {
	str, ok := p.object.(string)
	if !ok {
		m3lsh.Throw(&InvalidTypeException{Object: p.object}, "Object is not a string")
	}
	return str
}

func (p ParamsWrapper) AsString() string {
	return fmt.Sprintf("%v", p.object)
}

func (p ParamsWrapper) Integer() int {
	num, ok := p.object.(int)
	if !ok {
		m3lsh.Throw(&InvalidTypeException{Object: p.object}, "Object is not an int")
	}
	return num
}

func (p ParamsWrapper) Long() int64 {
	num, ok := p.object.(int64)
	if !ok {
		m3lsh.Throw(&InvalidTypeException{Object: p.object}, "Object is not an int64")
	}
	return num
}

func (p ParamsWrapper) Float() float32 {
	num, ok := p.object.(float32)
	if !ok {
		m3lsh.Throw(&InvalidTypeException{Object: p.object}, "Object is not a float32")
	}
	return num
}

func (p ParamsWrapper) Double() float64 {
	num, ok := p.object.(float64)
	if !ok {
		m3lsh.Throw(&InvalidTypeException{Object: p.object}, "Object is not a float64")
	}
	return num
}

func (p ParamsWrapper) AsInteger() int {
	switch p.object.(type) {
	case int:
		return p.Integer()
	case int64:
		return int(p.Long())
	case float32:
		return int(p.Float())
	case float64:
		return int(p.Double())
	case string:
		conv, err := strconv.ParseInt(p.StringValue(), 10, 32)
		if err != nil {
			m3lsh.Throw(&InvalidTypeException{Object: p.object}, "String cannot be converted to int")
		}
		return int(conv)
	}
	m3lsh.Throw(&InvalidTypeException{Object: p.object}, "Object cannot be converted to int")
	return 0
}

func (p ParamsWrapper) AsLong() int64 {
	switch p.object.(type) {
	case int:
		return int64(p.Integer())
	case int64:
		return p.Long()
	case float32:
		return int64(p.Float())
	case float64:
		return int64(p.Double())
	case string:
		conv, err := strconv.ParseInt(p.StringValue(), 10, 64)
		if err != nil {
			m3lsh.Throw(&InvalidTypeException{Object: p.object}, "String cannot be converted to int64")
		}
		return conv
	}
	m3lsh.Throw(&InvalidTypeException{Object: p.object}, "Object cannot be converted to int64")
	return 0
}

func (p ParamsWrapper) AsFloat() float32 {
	switch p.object.(type) {
	case int:
		return float32(p.Integer())
	case int64:
		return float32(p.Long())
	case float32:
		return p.Float()
	case float64:
		return float32(p.Double())
	case string:
		conv, err := strconv.ParseFloat(p.StringValue(), 64)
		if err != nil {
			m3lsh.Throw(&InvalidTypeException{Object: p.object}, "String cannot be converted to int64")
		}
		return float32(conv)
	}
	m3lsh.Throw(&InvalidTypeException{Object: p.object}, "Object cannot be converted to int64")
	return 0
}

func (p ParamsWrapper) AsDouble() float64 {
	switch p.object.(type) {
	case int:
		return float64(p.Integer())
	case int64:
		return float64(p.Long())
	case float32:
		return float64(p.Float())
	case float64:
		return p.Double()
	case string:
		conv, err := strconv.ParseFloat(p.StringValue(), 64)
		if err != nil {
			m3lsh.Throw(&InvalidTypeException{Object: p.object}, "String cannot be converted to int64")
		}
		return float64(conv)
	}
	m3lsh.Throw(&InvalidTypeException{Object: p.object}, "Object cannot be converted to int64")
	return 0
}

func (p *ParamsWrapper) addObject(key, value string) {
	mp, ok := p.object.(map[string]interface{})
	if !ok {
		return
	}
	mp[key] = value
}
