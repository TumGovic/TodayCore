// Package xerrors is the subset of Xray's common/errors used by FinalMask,
// logging through the sing-box logger.
package xerrors

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/sagernet/sing-box/log"
)

type Error struct {
	message string
	inner   error
}

func buildMessage(msg []any) string {
	var b strings.Builder
	for _, m := range msg {
		b.WriteString(fmt.Sprint(m))
	}
	return b.String()
}

// New creates an error with the concatenation of msg, like Xray.
func New(msg ...any) *Error {
	return &Error{message: buildMessage(msg)}
}

// Base sets the inner error.
func (e *Error) Base(err error) *Error {
	e.inner = err
	return e
}

func (e *Error) Error() string {
	if e.inner == nil {
		return e.message
	}
	return e.message + " > " + e.inner.Error()
}

func (e *Error) Unwrap() error {
	return e.inner
}

func Combine(maybeError ...error) error {
	return errors.Join(maybeError...)
}

func message(err error, msg []any) string {
	s := "[finalmask] " + buildMessage(msg)
	if err != nil {
		s += " > " + err.Error()
	}
	return s
}

func LogDebug(ctx context.Context, msg ...any) {
	log.StdLogger().DebugContext(ctx, message(nil, msg))
}

func LogDebugInner(ctx context.Context, err error, msg ...any) {
	log.StdLogger().DebugContext(ctx, message(err, msg))
}

func LogInfo(ctx context.Context, msg ...any) {
	log.StdLogger().InfoContext(ctx, message(nil, msg))
}

func LogInfoInner(ctx context.Context, err error, msg ...any) {
	log.StdLogger().InfoContext(ctx, message(err, msg))
}

func LogWarning(ctx context.Context, msg ...any) {
	log.StdLogger().WarnContext(ctx, message(nil, msg))
}

func LogWarningInner(ctx context.Context, err error, msg ...any) {
	log.StdLogger().WarnContext(ctx, message(err, msg))
}

func LogError(ctx context.Context, msg ...any) {
	log.StdLogger().ErrorContext(ctx, message(nil, msg))
}

func LogErrorInner(ctx context.Context, err error, msg ...any) {
	log.StdLogger().ErrorContext(ctx, message(err, msg))
}
