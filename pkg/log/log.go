package logger

import (
	"context"
	"log/slog"
	"net/url"
	"os"
	"slices"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Logger struct {
	slogger  *slog.Logger
	logLevel *slog.LevelVar
	area     string
}

const (
	LevelTrace   = slog.Level(-8)
	LevelVerbose = slog.Level(-4)
	LevelInfo    = slog.Level(0)
	LevelWarn    = slog.Level(4)
	LevelError   = slog.Level(8)
	LevelFatal   = slog.Level(12)
)

var LevelNames = map[slog.Leveler]string{
	LevelTrace:   "TRACE",
	LevelVerbose: "VERBOSE",
	LevelInfo:    "INFO",
	LevelWarn:    "WARN",
	LevelError:   "ERROR",
	LevelFatal:   "FATAL",
}

var filteredQueryParams = [...]string{
	"code",
	"state",
}

func New(area string, level slog.Level) *Logger {
	logLevel := new(slog.LevelVar)
	logLevel.Set(level)
	slogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				level := a.Value.Any().(slog.Level)
				levelLabel, exists := LevelNames[level]

				if !exists {
					levelLabel = level.String()
				}

				a.Value = slog.StringValue(levelLabel)
			}

			return a
		},
	}))

	return &Logger{
		area:     area,
		slogger:  slogger,
		logLevel: logLevel,
	}
}

func (l *Logger) SetLogLevel(level slog.Level) {
	l.logLevel.Set(level)
}

func (l *Logger) Info(msg string, args ...any) {
	ctx := context.Background()
	l.slogger.Log(ctx, LevelInfo, msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	ctx := context.Background()
	l.slogger.Log(ctx, LevelWarn, msg, args...)
}

func (l *Logger) Verbose(msg string, args ...any) {
	ctx := context.Background()
	l.slogger.Log(ctx, LevelVerbose, msg, args...)
}

// TODO fix issue with trace logging after config was loaded
func (l *Logger) Trace(msg string, args ...any) {
	ctx := context.Background()
	l.slogger.Log(ctx, LevelTrace, msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	ctx := context.Background()
	l.slogger.Log(ctx, LevelError, msg, args...)
}

func (l *Logger) Fatal(msg string, args ...any) {
	ctx := context.Background()
	l.slogger.Log(ctx, LevelFatal, msg, args...)
}

func GetEchoLogger() echo.MiddlewareFunc {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			parsedUrl, err := url.Parse(v.URI)
			if err != nil {
				return err
			}

			query := parsedUrl.Query()
			for key := range query {
				if slices.Contains(filteredQueryParams[:], key) {
					query[key] = []string{"REDACTED"}
				}
			}
			parsedUrl.RawQuery = query.Encode()

			if v.Error == nil {
				logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("method", v.Method),
					slog.String("uri", parsedUrl.RequestURI()),
					slog.Int("status", v.Status),
				)
			} else {
				logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String("method", v.Method),
					slog.String("uri", parsedUrl.RequestURI()),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	})
}
