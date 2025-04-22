/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log/slog"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.nuki.io/nuki/nukictl/internal"
)

var (

	// Logger is a dedicated Logger that should receive all messages that go to stdout
	Logger *slog.Logger

	cfgFile string
	verbose bool
	level   *slog.LevelVar
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:              "nukictl",
	Short:            "Command line tool to manage and control Nuki devices and online services.",
	TraverseChildren: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(l *slog.LevelVar) {
	level = l
	err := RootCmd.Execute()

	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig, setupLogger)
	cobra.OnFinalize(writeConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.nukictl)")
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose mode with increased and more detailed log output")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".nukictl" (without extension).
		viper.SetConfigType("yaml")
		viper.SetConfigName(".nukictl")
		viper.SetConfigFile(path.Join(home, ".nukictl"))
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		slog.Info("Using config file ", "path", viper.ConfigFileUsed())
	}
}

func setupLogger() {
	Logger = slog.New(internal.NewLogger(level, os.Stdout))
	if verbose {
		level.Set(slog.LevelDebug)
		slog.Info("Verbose mode enabled")
	}
}

func writeConfig() {
	err := viper.WriteConfig()
	if err != nil {
		slog.Error("Failed to persist configuration to file", "err", err)
	}
}
