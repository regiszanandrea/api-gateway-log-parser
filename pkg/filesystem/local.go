package filesystem

import (
	"bufio"
	"os"
	"path/filepath"
)

type Local struct{}

func NewLocalFileSystem() *Local {
	return &Local{}
}

func (l *Local) Open(path string) (*os.File, error) {
	absPath, _ := filepath.Abs(path)

	file, err := os.Open(absPath)

	if err != nil {
		return nil, err
	}

	return file, nil
}

func (l *Local) GetScanner(file *os.File) *bufio.Scanner {
	return bufio.NewScanner(file)
}

func (l *Local) GetLine(scanner *bufio.Scanner) string {
	return scanner.Text()
}

func (l *Local) Write(path string, data string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return err
	}

	defer f.Close()

	if _, err = f.WriteString(data); err != nil {
		return err
	}

	return nil
}
