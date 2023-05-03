package parser

import (
	"regexp"
	"strings"
)

func FormatArgs(artist, song string) (string, string) {
	artist = strings.ReplaceAll(strings.ToUpper(string(artist[0]))+artist[1:], " ", "-")
	song = strings.ReplaceAll(strings.ToLower(song), " ", "-")

	return artist, song
}

func Title(title string) string {
	space, _ := regexp.Compile("\u200b")

	_, title, _ = strings.Cut(title, "<title>")
	title, _, _ = strings.Cut(title, "</title>")
	title = strings.TrimSuffix(title, " Lyrics | Genius Lyrics")
	title = strings.TrimSuffix(title, " | Genius")
	title = space.ReplaceAllString(title, "")

	return title
}

func Lyrics(lyrics string) string {
	data := strings.Split(lyrics, "data-lyrics-container=\"true\" ")
	breaks, _ := regexp.Compile(`<br/>`)
	bold, _ := regexp.Compile(`<b>(.+)<\/b>`)
	italic, _ := regexp.Compile(`<i>(.+)<\/i>`)
	tags, _ := regexp.Compile(`<\/*.+?>`)
	single, _ := regexp.Compile(`&#x27;`)
	double, _ := regexp.Compile(`&quot;`)
	amp, _ := regexp.Compile(`&amp;`)
	new, _ := regexp.Compile(`\n\n`)
	sec, _ := regexp.Compile(`\n\[`)
	sections := []string{}

	for i := 1; i < len(data); i++ {
		str, _, _ := strings.Cut(data[i], "</div><div class=\"RightSidebar__Container-pajcl2-0 jOFKJt\"")
		str = breaks.ReplaceAllString(str[45:], "\n")
		str = bold.ReplaceAllString(str, "\x1b[1m$1\x1b[0m")
		str = italic.ReplaceAllString(str, "\x1b[3m$1\x1b[0m")
		str = tags.ReplaceAllString(str, "")
		str = single.ReplaceAllString(str, "'")
		str = double.ReplaceAllString(str, "\"")
		str = amp.ReplaceAllString(str, "&")
		str = new.ReplaceAllString(str, "\n")
		str = sec.ReplaceAllString(str, "\n\n[")

		if i == len(data)-1 {
			str, _, _ = strings.Cut(str, "Embed")
		}

		sections = append(sections, str)
	}

	return strings.Join(sections, "\n")
}

func AlbumList(list string) []string {
	data := strings.Split(list, "<div class=\"chart_row-content\">")
	sections := []string{}

	for i := 1; i < len(data); i++ {
		str, _, _ := strings.Cut(data[i], "\" class=\"u-display_block\">")
		_, str, _ = strings.Cut(str, "href=\"")
		sections = append(sections, str)
	}

	return sections
}

func WordWrap(lyrics string, width int) string {
	lines := strings.Split(lyrics, "\n")
	out := []string{}
	for _, line := range lines {
		newLine := []string{}
		for i, char := range line {
			i += 1
			if i%width == 0 {
				newLine = append(newLine, "\n")
			}

			newLine = append(newLine, string(char))
		}

		out = append(out, newLine...)
		out = append(out, "\n")
	}

	return strings.Join(out, "")
}
