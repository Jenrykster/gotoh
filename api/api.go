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
	Id       int
	Episodes int
	Title    struct {
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
	Episodes    int
}

type AnimeApi interface {
	ListByName(name string) (*[]Anime, error)
	GetById(id string) (*Anime, error)
	UpdateEpisodeCount(id string, episode int, markAsComplete bool) error
}

type AnimeClient struct {
	client graphql.Client
	url    string
}

const (
	NotFoundError = "Anime not found"
)

type AuthRoundTripper struct {
	token   string
	Proxied http.RoundTripper
}

func (art AuthRoundTripper) RoundTrip(req *http.Request) (res *http.Response, e error) {
	req.Header.Add("Authorization", "Bearer "+art.token)
	res, e = art.Proxied.RoundTrip(req)
	return
}

func newCustomAnimeClient(client *http.Client, url string) AnimeClient {
	graphqlClient := graphql.NewClient(url, client)

	return AnimeClient{
		client: *graphqlClient,
		url:    url,
	}
}

func NewAnilistClient(userToken string) AnimeClient {
	transporter := AuthRoundTripper{
		token:   userToken,
		Proxied: http.DefaultTransport,
	}
	httpClient := &http.Client{Timeout: 10 * time.Second, Transport: transporter}

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

type MediaListStatus string

const (
	CURRENT   MediaListStatus = "CURRENT"
	COMPLETED MediaListStatus = "COMPLETED"
)

type UpdateEpisode struct {
	SaveMediaListEntry struct {
		Id       int
		Status   MediaListStatus
		Progress int
	} `graphql:"SaveMediaListEntry(mediaId: $id, status: $status, progress: $progress)"`
}

func (a *AnimeClient) UpdateEpisodeCount(id string, episode int, markAsCompleted bool) error {
	convertedId, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("cannot query id: %w", err)
	}

	var mutation UpdateEpisode

	status := CURRENT
	if markAsCompleted {
		status = COMPLETED
	}

	params := QueryParams{
		"id":       convertedId,
		"progress": episode,
		"status":   status,
	}
	err = a.client.Mutate(context.Background(), &mutation, params)

	if err != nil {
		return err
	}

	return nil
}

func convertAnilistResponse(raw AnimeEntityData) Anime {
	var title string

	if len(strings.Trim(raw.Title.English, " ")) == 0 {
		title = raw.Title.Native
	} else {
		title = raw.Title.English
	}

	description := strings.ReplaceAll(raw.Description, "<br>", "")

	return Anime{Id: fmt.Sprint(raw.Id), Name: title, Description: description, Episodes: raw.Episodes}
}
