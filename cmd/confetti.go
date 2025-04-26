/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/maaslalani/confetty/confetti"
	"github.com/spf13/cobra"
)

// confettiCmd represents the confetti command
var confettiCmd = &cobra.Command{
	Use: "confetti",
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(confetti.InitialModel(), tea.WithAltScreen())
		_, err := p.Run()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(confettiCmd)
}
