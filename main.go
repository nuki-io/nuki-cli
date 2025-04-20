/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"log/slog"
	"os"

	"go.nuki.io/nuki/nukictl/cmd"
	logger "go.nuki.io/nuki/nukictl/internal"
)

func main() {
	logger := slog.New(logger.NewLogger(slog.LevelInfo, os.Stderr))
	slog.SetDefault(logger)

	cmd.Execute()
}
