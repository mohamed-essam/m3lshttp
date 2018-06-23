package m3lshttp

import (
	"encoding/json"
)

func parseBody(req Request) Params {
	var parsedBody interface{}

	contentType := req.ContentType()
	switch contentType {
	case "application/json":
		parsedBody = parseJson(req.Body())
	case "multipart/form-data":
		parsedBody = mapToInterfaceMap(req.MultipartForm())
	default:
		parsedBody = string(req.Body())
	}

	return newParams(parsedBody)
}

func parseJson(body []byte) interface{} {
	var ret interface{}
	json.Unmarshal(body, &ret)
	return ret
}
