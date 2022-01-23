package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ranguli/wordlistdb/internal/pkg/file"
	"github.com/ranguli/wordlistdb/internal/wordlistdb"
)

var subcommandUsage = map[string]string{}

const DefaultDatabaseName = "wordlistdb-data"

func main() {

	subcommandUsage["ingest"] = "Add a new wordlist into the database"
	ingestCmd := flag.NewFlagSet("ingest", flag.ExitOnError)

	if len(os.Args) < 2 {
		printBanner()
		printUsage()
	}

	switch os.Args[1] {

	case "ingest":
		ingestCmd.Parse(os.Args[2:])
		ingestSubcommand(DefaultDatabaseName, ingestCmd.Args())
	default:
		printBanner()
		printUsage()
	}
}

func ingestSubcommand(database string, args []string) {
	filename := args[0]

	if len(args) != 0 {
		if !file.Exists(filename) {
			log.Fatal(file.FileErrorMessage(filename))
		}
	} else {
		fmt.Println("Please provide a filename.\n")
		printUsage()
	}

	err := wordlistdb.Ingest(database, filename)
	if err != nil {
		log.Fatal(fmt.Sprintf("Could not ingest: %v\n", err))
		log.Fatal(err)
	}
}

func printBanner() {
	fmt.Println("wordlistdb\n")
}

func printUsage() {
	fmt.Println("Usage:\n")
	for key, value := range subcommandUsage {
		fmt.Printf("%s \t\t %s \n", key, value)
	}
	os.Exit(1)
}
