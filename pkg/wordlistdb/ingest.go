package wordlistdb

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/cockroachdb/pebble"
	"github.com/ranguli/wordlistdb/internal/pkg/file"
	"github.com/schollz/progressbar/v3"
)

func Ingest(database string, passwordFilename string) error {

	passwordFile, err := os.Open(passwordFilename)
	if err != nil {
		return err
	}

	lineCount, err := file.LineCounter(passwordFilename)
	fmt.Printf("Ingesting %d passwords from wordlist \"%s\"\n", lineCount, passwordFilename)

	client, err := pebble.Open(database, &pebble.Options{})
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(passwordFile)

	if err != nil {
		return err
	}

	batchSize := 0

	// In my testing, this number proved to anecdotally be the sweet spot for
	// maximizing performance, ingesting rockyou.txt in about 2:15.

	batchLimit := 250000
	batch := client.NewBatch()

	currentLine := 0

	progressBar := progressbar.Default(int64(lineCount))
	for scanner.Scan() {
		passwordString := scanner.Text()

		if !strings.HasPrefix(passwordString, "#!comment:") {
			// Instead of re-hashing every password a second time to check if
			// it exists in the datastore first, it's ~20x faster to add it
			// and simply overwrite it if it did exist.
			err := AddPassword(client, batch, passwordString)
			if err != nil {
				return err
			}
		}

		batchSize += 1
		currentLine += 1

		if batchSize == batchLimit || currentLine == lineCount {
			batch.Commit(pebble.Sync)
			batch.Reset()
			batchSize = 0
		}

		progressBar.Add(1)
	}

	if scanner.Err() != nil {
		return err
	}

	if err := client.Close(); err != nil {
		return err
	}

	return nil
}
