package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"triple-s/flags"
	"triple-s/handlers"
	"triple-s/utils"
)

func main() {
	port, dir, err := flags.Flags()
	if err != nil {
		fmt.Println(err)
		utils.PrintUsage()
		os.Exit(1)
	}

	err = utils.InitDirectory(dir)
	utils.ErrorPrinting(err)

	handlers.DirectoryPath = dir

	mux := http.NewServeMux()

	mux.HandleFunc("PUT /{BucketName}", handlers.PutBucket)
	mux.HandleFunc("GET /", handlers.GetBucket)
	mux.HandleFunc("DELETE /{BucketName}", handlers.DeleteBucket)

	mux.HandleFunc("PUT /{BucketName}/{ObjectName}", handlers.PutObject)
	mux.HandleFunc("GET /{BucketName}/{ObjectName}", handlers.GetObject)
	mux.HandleFunc("DELETE /{BucketName}/{ObjectName}", handlers.DeleteObject)

	fmt.Printf("Server running on port %d with BaseDir %s\n", port, dir)

	addr := ":" + strconv.Itoa(port)
	err = http.ListenAndServe(addr, mux)
	utils.ErrorPrinting(err)
}
