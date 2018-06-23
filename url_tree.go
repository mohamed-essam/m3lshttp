package m3lshttp

import (
	"fmt"
	"github.com/mohamed-essam/m3lsh"
	"strings"
)

type urlTree struct {
	Children []*urlNode
}

type DuplicateRouteException struct {
	m3lsh.BaseException
	pathParts []string
	path      string
}

func newUrlTree() *urlTree {
	return &urlTree{Children: make([]*urlNode, 0)}
}

func (t *urlTree) addPath(path, method string, fn handler) {
	pathParts := strings.Split(path, "/")
	// should only find one path with exact same name
	part := findPart(t.Children, pathParts[0])
	if part == nil {
		// new part, create it
		part = newUrlNode(pathParts[0])
		t.Children = append(t.Children, part)
	}
	m3lsh.TryCatch(func() {
		part.addPath(pathParts, 1, method, fn)
	}, m3lsh.Catcher(&DuplicateRouteException{}, func(e interface{}) {
		ex := e.(*DuplicateRouteException)
		ex.path = path
		m3lsh.Throw(ex, fmt.Sprintf("Path %s is duplicated", path))
	}))
}

func (t urlTree) handle(req Request) {
	path := req.Path()
	method := req.Method()
	pathParts := strings.Split(string(path), "/")
	parts := findParts(t.Children, pathParts[0])
	for _, part := range parts {
		if part.handle(pathParts, 1, string(method), req) {
			return
		}
	}
	m3lsh.Throw(&NotFound{}, "")
}

func findParts(nodes []*urlNode, part string) []*urlNode {
	parts := make([]*urlNode, 0)
	for _, v := range nodes {
		if v.part == part || v.isVariable {
			parts = append(parts, v)
		}
	}
	return parts
}

func findPart(nodes []*urlNode, part string) *urlNode {
	for _, v := range nodes {
		if v.part == part {
			return v
		}
	}
	return nil
}
