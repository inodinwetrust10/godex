package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.inodinwetrust10/goPhile/pkg"
)

var zipCmd = &cobra.Command{
	Use:   "zip [output.zip] [files...]",
	Short: "Zip files or directories into a .zip archive",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		dirFlag, err := cmd.Flags().GetBool("dir")
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		if dirFlag {
			if len(args) != 2 {
				fmt.Println("Error: when using -d, specify output file and directory")
				os.Exit(1)
			}
			outputFile := args[0]
			directory := args[1]

			fileInfo, err := os.Stat(directory)
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
			if !fileInfo.IsDir() {
				fmt.Println("Error:", directory, "is not a directory")
				os.Exit(1)
			}

			if err := pkg.ZipDirectory(outputFile, directory); err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
		} else {
			if len(args) < 2 {
				fmt.Println("Error: need output file and at least one file to zip")
				os.Exit(1)
			}
			outputFile := args[0]
			files := args[1:]

			for _, file := range files {
				fileInfo, err := os.Stat(file)
				if err != nil {
					fmt.Println("Error:", err)
					os.Exit(1)
				}
				if fileInfo.IsDir() {
					fmt.Println("Error:", file, "is a directory (use -d flag for directories)")
					os.Exit(1)
				}
			}

			if err := pkg.ZipFiles(outputFile, files); err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
		}
		fmt.Println("Zipped successfully:", args[0])
	},
}

func init() {
	zipCmd.Flags().BoolP("dir", "d", false, "Zip a directory instead of individual files")
	rootCmd.AddCommand(zipCmd)
}
