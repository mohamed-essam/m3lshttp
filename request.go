package m3lshttp

import "github.com/valyala/fasthttp"

type Request interface {
	pushPathParam(name, value string)
	popPathParam()
	Params() Params
	ContentType() string
	Body() []byte
	MultipartForm() map[string][]string
	Path() string
	Method() string
	context() *fasthttp.RequestCtx
}

type RequestWrapper struct {
	ctx                 *fasthttp.RequestCtx
	params              Params
	pathParams          []string
	pathValues          []string
	pathParamsEvaluated bool
}

func newRequest(ctx *fasthttp.RequestCtx) *RequestWrapper {
	req := &RequestWrapper{ctx: ctx, pathParams: make([]string, 0), pathValues: make([]string, 0)}
	req.params = parseBody(req)
	return req
}

func (r *RequestWrapper) pushPathParam(name, value string) {
	r.pathParams = append(r.pathParams, name)
	r.pathValues = append(r.pathValues, value)
}

func (r *RequestWrapper) popPathParam() {
	r.pathParams = r.pathParams[:len(r.pathParams)-1]
	r.pathValues = r.pathValues[:len(r.pathValues)-1]
}

func (r *RequestWrapper) Params() Params {
	if r.pathParamsEvaluated {
		return r.params
	}
	for idx, p := range r.pathParams {
		r.params.addObject(p, r.pathValues[idx])
	}
	return r.params
}

func (r *RequestWrapper) ContentType() string {
	return string(r.ctx.Request.Header.ContentType())
}

func (r *RequestWrapper) Body() []byte {
	return r.ctx.Request.Body()
}

func (r *RequestWrapper) MultipartForm() map[string][]string {
	form, _ := r.ctx.MultipartForm()
	return form.Value
}

func (r *RequestWrapper) Path() string {
	return string(r.ctx.Path())
}

func (r *RequestWrapper) Method() string {
	return string(r.ctx.Method())
}

func (r *RequestWrapper) context() *fasthttp.RequestCtx {
	return r.ctx
}
