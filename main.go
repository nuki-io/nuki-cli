/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"log/slog"
	"os"

	"github.com/nuki-io/nuki-cli/cmd"
	_ "github.com/nuki-io/nuki-cli/cmd/ble"
	"github.com/nuki-io/nuki-cli/internal"
)

func main() {
	l := &slog.LevelVar{}
	l.Set(slog.LevelInfo)
	logger := slog.New(internal.NewLogger(l, os.Stderr))
	slog.SetDefault(logger)

	cmd.Execute(l)
}
