package file

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
)

func Exists(filename string) bool {
	exists := true

	_, err := os.OpenFile(filename, os.O_RDONLY, 0o444)
	if errors.Is(err, os.ErrNotExist) {
		exists = false
	}

	return exists
}

func FileErrorMessage(file string) string {
	return fmt.Sprintf("Could not open file \"%s\".", file)
}

func LineCounter(filename string) (int, error) {
	// Derived from https://stackoverflow.com/a/24563853

	fileObject, err := os.Open(filename)
	if err != nil {
		return 0, err
	}

	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := fileObject.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			fileObject.Close()
			return count, nil

		case err != nil:
			fileObject.Close()
			return 0, err
		}
	}
}
