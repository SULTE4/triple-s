package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

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
	mux.Handle("/", &Router{})

	fmt.Printf("Server running on port %d with BaseDir %s\n", port, dir)

	addr := ":" + strconv.Itoa(port)
	err = http.ListenAndServe(addr, mux)
	utils.ErrorPrinting(err)
}

type Router struct{}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	pathParts := strings.Split(strings.TrimPrefix(req.URL.Path, "/"), "/")

	switch len(pathParts) {
	case 1:
		switch req.Method {
		case http.MethodGet:
			handlers.GetBucket(w, req)
		case http.MethodPut:
			handlers.PutBucket(w, req)
		case http.MethodDelete:
			handlers.DeleteBucket(w, req)
		}
	case 2:
		switch req.Method {
		case http.MethodGet:
			handlers.GetObject(w, req)
		case http.MethodPut:
			handlers.PutObject(w, req)
		case http.MethodDelete:
			handlers.DeleteObject(w, req)
		}

	default:
		handlers.SendResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}
