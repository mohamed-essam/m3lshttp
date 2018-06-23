package m3lshttp

import (
	"testing"

	"github.com/mohamed-essam/m3lsh"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUrlTree(t *testing.T) {
	tree := newUrlTree()
	if tree.Children == nil {
		t.Error("Tree children not initialized")
	}
}

func TestFindPartExists(t *testing.T) {
	node1 := newUrlNode("abc")
	node2 := newUrlNode("abcd")
	node3 := newUrlNode("abcde")
	parts := []*urlNode{node1, node2, node3}
	foundPart := findPart(parts, "abc")

	if foundPart != node1 {
		t.Errorf("Part abc not found, instead found %+v", foundPart)
	}
}

func TestFindPartNotExists(t *testing.T) {
	node1 := newUrlNode("abc")
	node2 := newUrlNode("abcd")
	node3 := newUrlNode("abcde")
	parts := []*urlNode{node1, node2, node3}
	foundPart := findPart(parts, "cba")

	if foundPart != nil {
		t.Errorf("Non-existent part found: %+v", foundPart)
	}
}

func TestFindParts(t *testing.T) {
	node1 := newUrlNode("abc")
	node2 := newUrlNode("abcd")
	node3 := newUrlNode(":id")
	parts := []*urlNode{node1, node2, node3}
	foundParts := findParts(parts, "abc")

	if !urlNodeArrayContains(foundParts, node1) {
		t.Errorf("Part abc not found, instead found %+v", foundParts)
	}
	if !urlNodeArrayContains(foundParts, node3) {
		t.Errorf("Part :id not found, instead found %+v", foundParts)
	}
}

func TestAddPath(t *testing.T) {
	tree := newUrlTree()
	path := "/api/sdk"
	method := "POST"
	fn := func(r Request) {}
	tree.addPath(path, method, fn)
	require.True(t, len(tree.Children) > 0)
	require.True(t, len(tree.Children[0].Children) > 0)
	require.True(t, len(tree.Children[0].Children[0].Children) > 0)
	require.True(t, len(tree.Children[0].Children[0].Children[0].handlers) > 0)
	assert.Equal(t, "", tree.Children[0].part)
	assert.Equal(t, "api", tree.Children[0].Children[0].part)
	assert.Equal(t, "sdk", tree.Children[0].Children[0].Children[0].part)
}

func TestAddPathDuplicateRoute(t *testing.T) {
	tree := newUrlTree()
	path := "/api/sdk"
	method := "POST"
	fn := func(r Request) {}
	tree.addPath(path, method, fn)
	m3lsh.Try(func() {
		tree.addPath(path, method, fn)
		t.Error("Duplicate route not reported")
	})
}

func TestHandle(t *testing.T) {
	tree := newUrlTree()
	path := "/api/sdk"
	method := "POST"
	handled := false
	fn := func(r Request) {
		handled = true
	}
	tree.addPath(path, method, fn)
	req := &RequestMock{}
	req.On("Path").Return(path)
	req.On("Method").Return(method)
	tree.handle(req)
	assert.True(t, handled)
}

func TestHandleNotFound(t *testing.T) {
	tree := newUrlTree()
	path := "/api/sdk"
	method := "POST"
	handled := false
	fn := func(r Request) {
		handled = true
	}
	tree.addPath(path, method, fn)
	req := &RequestMock{}
	req.On("Path").Return("/api/web")
	req.On("Method").Return(method)
	ex := m3lsh.Try(func() {
		tree.handle(req)
		t.Error("Did not throw not found")
	})
	assert.False(t, handled)
	assert.NotNil(t, ex)
}
