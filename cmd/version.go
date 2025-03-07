package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.inodinwetrust10/godex/pkg/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "File versioning operations",
	Long:  `Create, list, and restore versions of files`,
}

var createCmd = &cobra.Command{
	Use:   "create [filepath]",
	Short: "Create a new version of a file",
	Args:  cobra.ExactArgs(1),
	RunE:  createVersion,
}

func init() {
	versionCmd.AddCommand(createCmd)
	rootCmd.AddCommand(versionCmd)
}

func createVersion(cmd *cobra.Command, args []string) error {
	filePath, err := filepath.Abs(args[0])
	if err != nil {
		return err
	}
	id, err := version.GenerateVersionID(filePath)
	if err != nil {
		return err
	}
	meta, err := version.CreateFile(filePath, id, "First")
	if err != nil {
		return err
	}
	fmt.Println(meta.ID)
	return nil
}
