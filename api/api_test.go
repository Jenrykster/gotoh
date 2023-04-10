package api

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestAnimeClient(t *testing.T) {
	testId := "7785"

	mockClient, server := NewMockAnimeClient()
	defer server.Close()

	failingServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	failingClient := newCustomAnimeClient(failingServer.Client(), failingServer.URL)
	defer failingServer.Close()

	t.Run("It returns the requested Anime when the method is called", func(t *testing.T) {
		got, err := mockClient.GetById(testId)
		expected := &Anime{Id: testId, Name: "The Tatami Galaxy", Description: "BlaBlaBla..."}

		if err != nil {
			t.Errorf("Got unexpected error: %q", err)
		}

		if !reflect.DeepEqual(got, expected) {
			t.Errorf("Expected: %+v\n Got: %+v", expected, got)
		}
	})
	t.Run("It should return an error if the request fails", func(t *testing.T) {
		got, err := failingClient.GetById(testId)
		if err == nil {
			t.Errorf("Expected an error but got %+v", got)
		}
	})
	t.Run("It should return an error if no id is passed", func(t *testing.T) {
		got, err := mockClient.GetById("")
		if err == nil {
			t.Errorf("Expected an error but got %+v", got)
		}
	})

	t.Run("It should get a list of animes by name", func(t *testing.T) {
		parsedAnimeListResponse := []Anime{
			{
				Id:          "7785",
				Name:        "The Tatami Galaxy",
				Description: "...",
			},
			{
				Id:          "8985",
				Name:        "The Tatami Galaxy Specials",
				Description: "...",
			},
		}

		paginatedMockClient, server := NewMockPaginatedAnimeClient()
		defer server.Close()

		got, err := paginatedMockClient.ListByName("The Tatami Galaxy")

		expected := &parsedAnimeListResponse

		if err != nil {
			t.Errorf("Got unexpected error: %q", err)
		}

		if !reflect.DeepEqual(got, expected) {
			t.Errorf("Expected: %+v\n Got: %+v", expected, got)
		}
	})
}
