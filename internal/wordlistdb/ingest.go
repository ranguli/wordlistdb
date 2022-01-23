package wordlistdb

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/cockroachdb/pebble"
	"github.com/ranguli/wordlistdb/internal/pkg/password"
)

func Ingest(database string, passwordFile string) error {
	fmt.Printf("Ingesting wordlist \"%s\"\n", passwordFile)

	file, err := os.Open(passwordFile)
	if err != nil {
		return err
	}

	client, err := pebble.Open(database, &pebble.Options{})
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)

	batchSize := 0

	// In my testing, this number proved to anecdotally be the sweet spot for
	// maximizing performance, ingesting rockyou.txt in about 2:00 minutes.
	batchLimit := 250000
	batch := client.NewBatch()

	for scanner.Scan() {
		fmt.Println(batchSize)

		passwordString := scanner.Text()

		if !strings.HasPrefix(passwordString, "#!comment:") {
			fmt.Printf("Ingesting password:\t'%s'\n", passwordString)

			err := password.AddPassword(client, batch, passwordString)
			if err != nil {
				return err
			}
		}

		batchSize += 1

		if batchSize == batchLimit {
			batch.Commit(pebble.Sync)
			batch.Reset()
			batchSize = 0
		}
	}

	if scanner.Err() != nil {
		return err
	}

	if err := client.Close(); err != nil {
		return err
	}

	return nil
}
