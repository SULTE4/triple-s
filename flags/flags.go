package flags

import (
	"flag"
	"fmt"
	"os"

	"triple-s/utils"
)

func Flags() (int, string, error) {

	tripleS := flag.NewFlagSet("triple-s", flag.ExitOnError)

	port := tripleS.Int("port", 8080, "Port number")
	dir := tripleS.String("dir", "data", "Path to the directory")
	help := tripleS.Bool("help", false, "Show this screen")
	tripleS.Usage = utils.PrintUsage

	tripleS.Parse(os.Args[1:])

	if *help {
		utils.PrintUsage()
		os.Exit(0)
	}

	if *port < 1024 || *port > 65535 {
		return 0, "", fmt.Errorf("invalid port number: %d", *port)
	}

	positionalArgs := tripleS.Args()

	if len(positionalArgs) > 0 {
		return 0, "", fmt.Errorf("unexpected positional arguments: %v", positionalArgs)
	}

	if *dir == "" {
		return 0, "", fmt.Errorf("directory not provided")
	}

	return *port, *dir, nil
}
