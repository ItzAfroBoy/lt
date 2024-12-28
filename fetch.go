package main

import (
	"fmt"
	"io"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
)

type resMsg struct {
	title  string
	lyrics string
}

type albumMsg struct {
	titles []string
	lyrics []string
}

func get(url string) (string, string) {
	res, err := http.Get(url)
	if err != nil {
		return "Error", err.Error()
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "Error", err.Error()
	}

	data := string(body)
	title := formatTitle(data)
	lyrics := parseLyrics(data)

	return title, lyrics
}

func getSong() tea.Msg {
	url := fmt.Sprintf("https://genius.com/%s-%s-lyrics", *artist, *title)
	title, lyrics := get(url)
	if title == "" && lyrics == "" {
		return resMsg{"", ""}
	}

	return resMsg{title, lyrics}
}

func getAlbum() tea.Msg {
	res, err := http.Get(fmt.Sprintf("https://genius.com/albums/%s/%s", *artist, *title))
	if err != nil {
		return resMsg{"Error", err.Error()}
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return resMsg{"Error", err.Error()}
	}

	tracklist := albumList(string(body))
	var titles, lyrics []string

	for _, song := range tracklist {
		title, lyric := get(song)
		if title == "" && lyric == "" {
			continue
		}
		titles = append(titles, title)
		lyrics = append(lyrics, lyric)
	}

	return albumMsg{titles, lyrics}
}
