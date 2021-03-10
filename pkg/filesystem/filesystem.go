package filesystem

import (
	"bufio"
	"os"
)

type API interface {
	Open(path string) (*os.File, error)
	GetScanner(file *os.File) *bufio.Scanner
	GetLine(scanner *bufio.Scanner) string
}
