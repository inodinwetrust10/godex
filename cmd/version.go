package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.inodinwetrust10/godex/pkg/version"
)

var (
	message    string
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "File versioning operations",
		Long:  `Create, list, and restore versions of files`,
	}
)

var createCmd = &cobra.Command{
	Use:   "create [filepath]",
	Short: "Create a new version of a file",
	Args:  cobra.ExactArgs(1),
	RunE:  createVersion,
}

var listCmd = &cobra.Command{
	Use:   "list [filepath]",
	Short: "List all versions of a file",
	Args:  cobra.ExactArgs(1),
	RunE:  listVersion,
}

var restoreCmd = &cobra.Command{
	Use:   "restore [filepath] [versionID]",
	Short: "Restore your file to specified versionID",
	Args:  cobra.ExactArgs(2),
	RunE:  restoreVersion,
}

func init() {
	createCmd.Flags().StringVarP(&message, "message", "m", "commit", "Add a commit message")
	versionCmd.AddCommand(createCmd)
	versionCmd.AddCommand(listCmd)
	versionCmd.AddCommand(restoreCmd)
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
	meta, err := version.CreateFile(filePath, id, message)
	if err != nil {
		return err
	}
	fmt.Println(meta.ID)
	return nil
}

// ///////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////
func listVersion(cmd *cobra.Command, args []string) error {
	filePath, err := filepath.Abs((args[0]))
	if err != nil {
		return err
	}
	filePath, err = version.GetVersionPath(filePath)
	if err != nil {
		return err
	}
	list, err := version.ListAllVersions(filePath)
	if err != nil {
		return err
	}
	for _, data := range list {
		fmt.Printf("ID: %s\n", data.ID)
		fmt.Printf("Message: %s\n", data.Message)
		fmt.Printf("Created At: %s\n", data.CreatedAt)
		fmt.Printf("Size(in Bytes): %d\n", data.Size)
	}
	return nil
}

// /////////////////////////////////////////////////////////////////////
// /////////////////////////////////////////////////////////////////////
func restoreVersion(cmd *cobra.Command, args []string) error {
	filePath, err := filepath.Abs(args[0])
	if err != nil {
		return err
	}
	versionDir, err := version.GetVersionPath(filePath)
	if err != nil {
		return err
	}
	err = version.RestoreFile(versionDir, args[1], filePath)
	if err != nil {
		return err
	}
	fmt.Printf("File restored to version %s successfully", args[1])
	return nil
}
