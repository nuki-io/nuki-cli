/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"log/slog"
	"os"

	"go.nuki.io/nuki/nukictl/cmd"
	_ "go.nuki.io/nuki/nukictl/cmd/devices"
	logger "go.nuki.io/nuki/nukictl/internal"
)

func main() {
	l := &slog.LevelVar{}
	l.Set(slog.LevelInfo)
	logger := slog.New(logger.NewLogger(l, os.Stderr))
	slog.SetDefault(logger)

	cmd.Execute(l)
}
