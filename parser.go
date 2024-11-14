package main

import (
	"os"
	"regexp"
	"strings"
)

func formatArgs() {
	*artist = strings.ReplaceAll(strings.ToUpper(string((*artist)[0]))+(*artist)[1:], " ", "-")
	*title = strings.ReplaceAll(strings.ToLower(*title), " ", "-")
}

func formatTitle(title string) string {
	space, _ := regexp.Compile("\u200b")
	titleHTML, _ := regexp.Compile(`<title>(.+)( Lyrics.+)<\/title>`)
	defer func() {
		if r := recover(); r != nil {
			p.Quit()
		}
	}()
	title = titleHTML.FindAllStringSubmatch(title, -1)[0][1]
	title = space.ReplaceAllString(title, "")
	return title
}

func parseLyrics(lyrics string) string {
	data := strings.Split(lyrics, "data-lyrics-container=\"true\" ")
	breaks, _ := regexp.Compile(`<br/>`)
	bold, _ := regexp.Compile(`<b>([\s\S]+?)<\/b>`)
	italic, _ := regexp.Compile(`<i>([\s\S]+?)<\/i>`)
	tags, _ := regexp.Compile(`<\/*.+?>`)
	single, _ := regexp.Compile(`&#x27;`)
	double, _ := regexp.Compile(`&quot;`)
	amp, _ := regexp.Compile(`&amp;`)
	div, _ := regexp.Compile(`.*<div `)
	sec, _ := regexp.Compile(`>\[`)
	sections := []string{}

	for i := 1; i < len(data); i++ {
		str, _, _ := strings.Cut(data[i], "</div><div class=\"RightSidebar__Container-pajcl2-0 jOFKJt\"")
		str = breaks.ReplaceAllString(str[46:], "\n")
		str = bold.ReplaceAllString(str, "\x1b[1m$1\x1b[0m")
		str = italic.ReplaceAllString(str, "\x1b[2m$1\x1b[0m")
		str = tags.ReplaceAllString(str, "")
		str = single.ReplaceAllString(str, "'")
		str = double.ReplaceAllString(str, "\"")
		str = amp.ReplaceAllString(str, "&")
		str = div.ReplaceAllString(str, "")
		str = sec.ReplaceAllString(str, "\n[")

		if i == len(data)-1 {
			embed, _ := regexp.Compile(`\d+Embed`)
			str = embed.ReplaceAllString(str, "Embed")
			str, _, _ = strings.Cut(str, "Embed")
		}

		sections = append(sections, str)
	}

	return strings.Join(sections, "\n")
}

func parseFile(file string) (title, content string) {
	rawFile, _ := os.ReadFile(file)
	parsedFile := string(rawFile)
	title, content, _ = strings.Cut(parsedFile, "\n\n")
	return
}

func albumList(list string) []string {
	data := strings.Split(list, "<div class=\"chart_row-content\">")
	sections := []string{}

	for i := 1; i < len(data); i++ {
		str, _, _ := strings.Cut(data[i], "\" class=\"u-display_block\">")
		_, str, _ = strings.Cut(str, "href=\"")
		sections = append(sections, str)
	}

	return sections
}

func (m model) wordWrap() string {
	lines := strings.Split(m.content, "\n")
	out := []string{}
	for _, line := range lines {
		newLine := []string{}
		for i, char := range line {
			i++
			if i%m.viewport.Width == 0 {
				newLine = append(newLine, "\n")
			}
			newLine = append(newLine, string(char))
		}
		out = append(out, newLine...)
		out = append(out, "\n")
	}

	return strings.Join(out, "")
}

func userHomeDir() string {
	dir, _ := os.UserHomeDir()
	return dir
}
