package version

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// ///////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////
func GetVersionPath(filePath string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get home directory: %w", err)
	}
	hash := sha256.Sum256([]byte(filePath))
	hashedName := hex.EncodeToString(hash[:])

	versionDir := filepath.Join(homeDir, ".config", "godex", "versions")
	versionFilePath := filepath.Join(versionDir, hashedName)

	err = os.MkdirAll(versionFilePath, 0755)
	if err != nil {
		return "", fmt.Errorf("could not create directory: %w", err)
	}

	return versionFilePath, nil
}

/////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////

func GenerateVersionID(filePath string) (string, error) {
	dirPath, err := GetVersionPath(filePath)
	if err != nil {
		return "", err
	}

	versionFilePath := filepath.Join(dirPath, "version.json")

	nextVersion := 1

	if _, err := os.Stat(versionFilePath); err == nil {
		fileData, err := os.ReadFile(versionFilePath)
		if err != nil {
			return "", fmt.Errorf("failed to read version file: %w", err)
		}

		var versions []VersionMetaData
		if err := json.Unmarshal(fileData, &versions); err != nil {
			var singleVersion VersionMetaData
			if err := json.Unmarshal(fileData, &singleVersion); err != nil {
				return "", fmt.Errorf("failed to parse version data: %w", err)
			}
			versions = []VersionMetaData{singleVersion}
		}

		if len(versions) > 0 {
			highestVersion := 0
			for _, version := range versions {
				if len(version.ID) > 1 && version.ID[0] == 'v' {
					vNum, err := strconv.Atoi(version.ID[1:])
					if err == nil && vNum > highestVersion {
						highestVersion = vNum
					}
				}
			}
			nextVersion = highestVersion + 1
		}
	}

	versionID := fmt.Sprintf("v%d", nextVersion)
	return versionID, nil
}

/////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////

func checkDiffs(filePath1, filePath2 string) (bool, error) {
	if filePath2 == "No file found" {
		return true, nil
	}
	content1, err := os.ReadFile(filepath.Clean(filePath1))
	if err != nil {
		return false, fmt.Errorf("Error reading the file1 %w", err)
	}

	content2, err := os.ReadFile(filepath.Clean(filePath2))
	if err != nil {
		return false, fmt.Errorf("Error reading the file2 %w", err)
	}

	dmp := diffmatchpatch.New()

	diffs := dmp.DiffMain(string(content1), string(content2), false)

	return !(len(diffs) == 1 && diffs[0].Type == diffmatchpatch.DiffEqual), nil
}

///////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

func ReturnLastFilePath(jsonDirPath string) string {
	filePath := filepath.Join(jsonDirPath, "version.json")
	var elements []VersionMetaData

	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return "No file found"
		}
		return ""
	}
	file, err := os.ReadFile(filePath)
	if err != nil {
		return ""
	}
	err = json.Unmarshal(file, &elements)
	if err != nil {
		return ""
	}
	num := len(elements)
	filename := fmt.Sprintf("v%d", num)
	returnPath := filepath.Join(jsonDirPath, filename)
	return returnPath
}

// //////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////
func PrintDiffResults(diffRes *DiffResult) {
	if diffRes.Identical {
		fmt.Println("Files are identical")
		return
	}

	fmt.Printf("Diff Type: %s\n", diffRes.DiffType)

	if diffRes.Message != "" {
		fmt.Printf("Message: %s\n", diffRes.Message)
	}

	if len(diffRes.DiffLines) == 0 {
		fmt.Println("No line differences found")
		return
	}

	fmt.Println("\nDifferences:")
	fmt.Println("-------------------------------------------")

	for _, diff := range diffRes.DiffLines {
		fmt.Printf("Line %d:\n", diff.LineNumber)
		if diff.Line1 == "" {
			fmt.Printf("+ %s\n", diff.Line2)
		} else if diff.Line2 == "" {
			fmt.Printf("- %s\n", diff.Line1)
		} else {
			fmt.Printf("- %s\n+ %s\n", diff.Line1, diff.Line2)
		}
		fmt.Println("-------------------------------------------")
	}

	fmt.Printf("\nTotal differences: %d\n", len(diffRes.DiffLines))
}

////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////

func ReturnLastSecondFilePath(jsonDirPath string) string {
	filePath := filepath.Join(jsonDirPath, "version.json")
	var elements []VersionMetaData
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return "No file found"
		}
		return ""
	}
	file, err := os.ReadFile(filePath)
	if err != nil {
		return ""
	}
	err = json.Unmarshal(file, &elements)
	if err != nil {
		return ""
	}
	if len(elements) < 2 {
		return "No previous version found to check"
	}
	secondLastVersionData := elements[len(elements)-2]
	filename := secondLastVersionData.ID
	returnPath := filepath.Join(jsonDirPath, filename)
	return returnPath
}
