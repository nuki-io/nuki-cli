package cmd

import (
	"log/slog"
	"os"

	"go.nuki.io/nuki/nukictl/internal"
)

// logger is a dedicated logger that should receive all messages that go to stdout
var logger = slog.New(internal.NewLogger(slog.LevelInfo, os.Stdout))
