package wordlistdb

import (
	"errors"
	"fmt"

	"github.com/cockroachdb/pebble"
)

// Performs a search for a given key in database
func Search(database string, hash string) (string, error) {

	client, err := pebble.Open(database, &pebble.Options{})

	if err != nil {
		return "", err
	}

	plaintextBytes, closer, err := client.Get([]byte(hash))

	if err != nil {
		return "", errors.New(fmt.Sprintf("No plaintext password found with hash \"%s\"", hash))
	}

	err = closer.Close()
	if err != nil {
		return "", err
	}

	err = client.Close()
	if err != nil {
		return "", err
	}

	plaintext := string(plaintextBytes)
	return plaintext, nil
}
