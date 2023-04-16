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
		selectedId, err := selectAnimeFromList(w, r, animeAPI, animeQuery)
		if err != nil {
			if err.Error() == zeroResultsMessage {
				fmt.Fprint(w, err)
				return nil
			} else {
				fmt.Fprintf(w, "There was an error: %s ", err.Error())
			}
			return err
		}
		selectedAnime, err := getAnimeData(selectedId, animeAPI)
		if err != nil {
			fmt.Fprintf(w, "Error: %s", err)
		}
		fmt.Fprintf(w, endMessage, selectedAnime.Name, selectedAnime.Description)
	}

	return nil
}

func getAnimeData(id string, animeAPI api.AnimeApi) (anime *api.Anime, err error) {
	result, err := animeAPI.GetById(id)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func selectAnimeFromList(w io.Writer, r io.Reader, animeAPI api.AnimeApi, query string) (selectionId string, err error) {
	animeList, err := searchAnimes(query, animeAPI)
	orderedAnimeList := []string{}

	if err != nil {
		return "", fmt.Errorf("couldn't display list: %w", err)
	}

	for i, animeName := range *animeList {
		orderedAnimeList = append(orderedAnimeList, fmt.Sprintf("[%d] %s", i+1, animeName.Name))
	}
	fmt.Fprint(w, strings.Join(orderedAnimeList, "\n"))

	in := bufio.NewReader(r)
	userInput := ""
	selectedNumber := -1

	for selectedNumber < 0 || selectedNumber > len(orderedAnimeList)+1 {
		fmt.Println()
		fmt.Fprint(w, "Type the number of the anime and press ENTER: ")
		userInput, err = in.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("couldn't read user input: %w", err)
		}
		selectedNumber, err = strconv.Atoi(strings.Split(userInput, "\n")[0])
		fmt.Println()
	}

	if err != nil {
		return
	}

	selectedAnime := (*animeList)[selectedNumber-1]
	return string(selectedAnime.Id), nil
}

func searchAnimes(query string, animeAPI api.AnimeApi) (animeList *[]api.Anime, err error) {
	result, err := animeAPI.ListByName(query)

	if err != nil {
		return nil, err
	}

	if len(*result) < 1 {
		return nil, errors.New(zeroResultsMessage)
	}

	return result, nil
}
