package wordlistdb

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"

	"github.com/cockroachdb/pebble"
)

func getSHA256Sum(password string) string {
	data := sha256.New()
	io.WriteString(data, password)

	return fmt.Sprintf("%x", string(data.Sum(nil)))
}

func getMD5Sum(password string) string {
	data := md5.New()
	io.WriteString(data, password)

	return fmt.Sprintf("%x", string(data.Sum(nil)))
}

func getSHA1Sum(password string) string {
	data := sha1.New()
	io.WriteString(data, password)

	return fmt.Sprintf("%x", string(data.Sum(nil)))
}

func AddPassword(client *pebble.DB, batch *pebble.Batch, password string) error {
	hashes := [3]string{getSHA256Sum(password), getSHA1Sum(password), getMD5Sum(password)}

	for i := range hashes {
		if err := batch.Set([]byte(hashes[i]), []byte(password), pebble.Sync); err != nil {
			return err
		}
	}

	return nil
}

func GetPassword(client *pebble.DB, password string) error {
	hashes := [3]string{getSHA256Sum(password), getSHA1Sum(password), getMD5Sum(password)}

	for i := range hashes {
		_, _, err := client.Get([]byte(hashes[i]))
		if err != nil {
			return err
		}
	}

	return nil
}
