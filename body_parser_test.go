package m3lshttp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseBodyJson(t *testing.T) {
	testReq := new(RequestMock)
	testReq.On("ContentType").Return("application/json")
	testReq.On("Body").Return([]byte("{\"a\": \"b\"}"))
	output := parseBody(testReq)

	require.NotNil(t, output, "Params returned is nil")
	obj := output.Object()

	data, ok := obj.(map[string]interface{})
	require.True(t, ok, "Data not parsed to map")
	dataObject, ok := data["_data"]
	require.True(t, ok, "Data object not located")
	dataObjectMap, ok := dataObject.(map[string]interface{})
	require.True(t, ok, "Data object not map")
	aObject, ok := dataObjectMap["a"]
	require.True(t, ok, "Object not located")
	aObjectCasted, ok := aObject.(string)
	require.True(t, ok, "Object not string")
	assert.Equal(t, "b", aObjectCasted)
}

func TestParseBodyMultipart(t *testing.T) {
	testReq := new(RequestMock)
	testReq.On("ContentType").Return("multipart/form-data")
	testReq.On("MultipartForm").Return(map[string][]string{"a": []string{"b"}})
	output := parseBody(testReq)

	require.NotNil(t, output, "Params returned is nil")
	obj := output.Object()

	data, ok := obj.(map[string]interface{})
	require.True(t, ok, "Data not parsed to map")
	dataObject, ok := data["_data"]
	require.True(t, ok, "Data object not located")
	dataObjectMap, ok := dataObject.(map[string]interface{})
	require.True(t, ok, "Data object not map")
	aObject, ok := dataObjectMap["a"]
	require.True(t, ok, "Object not located")
	aFirst, ok := aObject.(string)
	require.True(t, ok, "Element not string")
	assert.Equal(t, "b", aFirst)
}

func TestParseBodyString(t *testing.T) {
	testReq := new(RequestMock)
	testReq.On("ContentType").Return("text/plain")
	testReq.On("Body").Return([]byte("7amada"))
	output := parseBody(testReq)

	require.NotNil(t, output, "Params returned is nil")
	obj := output.Object()

	data, ok := obj.(map[string]interface{})
	require.True(t, ok, "Data not parsed to map")
	dataObject, ok := data["_data"]
	require.True(t, ok, "Data object not located")
	dataObjectValue, ok := dataObject.(string)
	require.True(t, ok, "Data object not string")
	assert.Equal(t, "7amada", dataObjectValue)
}
