package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.inodinwetrust10/godex/pkg"
)

var unzipCmd = &cobra.Command{
	Use:   "unzip [input.zip] [destination]",
	Short: "Unzip a .zip archive to a destination directory",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		inputFile := args[0]
		destination := args[1]

		if err := pkg.UnzipFile(inputFile, destination); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Println("Unzipped successfully to:", destination)
	},
}

func init() {
	rootCmd.AddCommand(unzipCmd)
}
