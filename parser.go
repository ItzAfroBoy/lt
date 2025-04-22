package main

import (
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
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

func formatSpotify(_artist, _title string) {
	*artist = strings.ReplaceAll(strings.ToUpper(string((_artist)[0]))+(_artist)[1:], " ", "-")
	*title = strings.ReplaceAll(strings.ToLower(_title), " ", "-")
}

func parseLyrics(lyrics string) string {
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(lyrics))
	secs := doc.Find("div.Lyrics__Container-sc-78fb6627-1.hiRbsH")
	secs.Each(func(i int, s *goquery.Selection) {
		s.Find("div").Remove()
		s.Find("br").ReplaceWithHtml("\n")
		bold := s.Find("b")
		italics := s.Find("i")

		bold.Each(func(i int, s *goquery.Selection) {
			s.ReplaceWithHtml("\x1b[1m" + s.Text() + "\x1b[0m")
		})

		italics.Each(func(i int, s *goquery.Selection) {
			s.ReplaceWithHtml("\x1b[2m" + s.Text() + "\x1b[0m")
		})
		
		if i != secs.Size()-1 && i%2 == 0 {
			s.SetText(s.Text() + "\n")
		}
	})
	return secs.Text()
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
