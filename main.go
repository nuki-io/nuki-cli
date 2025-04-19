/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"go.nuki.io/nuki/nukictl/cmd"
)

func main() {
	start := time.Now()
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: false,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Value = slog.StringValue(fmt.Sprintf("[%d]", int(time.Since(start).Seconds())))
			}
			return a
		},
	}))
	slog.SetDefault(logger)
	cmd.Execute()
}
