/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	verbose *bool
	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "godex",
		Short: "A powerful file management CLI tool",
		Long:  `godex is a file management CLI built in Go, designed for advanced file operations.It supports zipping, file backup ,versioning and much more with a clean and extensible interface.`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		Run: func(cmd *cobra.Command, args []string) {
			if *verbose {
				fmt.Println("Verbose mode enabled")
			}
			fmt.Println("Welcome to godex CLI! Use --help to see available commands.")
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.godex.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	verbose = rootCmd.Flags().BoolP("verbose", "v", false, "Enable verbose output")
}
