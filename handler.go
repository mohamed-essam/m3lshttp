package m3lshttp

import (
	"encoding/json"

	"github.com/mohamed-essam/m3lsh"

	"github.com/valyala/fasthttp"
)

type HttpHandler struct {
	tree *urlTree
}

type handler func(Request)

func NewHttpHandler() *HttpHandler {
	return &HttpHandler{tree: newUrlTree()}
}

func (h *HttpHandler) POST(path string, fn handler) {
	h.tree.addPath(path, "POST", fn)
}

func (h *HttpHandler) GET(path string, fn handler) {
	h.tree.addPath(path, "GET", fn)
}

func (h *HttpHandler) PUT(path string, fn handler) {
	h.tree.addPath(path, "PUT", fn)
}

func (h *HttpHandler) DELETE(path string, fn handler) {
	h.tree.addPath(path, "DELETE", fn)
}

func (h *HttpHandler) PATCH(path string, fn handler) {
	h.tree.addPath(path, "PATCH", fn)
}

func (h HttpHandler) handle(ctx *fasthttp.RequestCtx) {
	m3lsh.TryCatch(func() {
		h.tree.handle(newRequest(ctx))
	}, m3lsh.Catcher(&BadRequest{}, func(e interface{}) {
		ex := e.(*BadRequest)
		ctx.Response.SetStatusCode(400)
		ctx.Response.SetBody([]byte(ex.Message))
	}), m3lsh.Catcher(&Unauthorized{}, func(e interface{}) {
		ex := e.(*Unauthorized)
		ctx.Response.SetStatusCode(401)
		ctx.Response.SetBody([]byte(ex.Message))
	}), m3lsh.Catcher(&PaymentRequired{}, func(e interface{}) {
		ex := e.(*PaymentRequired)
		ctx.Response.SetStatusCode(402)
		ctx.Response.SetBody([]byte(ex.Message))
	}), m3lsh.Catcher(&Forbidden{}, func(e interface{}) {
		ex := e.(*Forbidden)
		ctx.Response.SetStatusCode(403)
		ctx.Response.SetBody([]byte(ex.Message))
	}), m3lsh.Catcher(&NotFound{}, func(e interface{}) {
		ex := e.(*NotFound)
		ctx.Response.SetStatusCode(404)
		ctx.Response.SetBody([]byte(ex.Message))
	}), m3lsh.Catcher(&MethodNotAllowed{}, func(e interface{}) {
		ex := e.(*MethodNotAllowed)
		ctx.Response.SetStatusCode(405)
		ctx.Response.SetBody([]byte(ex.Message))
	}), m3lsh.Catcher(&TimedOut{}, func(e interface{}) {
		ex := e.(*TimedOut)
		ctx.Response.SetStatusCode(408)
		ctx.Response.SetBody([]byte(ex.Message))
	}), m3lsh.Catcher(&UnprocessableEntity{}, func(e interface{}) {
		ex := e.(*UnprocessableEntity)
		ctx.Response.SetStatusCode(422)
		ctx.Response.SetBody([]byte(ex.Message))
	}), m3lsh.Catcher(&InternalServerError{}, func(e interface{}) {
		ex := e.(*InternalServerError)
		ctx.Response.SetStatusCode(500)
		ctx.Response.SetBody([]byte(ex.Message))
	}))
}

func (h HttpHandler) ListenAndServe(port string) error {
	return fasthttp.ListenAndServe(port, h.handle)
}

type ResponseType int

const (
	Json = iota
)

func Respond(r Request, responseType ResponseType, data interface{}) {
	var writtenData []byte
	switch responseType {
	case Json:
		writtenData, _ = json.Marshal(data)
	}
	r.context().Response.SetBody(writtenData)
}
