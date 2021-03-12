package mock

import (
	"bufio"
	"io"
	"os"

	"github.com/stretchr/testify/mock"
)

type FileSystemMock struct {
	mock.Mock
}

func (f *FileSystemMock) Open(path string) (*os.File, error) {
	args := f.Called(path)
	if args.Get(0) == nil {
		return nil, args.Error(0)
	}

	file := args.Get(0).(*os.File)

	return file, nil
}

func (f *FileSystemMock) GetScanner(file *os.File) *bufio.Scanner {
	args := f.Called(file)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*bufio.Scanner)
}

func (f *FileSystemMock) GetLine(scanner *bufio.Scanner) string {
	args := f.Called(scanner)

	var getLineCalls []mock.Call

	for _, call := range f.Calls {
		if call.Method == "GetLine" {
			getLineCalls = append(getLineCalls, call)
		}
	}

	if len(getLineCalls) == 2 {
		return ""
	}

	if args.Get(0) == nil {
		return ""
	}

	return args.Get(0).(string)
}

func (f *FileSystemMock) Write(path string, data string) error {
	args := f.Called(path, data)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(error)
}

type ReaderMock struct {
	Data string
	done bool
}

func NewReaderMock(data string) *ReaderMock {
	return &ReaderMock{data, false}
}

func (r *ReaderMock) Read(p []byte) (n int, err error) {
	copy(p, []byte(r.Data))
	if r.done {
		return 0, io.EOF
	}
	r.done = true
	return len([]byte(r.Data)), nil
}
