package utils

import (
	"fmt"
	"os"
)

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

func ErrorPrinting(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: "+err.Error())
		os.Exit(1)
	}
}
