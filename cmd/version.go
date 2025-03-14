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
		Long:  `Create,list,check diffs,restore and remove versions of files`,
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
	Short: "Restore your file to a specific versionID",
	Args:  cobra.ExactArgs(2),
	RunE:  restoreVersion,
}

var (
	useLastVersion bool
	seeDiffCmd     = &cobra.Command{
		Use:   "diff [filepath1] [filepath2]",
		Short: "Check diffs between two files",
		Args:  cobra.MinimumNArgs(1),
		RunE:  seeDiff,
	}
)

var (
	versionToRemove string
	removeCmd       = &cobra.Command{
		Use:   "remove [filepath] [verionFlag]",
		Short: "Remove a specific version or all versions of a file",
		Args:  cobra.ExactArgs(1),
		RunE:  removeVersion,
	}
)

func init() {
	createCmd.Flags().StringVarP(&message, "message", "m", "commit", "Add a commit message")
	seeDiffCmd.Flags().
		BoolVarP(&useLastVersion, "default", "d", false, "Compare with the last version")
	removeCmd.Flags().StringVarP(&versionToRemove, "version", "v", "", "Remove a specific version")
	versionCmd.AddCommand(removeCmd)
	versionCmd.AddCommand(createCmd)
	versionCmd.AddCommand(listCmd)
	versionCmd.AddCommand(restoreCmd)
	versionCmd.AddCommand(seeDiffCmd)
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
	fmt.Printf("A version was successfully created with version ID: %s\n", meta.ID)
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
	for _, data := range *list {
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

////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////

func seeDiff(cmd *cobra.Command, args []string) error {
	var diffRes version.DiffResult
	if useLastVersion && len(args) == 1 {
		filePath, err := filepath.Abs(args[0])
		if err != nil {
			return err
		}
		fileDir, err := version.GetVersionPath(filePath)
		lastVersionPath := version.ReturnLastFilePath(fileDir)
		if lastVersionPath == "No file found" {
			return fmt.Errorf("No last version found. Make a version first to check")
		}

		diffRes, err = version.FileDiff(filePath, lastVersionPath)
		if err != nil {
			return err
		}
	} else if len(args) == 2 {
		filePath1, err := filepath.Abs(args[0])
		if err != nil {
			return err
		}

		filePath2, err := filepath.Abs(args[1])
		if err != nil {
			return err
		}

		diffRes, err = version.FileDiff(filePath1, filePath2)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("incorrect number of arguments: provide either one file with -d flag or two files to compare")
	}
	version.PrintDiffResults(&diffRes)
	return nil
}

// //////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////
func removeVersion(cmd *cobra.Command, args []string) error {
	filePath, err := filepath.Abs(args[0])
	if err != nil {
		return err
	}
	fileDir, err := version.GetVersionPath(filePath)
	if err != nil {
		return err
	}
	if cmd.Flags().Changed("version") {
		err := version.ClearVersion(fileDir, versionToRemove)
		if err != nil {
			return err
		}

		fmt.Printf("Version with verisonID %s is cleared", versionToRemove)
	} else {
		err := version.ClearAllVersion(fileDir)
		if err != nil {
			return err
		}
		fmt.Println("All versions of this file are now cleared")
	}
	return nil
}
