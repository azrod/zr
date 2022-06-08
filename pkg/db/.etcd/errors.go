package etcd

import "errors"

var (
	ErrKeyInvalid = errors.New("key is invalid")

	ErrValueInvalid = errors.New("value is invalid")
)
