package middleware

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Event is the structured log record emitted once per HTTP request.
// It implements slog.LogValuer so its fields are logged as attributes.
type Event struct {
	RequestID string
	Path      string
	Method    string
	// ActorID is the authenticated user's ID, or nil if the request was
	// unauthenticated (e.g. /auth routes).
	ActorID *uint
}

// LoggingMiddleware emits one structured log line per request.
type LoggingMiddleware struct{}

// NewLoggingMiddleware returns a LoggingMiddleware.
func NewLoggingMiddleware() *LoggingMiddleware {
	return &LoggingMiddleware{}
}

// LogMiddleware returns a handler that assigns each request a UUID,
// exposes it via the X-Request-ID response header, and logs the request
// after all downstream handlers have run.
//
// The actor ID is read from the UserIDKey context value set by the auth
// middleware. Because logging happens after c.Next(), this middleware can
// be registered before the auth middleware and still see the actor.
//
// The log level is derived from the response status: 5xx is logged as
// Error, 4xx as Warn and everything else as Info.
func (m *LoggingMiddleware) LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := uuid.NewString()

		c.Header("X-Request-ID", requestId)

		c.Next()

		e := Event{
			RequestID: requestId,
			Path:      c.Request.URL.Path,
			Method:    c.Request.Method,
		}

		if id, ok := c.Get(UserIDKey); ok {
			if v, ok := id.(uint); ok {
				e.ActorID = &v
			}
		}

		// Inline the log event, for direct access to event attributes.
		slog.Log(c.Request.Context(), levelFor(c.Writer.Status()), "request", slog.Any("", e))
	}
}

// LogValue implements slog.LogValuer.
func (e Event) LogValue() slog.Value {
	attrs := []slog.Attr{
		slog.String("requestId", e.RequestID),
		slog.String("path", e.Path),
		slog.String("method", e.Method),
	}
	if e.ActorID != nil {
		attrs = append(attrs, slog.Uint64("actorId", uint64(*e.ActorID)))
	}
	return slog.GroupValue(attrs...)
}

// levelFor maps an HTTP status code to the level it should be logged at.
func levelFor(status int) slog.Level {
	if status >= 500 {
		return slog.LevelError
	} else if status >= 400 {
		return slog.LevelWarn
	}
	return slog.LevelInfo
}
