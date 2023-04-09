package app

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"
)

type ErrorMockReader struct{}

func (e ErrorMockReader) Read([]byte) (n int, err error) {
	return 0, errors.New("Dummy Error")
}

func TestInit(t *testing.T) {
	testString := "Tatami Galaxy\n"
	selectedAnimeMessage := fmt.Sprintf(endMessage, testString)

	t.Run("It returns a message when run", func(t *testing.T) {
		mockReader := strings.NewReader(testString)
		mockWriter := &bytes.Buffer{}

		Init(mockWriter, mockReader)

		want := startMessage + selectedAnimeMessage
		got := mockWriter.String()

		assertString(t, got, want)
	})

	t.Run("It returns an error if could not read user input", func(t *testing.T) {
		mockReader := ErrorMockReader{}
		mockWriter := &bytes.Buffer{}

		err := Init(mockWriter, mockReader)

		expectError(t, err)
	})
}

func assertString(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("Expected %q but got %q", want, got)
	}
}

func expectError(t *testing.T, got error) {
	t.Helper()
	if got == nil {
		t.Errorf("Expected an error but got none")
	}
}
