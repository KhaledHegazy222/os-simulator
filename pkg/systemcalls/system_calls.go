package systemcalls

import (
	"fmt"
	"os"
	"strings"
)

type OS struct{}

// NewOS creates new object of os struct.
func NewOS() *OS {
	return &OS{}
}

// ReadFile read file from disk given its path.
func (o *OS) ReadFile(path string) ([]string, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return []string{}, err
	}

	lines := strings.Split(string(bytes), "\n")
	return lines, nil
}

// WriteToFile write data to existing file on disk and create it if not existed.
func (o *OS) WriteToFile(path string, data string) error {
	return os.WriteFile(path, []byte(data), 0666)
}

// WriteToFile delete file from file system.
func (o *OS) DeleteFile(path string) error {
	return os.Remove(path)
}

// PrintToStdOut print given data to the screen
func (o *OS) PrintToStdOut(data string) {
	fmt.Println(data)
}

// GetInput takes input from the user
func (o *OS) GetInput() string {
	var input string
	fmt.Scanln(&input)
	return input
}

// ReadFromMemory read specific location of memory
func (o *OS) ReadFromMemory() {
}

// WriteToMemory write to specific location of memory
func (o *OS) WriteToMemory() {
}
