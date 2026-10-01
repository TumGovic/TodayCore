// Package xcommon is the subset of Xray's common package used by FinalMask.
package xcommon

import "io"

// Must panics if err is not nil.
func Must(err error) {
	if err != nil {
		panic(err)
	}
}

// Must2 panics if the second parameter is not nil, otherwise returns the first parameter.
func Must2[T any](v T, err error) T {
	Must(err)
	return v
}

// CloseIfExists closes v if it is an io.Closer.
func CloseIfExists(v any) error {
	if closer, ok := v.(io.Closer); ok && closer != nil {
		return closer.Close()
	}
	return nil
}
