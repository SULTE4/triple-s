package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"triple-s/flags"
	"triple-s/utils"
)

func main() {
	port, dir, err := flags.Flags()
	if err != nil {
		fmt.Println(err)
		utils.PrintUsage()
		os.Exit(1)
	}

	// err = utils.InitDirectory(dir)
	// utils.ErrorPrinting(err)

	http.HandleFunc("/", handleRequest)

	fmt.Printf("Server running on port %d with BaseDir %s\n", port, dir)

	addr := ":" + strconv.Itoa(port)
	err = http.ListenAndServe(addr, nil)
	utils.ErrorPrinting(err)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println("-----------------------------------------------------------------------")
	fmt.Println()
	fmt.Println(w)
	fmt.Println()
	fmt.Println("-----------------------------------------------------------------------")
	fmt.Println()
	fmt.Println("Request method:", r.Method)
	fmt.Println("Request URL:", r.URL.Path)
	fmt.Println()
	fmt.Println("-----------------------------------------------------------------------")
	fmt.Println()
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
	fmt.Println()
	fmt.Println()
}
