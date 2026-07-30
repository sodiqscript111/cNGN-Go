package service

import (
	"github.com/theobiabo/cNGN-Go/envelope"
	"github.com/theobiabo/cNGN-Go/error"
)

type sendFunc func(method, path string, query map[string]string, body any, result any) *error.Error

func sendRequest[T any](send sendFunc, method, path string, query map[string]string, body any) (*envelope.Response[T], *error.Error) {
	result := &envelope.Response[T]{}
	err := send(method, path, query, body, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func sendJSON[T any](send sendFunc, method, path string, query map[string]string, body any) (T, *error.Error) {
	var result T
	err := send(method, path, query, body, &result)
	if err != nil {
		var zero T
		return zero, err
	}
	return result, nil
}

func fmtUint32(n uint32) string {
	if n == 0 {
		return "0"
	}
	digits := make([]byte, 0, 10)
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}
