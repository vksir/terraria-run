package log

import (
	"golang.org/x/exp/slog"
	"io"
	"os"
	"path/filepath"
	"terraria-run/internal/common/constant"
)

func init() {
	w := getLogWriter()
	slog.SetDefault(slog.New(slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				a.Value = slog.StringValue(filepath.Base(a.Value.String()))
			}
			return a
		},
	}.NewTextHandler(w)))
}

func getLogWriter() io.Writer {
	f, err := os.OpenFile(constant.ServerLogPath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	return io.MultiWriter(f, os.Stdout)
}
