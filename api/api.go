package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type AnilistAnime struct {
	Id    int `json:"id"`
	Title struct {
		English string `json:"english"`
	} `json:"title"`

	Description string `json:"description"`
}

type SuccessResponse struct {
	Data struct {
		Media AnilistAnime `json:"Media"`
	} `json:"data"`
}

type PaginatedSuccessResponse struct {
	Data struct {
		Page struct {
			Media []AnilistAnime `json:"media"`
		} `json:"Page"`
	} `json:"data"`
}

type Anime struct {
	Id          string
	Name        string
	Description string
}

type AnimeApi interface {
	ListByName(name string) (*[]Anime, error)
	GetById(id string) (*Anime, error)
}

type AnimeClient struct {
	client *http.Client
	url    string
}

const (
	NotFoundError = "Anime not found"
)

type jsonQuery map[string]string

func newCustomAnimeClient(client *http.Client, url string) AnimeClient {
	return AnimeClient{
		client: client,
		url:    url,
	}
}

func NewAnilistClient() AnimeClient {
	httpClient := &http.Client{Timeout: 10 * time.Second}

	return newCustomAnimeClient(httpClient, "https://graphql.anilist.co")
}

func (a *AnimeClient) getResource(query jsonQuery, v any) error {
	jsonValue, _ := json.Marshal(query)

	response, err := a.client.Post(a.url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil || response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusNotFound {
			return fmt.Errorf(NotFoundError)
		}
		return fmt.Errorf("there was an error during the request, %v", response.StatusCode)
	}

	data, err := io.ReadAll(response.Body)
	defer response.Body.Close()

	if err != nil {
		return errors.New("couldn't read the response")
	}

	unMarshalError := json.Unmarshal(data, v)

	if unMarshalError != nil {
		return errors.New("error during unmarshal process")
	}

	return nil
}

func (a *AnimeClient) GetById(id string) (anime *Anime, err error) {
	if len(id) < 1 {
		return nil, errors.New("no id provided")
	}

	query := jsonQuery{
		"query":     "query($id: Int){Media(id: $id, type: ANIME){id title{english} description}}",
		"variables": fmt.Sprintf(`{"id": %s}`, id),
	}

	var parsedResponse SuccessResponse
	error := a.getResource(query, &parsedResponse)

	if error != nil {
		return nil, error
	}

	selectedAnime := convertAnilistResponse(parsedResponse.Data.Media)

	return &selectedAnime, nil
}

func (a *AnimeClient) ListByName(name string) (*[]Anime, error) {
	if len(name) < 1 {
		return nil, errors.New("no name provided")
	}

	query := jsonQuery{
		"query":     "query($name: String){Page(perPage: 5){media(search: $name, type: ANIME){id title{english} description}}}",
		"variables": fmt.Sprintf(`{"name": "%s"}`, name),
	}

	var parsedResponse PaginatedSuccessResponse
	error := a.getResource(query, &parsedResponse)

	if error != nil {
		return nil, error
	}

	animeList := []Anime{}

	for _, e := range parsedResponse.Data.Page.Media {
		animeList = append(animeList, convertAnilistResponse(e))
	}

	return &animeList, nil
}

func convertAnilistResponse(raw AnilistAnime) Anime {
	return Anime{Id: fmt.Sprint(raw.Id), Name: raw.Title.English, Description: raw.Description}
}
