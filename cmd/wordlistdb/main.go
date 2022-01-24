package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ranguli/wordlistdb/internal/pkg/file"
	"github.com/ranguli/wordlistdb/pkg/wordlistdb"
)

var subcommandUsage = map[string]string{}

const DefaultDatabaseName = "wordlistdb-data"

func main() {

	subcommandUsage["ingest"] = "Add a new wordlist into the database"
	subcommandUsage["search"] = "Search for a single hash in the database"

	ingestCmd := flag.NewFlagSet("ingest", flag.ExitOnError)
	searchCmd := flag.NewFlagSet("search", flag.ExitOnError)

	if len(os.Args) < 2 {
		printBanner()
		printUsage()
	}

	switch os.Args[1] {

	case "ingest":
		ingestCmd.Parse(os.Args[2:])
		ingestSubcommand(DefaultDatabaseName, ingestCmd.Args())
	case "search":
		searchCmd.Parse(os.Args[2:])
		searchSubcommand(DefaultDatabaseName, searchCmd.Args())
	default:
		printBanner()
		printUsage()
	}
}

func ingestSubcommand(database string, args []string) {

	if len(args) == 0 {
		fmt.Println("Please provide a filename.\n")
		printUsage()
	}

	filename := args[0]

	if !file.Exists(filename) {
		log.Fatal(file.FileErrorMessage(filename))
	}

	err := wordlistdb.Ingest(database, filename)
	if err != nil {
		log.Fatal(fmt.Sprintf("Could not ingest: %v\n", err))
	}
}

func searchSubcommand(database string, args []string) {

	if len(args) == 0 {
		fmt.Println("Please provide a hash to search for.\n")
		printUsage()
	}

	hash := args[0]
	plaintext, err := wordlistdb.Search(database, hash)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println(plaintext)
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
