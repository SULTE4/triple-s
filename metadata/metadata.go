package metadata

import (
	"encoding/csv"
	"os"
	"strings"
	"time"
)

func RecordBucket(bucketName string, DirectoryPath string) error {
	metadataFile := DirectoryPath + "/" + "buckets.csv"
	file, err := os.OpenFile(metadataFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	timestamp := time.Now().Format(time.RFC3339)
	record := []string{
		bucketName,
		timestamp,
		timestamp,
	}

	if err := writer.Write(record); err != nil {
		return err
	}

	return nil
}

func DeleteBucket(bucketName string, DirectoryPath string) error {
	metadataFile := DirectoryPath + "/" + "buckets.csv"

	data, err := os.ReadFile(metadataFile)
	if os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	var newLines []string

	for _, line := range lines {
		if !strings.HasPrefix(line, bucketName) && !(line == "") {
			newLines = append(newLines, line)
		}
	}

	os.WriteFile(metadataFile, []byte(strings.Join(newLines, "\n")), 0o644)
	return nil
}

func UpdateBucket(bucketName string, DirectoryPath string) {
	metadataFile := DirectoryPath + "/" + "buckets.csv"

	data, err := os.ReadFile(metadataFile)
	if os.IsNotExist(err) {
		return
	} else if err != nil {
		return
	}

	lines := strings.Split(string(data), "\n")
	var newLines []string

	for _, line := range lines {
		if strings.HasPrefix(line, bucketName) && !(line == "") {
			parts := strings.Split(line, ",")
			newLine := parts[0] + "," + parts[1] + "," + time.Now().Format(time.RFC3339)
			newLines = append(newLines, newLine)
		} else {
			newLines = append(newLines, line)
		}
	}

	os.WriteFile(metadataFile, []byte(strings.Join(newLines, "\n")), 0o644)
}
