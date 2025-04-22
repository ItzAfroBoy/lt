package main

import (
	gospotti "github.com/ItzAfroBoy/go-spotti"
)

func getSpotifyInfo() (title, artist string) {
	spotti := gospotti.Init()
	spotti.Auth.RedirectURI = "http://localhost:7171/callback"
	spotti.ClientID = "f1b6295487874fafb175fb5818c5abcf"
	spotti.Authorize(false)

	info := spotti.Playback.GetPlaybackInfo()
	title = info.Track.Name
	artist = info.Track.Artists[0].Name

	formatSpotify(artist, title)
	return
}
