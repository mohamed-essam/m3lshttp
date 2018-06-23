package m3lshttp

import "github.com/mohamed-essam/m3lsh"

type BadRequest struct {
	m3lsh.BaseException
}

type Unauthorized struct {
	m3lsh.BaseException
}

type PaymentRequired struct {
	m3lsh.BaseException
}

type Forbidden struct {
	m3lsh.BaseException
}

type NotFound struct {
	m3lsh.BaseException
}

type MethodNotAllowed struct {
	m3lsh.BaseException
}

type TimedOut struct {
	m3lsh.BaseException
}

type UnprocessableEntity struct {
	m3lsh.BaseException
}

type InternalServerError struct {
	m3lsh.BaseException
}
