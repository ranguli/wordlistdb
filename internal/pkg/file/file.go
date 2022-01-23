package file

import (
	"errors"
	"fmt"
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
