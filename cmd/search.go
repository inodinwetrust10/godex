package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"

	"github.inodinwetrust10/godex/pkg"
)

var (
	rootDir        string
	name           string
	minSize        int64
	maxSize        int64
	modifiedAfter  string
	modifiedBefore string
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search files with various criteria",
	Long: `Search for files in the specified root directory using various criteria:
- exact name match
- file size range
- modification date range`,
	Run: func(cmd *cobra.Command, args []string) {
		if rootDir == "" {
			rootDir = "."
		}

		var afterTime, beforeTime time.Time
		var err error

		if modifiedAfter != "" {
			afterTime, err = time.Parse("2006-01-02", modifiedAfter)
			if err != nil {
				log.Fatalf(
					"Invalid date format for --modified-after: %v. Use YYYY-MM-DD format.",
					err,
				)
			}
		}

		if modifiedBefore != "" {
			beforeTime, err = time.Parse("2006-01-02", modifiedBefore)
			if err != nil {
				log.Fatalf(
					"Invalid date format for --modified-before: %v. Use YYYY-MM-DD format.",
					err,
				)
			}
		}

		criteria := pkg.SearchCriteria{
			Name:    name,
			MinSize: minSize,
			MaxSize: maxSize,
			After:   afterTime,
			Before:  beforeTime,
		}

		// Perform search using the provided filters
		results, err := pkg.SearchFiles(rootDir, criteria)
		if err != nil {
			log.Fatalf("Error searching files: %v", err)
		}

		if len(results) == 0 {
			fmt.Println("No files found matching the criteria")
			return
		}

		fmt.Println("Found files:")
		for _, result := range results {
			fmt.Println(result)
		}
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().StringVarP(&rootDir, "path", "p", "",
		"Root path for the search (default is current directory)")

	searchCmd.Flags().StringVarP(&name, "name", "n", "",
		"Search by exact file name")

	searchCmd.Flags().Int64VarP(&minSize, "min-size", "m", 0,
		"Minimum file size in bytes")
	searchCmd.Flags().Int64VarP(&maxSize, "max-size", "M", 0,
		"Maximum file size in bytes")

	searchCmd.Flags().StringVarP(&modifiedAfter, "modified-after", "a", "",
		"Find files modified after this date (YYYY-MM-DD)")
	searchCmd.Flags().StringVarP(&modifiedBefore, "modified-before", "b", "",
		"Find files modified before this date (YYYY-MM-DD)")
}
