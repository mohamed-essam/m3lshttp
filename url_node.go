package m3lshttp

import (
	"strings"

	"github.com/mohamed-essam/m3lsh"
)

type urlNode struct {
	Children   []*urlNode
	part       string
	isVariable bool
	handlers   map[string]handler
}

func newUrlNode(pathPart string) *urlNode {
	return &urlNode{Children: make([]*urlNode, 0), part: pathPart, isVariable: isVariable(pathPart), handlers: make(map[string]handler, 0)}
}

func (t *urlNode) addPath(path []string, idx int, method string, fn handler) {
	if idx >= len(path) {
		if handlerMapHasKey(t.handlers, method) {
			m3lsh.Throw(&DuplicateRouteException{pathParts: path}, "")
		}
		t.handlers[method] = fn
		return
	}
	applicable := findPart(t.Children, path[idx])

	if applicable == nil { // nothing applicable, create node
		part := newUrlNode(path[idx])
		t.Children = append(t.Children, part)
		part.addPath(path, idx+1, method, fn)
		return
	}
	applicable.addPath(path, idx+1, method, fn)
}

func (t *urlNode) handle(path []string, idx int, method string, req Request) bool {
	if t.isVariable {
		req.pushPathParam(t.part[1:], path[idx-1])
	}

	defer func() {
		if t.isVariable {
			req.popPathParam()
		}
	}()

	if idx >= len(path) {
		if method == "HEAD" || method == "OPTIONS" {
			return true
		}
		if !handlerMapHasKey(t.handlers, method) && len(t.handlers) > 0 {
			m3lsh.Throw(&MethodNotAllowed{}, "")
		} else if !handlerMapHasKey(t.handlers, method) {
			m3lsh.Throw(&NotFound{}, "")
		}
		t.handlers[method](req)
		return true
	}
	parts := findParts(t.Children, path[idx])

	for _, part := range parts {
		if part.handle(path, idx+1, method, req) {
			return true
		}
	}

	return false
}

func isVariable(pathPart string) bool {
	return strings.HasPrefix(pathPart, ":")
}
