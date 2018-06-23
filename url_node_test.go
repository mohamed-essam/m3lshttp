package m3lshttp

import (
	"github.com/mohamed-essam/m3lsh"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewUrlNode(t *testing.T) {
	node := newUrlNode("partA")
	if node.Children == nil {
		t.Error("newUrlNode should initialize children array")
	}
	if node.handlers == nil {
		t.Error("newUrlNode should initialize handlers map")
	}
	if node.part != "partA" {
		t.Error("Wrong part string")
	}
	if node.isVariable {
		t.Error("Node considered variable when not")
	}

	nodeB := newUrlNode(":partB")
	if !nodeB.isVariable {
		t.Error("Node considered non-variable when it is")
	}
}

func TestIsVariable(t *testing.T) {
	if isVariable("abcd") {
		t.Error("abcd is not variadic")
	}

	if !isVariable(":abc") {
		t.Error(":abc is not variadic")
	}
}

func TestAddPathTerminalNode(t *testing.T) {
	var path []string
	var idx int
	var node *urlNode
	method := "POST"
	fn := func(r Request) {}

	node = newUrlNode("api")
	path = []string{"", "api"}
	idx = 2

	node.addPath(path, idx, method, fn)
	t.Log("node should have POST handler set")
	if node.handlers["POST"] == nil {
		t.Error("Handler not set")
	}
	t.Log("node should not create any children")
	if len(node.Children) > 0 {
		t.Error("Node created child when it should not")
	}

	t.Log("node should throw error on duplicate route")
	m3lsh.Try(func() {
		node.addPath(path, idx, method, fn)
		t.Error("Error not thrown")
	})
}

func TestAddPathIntermediate(t *testing.T) {
	var path []string
	var idx int
	var node *urlNode
	method := "POST"
	fn := func(r Request) {}

	t.Log("when node has no applicable children")

	node = newUrlNode("api")
	path = []string{"", "api"}
	idx = 1

	node.addPath(path, idx, method, fn)

	t.Log(">node should not set any handlers")
	if len(node.handlers) > 0 {
		t.Error("Node created handler for itself")
	}
	t.Log(">node should create one child")
	if len(node.Children) != 1 {
		t.Fatalf("Expected 1 child, found %d", len(node.Children))
	}
	t.Log(">child node should not have children")
	if len(node.Children[0].Children) > 0 {
		t.Errorf("Expected no children, found %d", len(node.Children[0].Children))
	}
	t.Log(">child node should have POST handler set")
	if node.Children[0].handlers["POST"] == nil {
		t.Error("Child handler not set")
	}

	t.Log("when node has an applicable child")

	path = []string{"", "api", "sdk"}
	method = "GET"

	node.addPath(path, idx, method, fn)

	t.Log(">child node should have a child of its own")
	if len(node.Children[0].Children) != 1 {
		t.Fatalf("Expected 1 child, found %d", len(node.Children[0].Children))
	}

	t.Log(">grandchild node should handle new GET path")
	if node.Children[0].Children[0].handlers["GET"] == nil {
		t.Error("Grandchild handler not set")
	}
}

func TestHandleExists(t *testing.T) {
	var path []string
	var idx int
	var node *urlNode
	method := "POST"
	handled := false
	fn := func(r Request) {
		assert.NotNil(t, r)
		handled = true
	}

	node = newUrlNode("api")
	path = []string{"", "api"}
	idx = 1

	node.addPath(path, idx, method, fn)

	req := &RequestMock{}
	req.On("ContentType").Return("application/json")
	req.On("Body").Return("{\"a\": \"b\"}")

	node.handle(path, idx, method, req)
	assert.True(t, handled, "Request not handled")
}

func TestHandleHead(t *testing.T) {
	var path []string
	var idx int
	var node *urlNode
	method := "POST"
	fn := func(r Request) {
		t.Error("HEAD request handled")
	}

	node = newUrlNode("api")
	path = []string{"", "api"}
	idx = 1

	node.addPath(path, idx, method, fn)

	req := &RequestMock{}
	req.On("ContentType").Return("application/json")
	req.On("Body").Return("{\"a\": \"b\"}")

	resp := node.handle(path, idx, "HEAD", req)

	assert.True(t, resp)
}

func TestHandleMethodNotAllowed(t *testing.T) {
	var path []string
	var idx int
	var node *urlNode
	method := "POST"
	fn := func(r Request) {}

	node = newUrlNode("api")
	path = []string{"", "api"}
	idx = 1

	node.addPath(path, idx, method, fn)

	req := &RequestMock{}
	req.On("ContentType").Return("application/json")
	req.On("Body").Return("{\"a\": \"b\"}")

	ex := m3lsh.Try(func() {
		node.handle(path, idx, "GET", req)
		t.Error("Error not thrown")
	})
	assert.NotNil(t, ex)
	assert.IsType(t, &MethodNotAllowed{}, ex)
}

func TestHandleVariable(t *testing.T) {
	var path []string
	var idx int
	var node *urlNode
	method := "GET"
	handled := false
	fn := func(r Request) {
		assert.NotNil(t, r)
		handled = true
	}

	node = newUrlNode("api")
	path = []string{"", "api", ":version"}
	callPath := []string{"", "api", "v3"}
	idx = 1

	node.addPath(path, idx, method, fn)

	req := &RequestMock{}
	req.On("ContentType").Return("application/json")
	req.On("Body").Return("{\"a\": \"b\"}")
	pushPathParamCalled := false
	popPathParamCalled := false
	req.On("pushPathParam", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		pushPathParamCalled = true
		assert.Equal(t, "version", args.Get(0))
		assert.Equal(t, "v3", args.Get(1))
	})
	req.On("popPathParam").Run(func(args mock.Arguments) {
		popPathParamCalled = true
	})
	params := &ParamsMock{}
	req.On("Params").Return(params)

	resp := node.handle(callPath, idx, method, req)

	assert.True(t, resp, "Returned unhandled")
	assert.True(t, handled, "Not handled")
	assert.True(t, pushPathParamCalled, "Variable not pushed")
	assert.True(t, pushPathParamCalled, "Variable not popped")
}

func TestHandleNotExists(t *testing.T) {
	var path []string
	var idx int
	var node *urlNode
	method := "GET"
	fn := func(r Request) {}

	node = newUrlNode("api")
	path = []string{"", "api"}
	callPath := []string{"", "admin"}
	idx = 1

	node.addPath(path, idx, method, fn)

	req := &RequestMock{}
	req.On("ContentType").Return("application/json")
	req.On("Body").Return("{\"a\": \"b\"}")

	resp := node.handle(callPath, idx, method, req)

	assert.False(t, resp)
}
