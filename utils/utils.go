package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func ErrorPrinting(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: "+err.Error())
		os.Exit(1)
	}
}

func InitDirectory(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	if path == "." || path == ".." || strings.Contains(path, "..") || strings.Contains(path, "/") || strings.Contains(path, "\\") {
		return fmt.Errorf("invalid directory path: %s", path)
	}

	bannedDirNames := []string{"cmd", "handlers", "utils", "flags", "metadata"}
	for _, bannedDirName := range bannedDirNames {
		if path == bannedDirName {
			return fmt.Errorf("invalid directory path: %s", path)
		}
	}

	info, err := os.Stat(absPath)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(absPath, 0o755); err != nil {
			return fmt.Errorf("failed to create dir: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("error accessing the directory: %w", err)
	} else if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", absPath)
	}

	testFile := filepath.Join(absPath, ".testfile")
	file, err := os.Create(testFile)
	if err != nil {
		return fmt.Errorf("problem with writing in the directory: %w", err)
	}

	file.Close()
	os.Remove(testFile)
	return nil
}

func ValidateBucketName(bucketName string) error {
	if len(bucketName) < 3 || len(bucketName) > 63 {
		return fmt.Errorf("Bucket name should be between 3 and 63 characters")
	}

	if strings.Contains(bucketName, "..") {
		return fmt.Errorf("Bucket name should not contain '..'")
	}

	namePattern := `^[a-z0-9]([a-z0-9\-\.]*[a-z0-9])?$`
	if matched, _ := regexp.MatchString(namePattern, bucketName); !matched {
		return errors.New("bucket name contains invalid characters or format")
	}

	return nil
}

func RemoveCSV(filePath string) error {
	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if info.Size() == 0 {
		return os.Remove(filePath)
	}
	return nil
}

func PrintUsage() {
	fmt.Println(`Simple Storage Service.
	
	**Usage:**
		triple-s [-port <N>] [-dir <S>]  
		triple-s --help
	
	**Options:**
	- --help     Show this screen.
	- --port N   Port number
	- --dir S    Path to the directory`)
}
