package systemcalls

import (
	"reflect"
	"testing"
)

const (
	body = `first line
second line
third line`
)

func TestReadFile(t *testing.T) {

	t.Run("read not existing file", func(t *testing.T) {
		os := NewOS()
		_, err := os.ReadFile("../../notExistingFile")
		if err == nil {
			t.Errorf("err is not nil, but %v", err)
		}
	})

	t.Run("read one-liner file", func(t *testing.T) {
		os := NewOS()
		found, err := os.ReadFile("../../scripts/sample1")
		expected := []string{"test"}
		if err != nil {
			t.Errorf("err is not nil, but %v", err)
		}

		if !reflect.DeepEqual(found, expected) {
			t.Errorf("expected %v, found %v", expected, found)
		}

	})
	t.Run("read multi-liner file", func(t *testing.T) {
		os := NewOS()
		found, err := os.ReadFile("../../scripts/sample2")
		expected := []string{"assign x 1", "", "print x", "writeFile filename data", "semWait mutex1"}
		if err != nil {
			t.Errorf("err is not nil, but %v", err)
		}

		if !reflect.DeepEqual(found, expected) {
			t.Errorf("expected %v, found %v", expected, found)
		}
	})
}

func TestWriteToFile(t *testing.T) {

	t.Run("write to file file", func(t *testing.T) {

		filepath := "../../scripts/new_file"

		os := NewOS()
		err := os.WriteToFile(filepath, body)
		if err != nil {
			t.Errorf("expected nil, found %v", err)
		}

		found, err := os.ReadFile(filepath)
		if err != nil {
			t.Errorf("expected nil, found %v", err)
		}

		expected := []string{"first line", "second line", "third line"}

		if !reflect.DeepEqual(expected, found) {
			t.Errorf("expected %v, found %v", expected, found)
		}

		err = os.DeleteFile(filepath)
		if err != nil {
			t.Errorf("expected nil, found %v", err)
		}
	})
}
