package logger

import (
	"context"
	"log/slog"
)

func InfoCtx(ctx context.Context, msg string, attr ...slog.Attr) {
	slog.LogAttrs(ctx, slog.LevelInfo, msg, attr...)
}
