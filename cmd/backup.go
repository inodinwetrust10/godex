package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.inodinwetrust10/godex/pkg"
)

var backupCmd = &cobra.Command{
	Use:   "backup [file]",
	Short: "Backup file to Google Drive",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := pkg.UploadFile(args[0])
		if err != nil {
			log.Fatalf("Backup failed: %v", err)
		}
		fmt.Println("Backup successful!")
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
}
