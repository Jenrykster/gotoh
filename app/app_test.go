package app

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/Jenrykster/gotoh/api"
)

type ErrorMockReader struct{}

func (e ErrorMockReader) Read([]byte) (n int, err error) {
	return 0, errors.New("Dummy Error")
}

var dummyAnime = api.Anime{
	Id:          "7785",
	Name:        "The Tatami Galaxy",
	Description: "...",
}
var dummyAnime2 = api.Anime{
	Id:          "8985",
	Name:        "The Tatami Galaxy Specials",
	Description: "...",
}

var dummyAnimeList = []api.Anime{dummyAnime, dummyAnime2}

func TestInit(t *testing.T) {
	mockAnimeClient, server := api.NewMockAnimeClient()
	defer server.Close()

	selectedAnimeMessage := fmt.Sprintf(endMessage, dummyAnime.Name, dummyAnime.Description)

	t.Run("It returns a message when run", func(t *testing.T) {
		mockReader := strings.NewReader("12345" + "\n")
		mockWriter := &bytes.Buffer{}

		Init(mockWriter, mockReader, &mockAnimeClient)

		want := startMessage + selectedAnimeMessage
		got := mockWriter.String()

		assertString(t, got, want)
	})

	t.Run("It returns an error if could not read user input", func(t *testing.T) {
		mockReader := ErrorMockReader{}
		mockWriter := &bytes.Buffer{}

		err := Init(mockWriter, mockReader, &mockAnimeClient)

		expectError(t, err)
	})
}

func TestGetAnimeData(t *testing.T) {
	mockAnimeClient, server := api.NewMockAnimeClient()
	defer server.Close()

	testId := "7785"

	t.Run("It should return the anime data when a id is passed", func(t *testing.T) {
		got, err := getAnimeData(testId, &mockAnimeClient)
		expected := &dummyAnime

		if err != nil {
			t.Errorf("Unexpected error: %q", err)
		}

		if !reflect.DeepEqual(got, expected) {
			t.Errorf("Expected: %+v\n Got: %+v", expected, got)
		}
	})
}

func TestListAnime(t *testing.T) {
	mockPaginatedClient, server := api.NewMockPaginatedAnimeClient()
	defer server.Close()
	t.Run("It returns a list of animes when a name is passed", func(t *testing.T) {
		got, err := searchAnimes("Tatami", &mockPaginatedClient)
		expected := &dummyAnimeList

		if err != nil {
			t.Errorf("Unexpected error: %q", err)
		}

		if !reflect.DeepEqual(got, expected) {
			t.Errorf("Expected: %+v\nGot: %+v", expected, got)
		}
	})
}

func assertString(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("Expected %q \nGot %q", want, got)
	}
}

func expectError(t *testing.T, got error) {
	t.Helper()
	if got == nil {
		t.Errorf("Expected an error but got none")
	}
}
