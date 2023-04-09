package app

import (
	"bufio"
	"fmt"
	"io"
)

const (
	startMessage = "Insert the anime name: "
	endMessage   = "Selected anime: %s"
)

func Init(w io.Writer, r io.Reader) (err error) {
	in := bufio.NewReader(r)

	fmt.Fprint(w, startMessage)

	selectedAnime, err := in.ReadString('\n')

	if err != nil {
		fmt.Fprintf(w, "There was an error: %s", err.Error())
		return err
	}

	fmt.Fprintf(w, endMessage, selectedAnime)
	return nil
}
