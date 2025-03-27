package handlers

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"triple-s/metadata"
	"triple-s/utils"
)

var (
	objectPath string
	bucketName string
)

func objectValidation(w http.ResponseWriter, r *http.Request) bool {
	parts := strings.SplitN(strings.TrimPrefix(r.URL.Path, "/"), "/", 2)
	counts := len(parts)
	if counts != 2 {
		SendResponse(w, http.StatusBadRequest, "Invalid request: object key required")
		return false
	}

	bucketName = parts[0]
	objectName := parts[1]

	objectPath = filepath.Join(DirectoryPath, bucketName, objectName)
	return true
}

func PutObject(w http.ResponseWriter, r *http.Request) {
	ok := objectValidation(w, r)
	if !ok {
		return
	}
	bucketPath := filepath.Join(DirectoryPath, bucketName)
	if _, err := os.Stat(bucketPath); os.IsNotExist(err) {
		SendResponse(w, http.StatusNotFound, "Bucket does not exist")
		return
	}

	file, err := os.Create(objectPath)
	if err != nil {
		SendResponse(w, http.StatusInternalServerError, "Failed to create object")
		return
	}

	if strings.HasSuffix(objectPath, "object.csv") {
		SendResponse(w, http.StatusBadRequest, "Object.csv cannot be the object")
		return
	}

	defer file.Close()

	written, err := io.Copy(file, r.Body)
	if err != nil {
		SendResponse(w, http.StatusInternalServerError, "Failed to write object")
		return
	}

	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	err = recordObjectMetadata(bucketName, filepath.Base(objectPath), written, contentType)
	if err != nil {
		SendResponse(w, http.StatusInternalServerError, "Failed to update object metadata")
		return
	}
	metadata.UpdateBucket(bucketName, DirectoryPath)
	SendResponse(w, http.StatusCreated, "Object uploaded successfully")
}

func GetObject(w http.ResponseWriter, r *http.Request) {
	ok := objectValidation(w, r)
	if !ok {
		return
	}
	file, err := os.Open(objectPath)
	if os.IsNotExist(err) {
		SendResponse(w, http.StatusNotFound, "Object not found")
		return
	} else if err != nil {
		SendResponse(w, http.StatusInternalServerError, "Failed to retrieve object")
		return
	}
	defer file.Close()

	http.ServeContent(w, r, filepath.Base(objectPath), time.Now(), file)
}

func DeleteObject(w http.ResponseWriter, r *http.Request) {
	ok := objectValidation(w, r)
	if !ok {
		return
	}
	objectKey := filepath.Base(objectPath)

	if !isObjectInMetadata(bucketName, objectKey) {
		SendResponse(w, http.StatusNotFound, "Object not found in metadata")
		return
	}

	if _, err := os.Stat(objectPath); os.IsNotExist(err) {
		SendResponse(w, http.StatusNotFound, "Object file not found")
		return
	}

	err := os.Remove(objectPath)
	if err != nil {
		SendResponse(w, http.StatusInternalServerError, "Failed to delete object file")
		return
	}

	err = removeObjectMetadata(bucketName, objectKey)
	if err != nil {
		SendResponse(w, http.StatusInternalServerError, "Failed to update object metadata")
		return
	}

	err = utils.RemoveCSV(filepath.Join(DirectoryPath, bucketName, "objects.csv"))
	if err != nil {
		SendResponse(w, http.StatusInternalServerError, "Failed to clean up empty metadata file")
		return
	}
	metadata.UpdateBucket(bucketName, DirectoryPath)
	SendResponse(w, http.StatusNoContent, "Object deleted successfully")
}

func recordObjectMetadata(bucketName, objectKey string, size int64, contentType string) error {
	metadataFile := filepath.Join(DirectoryPath, bucketName, "objects.csv")

	var updatedRecords [][]string
	if _, err := os.Stat(metadataFile); err == nil {
		file, err := os.Open(metadataFile)
		if err != nil {
			return err
		}
		defer file.Close()

		reader := csv.NewReader(file)
		records, err := reader.ReadAll()
		if err != nil {
			return err
		}

		for _, record := range records {
			if len(record) > 0 && record[0] != objectKey {
				updatedRecords = append(updatedRecords, record)
			}
		}
	}

	newRecord := []string{
		objectKey,
		fmt.Sprintf("%d", size),
		contentType,
		time.Now().Format(time.RFC3339),
	}
	updatedRecords = append(updatedRecords, newRecord)

	file, err := os.Create(metadataFile)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	return writer.WriteAll(updatedRecords)
}

func removeObjectMetadata(bucketName, objectKey string) error {
	metadataFile := filepath.Join(DirectoryPath, bucketName, "objects.csv")

	content, err := os.ReadFile(metadataFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	lines := strings.Split(string(content), "\n")
	var updatedLines []string

	for _, line := range lines {
		if !strings.HasPrefix(line, objectKey+",") && line != "" {
			updatedLines = append(updatedLines, line)
		}
	}
	joined := []byte(strings.Join(updatedLines, "\n"))
	return os.WriteFile(metadataFile, joined, 0o644)
}

func isObjectInMetadata(bucketName, objectKey string) bool {
	metadataFile := filepath.Join(DirectoryPath, bucketName, "objects.csv")

	file, err := os.Open(metadataFile)
	if os.IsNotExist(err) {
		return false
	} else if err != nil {
		fmt.Println("Error opening metadata file:", err)
		return false
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading metadata file:", err)
		return false
	}

	for _, record := range records {
		if len(record) > 0 && record[0] == objectKey {
			return true
		}
	}

	return false
}
