package m3lshttp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
)

func tempServer(handler func(*fasthttp.RequestCtx), path, body, method string, t *testing.T) (int, []byte) {
	s := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			handler(ctx)
		},
	}

	ln := fasthttputil.NewInmemoryListener()

	serverCh := make(chan struct{})
	go func() {
		if err := s.Serve(ln); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		close(serverCh)
	}()

	clientCh := make(chan struct{})
	statusCode := 0
	var respBody []byte
	go func() {
		c, err := ln.Dial()
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		var req string
		if len(body) > 0 {
			reqs := "%s %s HTTP/1.1\r\nContent-Type: application/json\r\nContent-Length: %d\r\n\r\n%s"
			req = strings.Trim(fmt.Sprintf(reqs, method, path, len(body), body), "\r\n ")
		} else {
			reqs := "%s %s HTTP/1.1\r\nHost: 7amada\r\n\r\n"
			req = fmt.Sprintf(reqs, method, path)
		}

		if _, err = c.Write([]byte(req)); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		br := bufio.NewReader(c)
		var resp fasthttp.Response
		if err = resp.Read(br); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		statusCode = resp.StatusCode()
		respBody = resp.Body()
		close(clientCh)
	}()

	select {
	case <-clientCh:
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout")
	}

	if err := ln.Close(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	select {
	case <-serverCh:
	case <-time.After(time.Second):
		t.Fatalf("timeout")
	}
	return statusCode, respBody
}

func TestIntegration1(t *testing.T) {
	handler := NewHttpHandler()
	handler.GET("/api/sdk/v3/latest_version/:os", func(r Request) {
		os := r.Params().GetObject("os").StringValue()
		var version string
		switch os {
		case "android":
			version = "1.0"
		case "ios":
			version = "1.2"
		default:
			version = "0.0"
		}
		Respond(r, Json, map[string]string{"version": version})
	})
	status, body := tempServer(func(ctx *fasthttp.RequestCtx) {
		handler.handle(ctx)
	}, "/api/sdk/v3/latest_version/android", "", "GET", t)
	assert.Equal(t, 200, status)
	var parsedBody map[string]interface{}
	err := json.Unmarshal(body, &parsedBody)
	assert.NoError(t, err)
	assert.NotNil(t, parsedBody["version"])
	assert.Equal(t, "1.0", parsedBody["version"])
}

func TestIntegration2(t *testing.T) {
	handler := NewHttpHandler()
	handler.POST("/api/sdk/v3/bugs", func(r Request) {
		bugState := r.Params().GetObject("_data").GetObject("state_hash").StringValue()
		Respond(r, Json, map[string]string{"id": bugState})
	})
	status, body := tempServer(func(ctx *fasthttp.RequestCtx) {
		handler.handle(ctx)
	}, "/api/sdk/v3/bugs", "{\"state_hash\": \"7amada\"}", "POST", t)
	assert.Equal(t, 200, status)
	var parsedBody map[string]interface{}
	err := json.Unmarshal(body, &parsedBody)
	assert.NoError(t, err)
	assert.NotNil(t, parsedBody["id"])
	assert.Equal(t, "7amada", parsedBody["id"])
}

func TestIntegration3(t *testing.T) {
	handler := NewHttpHandler()
	handler.POST("/api/sdk/v3/bugs", func(r Request) {})
	status, _ := tempServer(func(ctx *fasthttp.RequestCtx) {
		handler.handle(ctx)
	}, "/api/sdk/v3", "", "POST", t)
	assert.Equal(t, 404, status)
}

func TestIntegration4(t *testing.T) {
	handler := NewHttpHandler()
	handler.POST("/api/sdk/v3/bugs", func(r Request) {})
	status, _ := tempServer(func(ctx *fasthttp.RequestCtx) {
		handler.handle(ctx)
	}, "/api/sdk/v3/bugs", "", "GET", t)
	assert.Equal(t, 405, status)
}
