package log

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

const key = "logger"

func Middleware(logger *zap.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger = logger.With(zap.String("request-id", r.Header.Get("X-Request-ID")))
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), key, logger)))
		})
	}
}

func Parse(ctx context.Context) *zap.Logger {
	i := ctx.Value(key)
	if logger, ok := i.(*zap.Logger); ok {
		return logger
	}

	return zap.NewExample()
}

type GormLogger struct {
	*zap.Logger
}

func (logger GormLogger) Print(v ...interface{}) {
	logger.Logger.Sugar().Info(v...)
}
