package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hasura/go-graphql-client"
)

type AnimeEntityData struct {
	Id    int
	Title struct {
		English string
		Native  string
	}
	Description string
}
type AnilistAnime struct {
	Media AnimeEntityData `graphql:"Media(id: $id)"`
}

type QueryParams map[string]interface{}

type PaginatedSuccessResponse struct {
	Page struct {
		Media []AnimeEntityData `graphql:"media(search: $name, type: ANIME)"`
	} `graphql:"Page(perPage: 5)"`
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
	client graphql.Client
	url    string
}

const (
	NotFoundError = "Anime not found"
)

func newCustomAnimeClient(client *http.Client, url string) AnimeClient {
	graphqlClient := graphql.NewClient(url, client)

	return AnimeClient{
		client: *graphqlClient,
		url:    url,
	}
}

func NewAnilistClient() AnimeClient {
	httpClient := &http.Client{Timeout: 10 * time.Second}

	return newCustomAnimeClient(httpClient, "https://graphql.anilist.co")
}

func (a *AnimeClient) GetById(id string) (anime *Anime, err error) {
	if len(id) < 1 {
		return nil, errors.New("no id provided")
	}

	convertedId, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("cannot query id: %w", err)
	}

	var response AnilistAnime
	params := QueryParams{"id": convertedId}
	error := a.client.Query(context.Background(), &response, params)
	if error != nil {
		return nil, error
	}

	selectedAnime := convertAnilistResponse(response.Media)

	return &selectedAnime, nil
}

func (a *AnimeClient) ListByName(name string) (*[]Anime, error) {
	if len(name) < 1 {
		return nil, errors.New("no name provided")
	}

	var parsedResponse PaginatedSuccessResponse

	params := QueryParams{"name": name}
	error := a.client.Query(context.Background(), &parsedResponse, params)
	if error != nil {
		return nil, error
	}

	animeList := []Anime{}

	for _, e := range parsedResponse.Page.Media {
		animeList = append(animeList, convertAnilistResponse(e))
	}

	return &animeList, nil
}

func convertAnilistResponse(raw AnimeEntityData) Anime {
	var title string

	if len(strings.Trim(raw.Title.English, " ")) == 0 {
		title = raw.Title.Native
	} else {
		title = raw.Title.English
	}

	return Anime{Id: fmt.Sprint(raw.Id), Name: title, Description: raw.Description}
}
