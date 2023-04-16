package api

import (
	"net/http"
	"net/http/httptest"
)

const getByIdResponse = `{
	"data": {
		"Media": {
			"id": 7785,
			"title": {
				"english": "The Tatami Galaxy"	
			},
			"description": "..."
		}
	}
}`

func NewMockAnimeClient() (AnimeClient, *httptest.Server) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(getByIdResponse))
	}))

	return newCustomAnimeClient(server.Client(), server.URL), server
}

const getListByNameResponse = `{
	"data": {
		"Page": {
			"media": [
				{
					"id": 7785,
					"title": {
						"english": "The Tatami Galaxy"
					},
					"description": "..."
				},
				{
					"id": 8985,
					"title": {
						"english": "The Tatami Galaxy Specials"
					},
					"description": "..."
				}
			]
		}
	}
}`

func NewMockPaginatedAnimeClient() (AnimeClient, *httptest.Server) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(getListByNameResponse))
	}))

	return newCustomAnimeClient(server.Client(), server.URL), server
}
