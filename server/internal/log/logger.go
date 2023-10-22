package log

import (
	"context"
	"github.com/dqso/mincer/server/internal/entity"
	"log/slog"
	"os"
)

type Logger interface {
	With(args ...any) *slog.Logger
	WithGroup(name string) *slog.Logger
	Enabled(ctx context.Context, level slog.Level) bool
	Log(ctx context.Context, level slog.Level, msg string, args ...any)
	LogAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr)
	Debug(msg string, args ...any)
	DebugContext(ctx context.Context, msg string, args ...any)
	Info(msg string, args ...any)
	InfoContext(ctx context.Context, msg string, args ...any)
	Warn(msg string, args ...any)
	WarnContext(ctx context.Context, msg string, args ...any)
	Error(msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)
}

type config interface {
	LogLevel() slog.Level
}

func NewWithConfig(config config) *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     config.LogLevel(),
		//ReplaceAttr: nil,
	}))
}

func New() *slog.Logger {
	return NewWithConfig(&stdConfig{
		level: slog.LevelInfo,
	})
}

type stdConfig struct {
	level slog.Level
}

func (c stdConfig) LogLevel() slog.Level { return c.level }

func Module(module string) slog.Attr {
	return slog.String("module", module)
}

func Err(err error) slog.Attr {
	return slog.Any("error", err)
}

func Damage(damage entity.Damage) slog.Attr {
	return slog.Group("damage",
		slog.Int64("physical", int64(damage.Physical())),
		slog.Int64("magical", int64(damage.Magical())),
	)
}

func Point(point entity.Point) slog.Attr {
	return slog.Group("point",
		slog.Float64("x", point.X),
		slog.Float64("y", point.Y),
	)
}

func Stats(stats entity.PlayerStats) slog.Attr {
	return slog.Group("stats",
		slog.String("class", stats.Class().String()),
		Resist(stats),
		slog.Float64("radius", stats.Radius()),
		slog.Float64("speed", stats.Speed()),
		slog.Int64("max_hp", int64(stats.MaxHP())),
	)
}

func Resist(resist entity.Resist) slog.Attr {
	return slog.Group("resist",
		slog.Float64("physical", resist.PhysicalResist()),
		slog.Float64("magical", resist.MagicalResist()),
	)
}
