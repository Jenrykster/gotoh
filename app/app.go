package app

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/Jenrykster/gotoh/api"
)

const (
	startMessage       = "Insert the anime name: "
	endMessage         = "Selected anime:\nname: %s\ndescription: %s\n"
	zeroResultsMessage = "no animes found"
)

type AnimeData struct {
	Name        string
	Description string
}

func Init(w io.Writer, r io.Reader, animeAPI api.AnimeApi) (err error) {
	in := bufio.NewReader(r)

	fmt.Fprint(w, startMessage)
	userInput, err := in.ReadString('\n')
	fmt.Println()

	animeQuery := strings.ReplaceAll(userInput, "\n", "")
	if err != nil {
		fmt.Fprintf(w, "There was an error getting this anime: %s", err.Error())
	}

	// User is searching by ID
	if _, err := strconv.Atoi(animeQuery); err == nil {
		selectedAnime, err := getAnimeData(animeQuery, animeAPI)

		if err != nil {
			if err.Error() == api.NotFoundError {
				fmt.Fprint(w, err)
				return nil
			} else {
				fmt.Fprintf(w, "There was an error: %s ", err.Error())
			}
			return err
		}
		fmt.Fprintf(w, endMessage, selectedAnime.Name, selectedAnime.Description)
	} else {
		animeList, err := searchAnimes(animeQuery, animeAPI)

		if err != nil {
			if err.Error() == zeroResultsMessage {
				fmt.Fprint(w, err)
				return nil
			} else {
				fmt.Fprintf(w, "There was an error: %s ", err.Error())
			}
			return err
		}
		fmt.Fprint(w, strings.Join(animeList, "\n"))

	}

	return nil
}

func getAnimeData(id string, animeAPI api.AnimeApi) (anime *AnimeData, err error) {
	result, err := animeAPI.GetById(id)

	if err != nil {
		return nil, err
	}
	return &AnimeData{Name: result.Name, Description: result.Description}, nil
}

func searchAnimes(query string, animeAPI api.AnimeApi) (animeList []string, err error) {
	result, err := animeAPI.ListByName(query)

	if err != nil {
		return nil, err
	}

	if len(*result) < 1 {
		return nil, errors.New(zeroResultsMessage)
	}

	for _, animeName := range *result {
		animeList = append(animeList, animeName.Name)
	}

	return
}
