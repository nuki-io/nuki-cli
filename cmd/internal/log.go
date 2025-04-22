package internal

import (
	"log/slog"
	"os"

	"go.nuki.io/nuki/nukictl/internal"
)

// Logger is a dedicated Logger that should receive all messages that go to stdout
var Logger = slog.New(internal.NewLogger(slog.LevelInfo, os.Stdout))
