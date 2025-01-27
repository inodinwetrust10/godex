package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.inodinwetrust10/goPhile/pkg"
)

var zipCmd = &cobra.Command{
	Use:   "zip [output.zip] [files...]",
	Short: "Zip one or more files into a .zip archive",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		outputFile := args[0]
		files := args[1:]

		if err := pkg.ZipFiles(outputFile, files); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Println("Zipped successfully:", outputFile)
	},
}

func init() {
	rootCmd.AddCommand(zipCmd)
}
