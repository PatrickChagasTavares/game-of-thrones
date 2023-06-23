package logger

import "context"

type (
	Logger interface {
		Info(arg ...any)
		InfoContext(ctx context.Context, arg ...any)
		Error(arg ...any)
		ErrorContext(ctx context.Context, arg ...any)
		Fatal(arg ...any)
	}
)
