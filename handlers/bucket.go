package handlers

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"strings"
	"triple-s/metadata"
	"triple-s/utils"
)

var DirectoryPath string

type Bucket struct {
	Name         string `xml:"Name"`
	CreationDate string `xml:"CreationDate"`
	LastModified string `xml:"LastModified"`
}

type BucketList struct {
	XMLName xml.Name `xml:"ListBucketResult"`
	Buckets []Bucket `xml:"Buckets"`
}

type ResponseMessage struct {
	XMLName xml.Name `xml:"ResponseMessage"`
	Status  int      `xml:"Status"`
	Message string   `xml:"Message"`
}

func PutBucket(w http.ResponseWriter, r *http.Request) {
	bucketName := strings.TrimPrefix(r.URL.Path, "/")
	if bucketName == "" {
		SendResponse(w, http.StatusBadRequest, "Bucket name is required")
		return
	}

	if err := utils.ValidateBucketName(bucketName); err != nil {
		SendResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	bucketPath := DirectoryPath + "/" + bucketName
	if _, err := os.Stat(bucketPath); !os.IsNotExist(err) {
		SendResponse(w, http.StatusNotFound, "Bucket already exists")
		return
	}

	err := os.Mkdir(bucketPath, 0o755)
	if err != nil {
		SendResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	err = metadata.RecordBucket(bucketName, DirectoryPath)
	if err != nil {
		fmt.Println(err)
		SendResponse(w, http.StatusInternalServerError, "Failed to update bucket metadata")
		return
	}

	SendResponse(w, http.StatusCreated, "Bucket created succesfully")
}

func GetBucket(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Path != "/" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	metadataFile := DirectoryPath + "/buckets.csv"

	file, err := os.Open(metadataFile)
	if os.IsNotExist(err) {
		sendResponseWithBuckets(w, http.StatusOK, nil)
		return
	} else if err != nil {
		SendResponse(w, http.StatusInternalServerError, "Failed to read bucket metadata")
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		SendResponse(w, http.StatusInternalServerError, "Failed to parse bucket metadata")
		return
	}

	var buckets []Bucket
	for _, record := range records {
		if len(record) >= 3 {
			buckets = append(buckets, Bucket{
				Name:         record[0],
				CreationDate: record[1],
				LastModified: record[2],
			})
		}
	}

	sendResponseWithBuckets(w, http.StatusOK, buckets)
}

func DeleteBucket(w http.ResponseWriter, r *http.Request) {
	bucketName := strings.TrimPrefix(r.URL.Path, "/")
	if bucketName == "" {
		SendResponse(w, http.StatusBadRequest, "Bucket name is required")
		return
	}

	if !isBucketExists(bucketName) {
		SendResponse(w, http.StatusNotFound, "Bucket not found")
		return
	}

	bucketPath := DirectoryPath + "/" + bucketName

	if _, err := os.Stat(bucketPath); os.IsNotExist(err) {
		SendResponse(w, http.StatusNotFound, "Bucket not found")
		return
	}

	files, err := os.ReadDir(bucketPath)
	if err != nil {
		SendResponse(w, http.StatusInternalServerError, "Bucket reading error")
		return
	}

	if len(files) > 0 {
		SendResponse(w, http.StatusConflict, "Bucket is not empty")
		return
	}

	err = os.Remove(bucketPath)
	if err != nil {
		SendResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	err = metadata.DeleteBucket(bucketName, DirectoryPath)
	if err != nil {
		SendResponse(w, http.StatusInternalServerError, "Failed to update bucket metadata")
		return
	}

	SendResponse(w, http.StatusNoContent, "Bucket deleted succesfully")
}

func isBucketExists(bucketName string) bool {
	metadataFile := DirectoryPath + "/buckets.csv"

	file, err := os.Open(metadataFile)
	if os.IsNotExist(err) {
		return false
	} else if err != nil {
		fmt.Println("Failed to read bucket metadata")
		return false
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Failed to parse bucket metadata")
		return false
	}

	for _, record := range records {
		if len(record) > 0 && record[0] == bucketName {
			return true
		}
	}

	return false
}

func SendResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(statusCode)
	xml.NewEncoder(w).Encode(ResponseMessage{Status: statusCode, Message: message})
}

func sendResponseWithBuckets(w http.ResponseWriter, statusCode int, buckets []Bucket) {
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(statusCode)
	xml.NewEncoder(w).Encode(BucketList{Buckets: buckets})
}
