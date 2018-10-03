package main

type lightShowStatePayload struct {
	Songs            []song `json:"songs"`
	Running          bool   `json:"running"`
	CurrentSongIndex int    `json:"currentSongIndex"`
	CurrentSong      song   `json:"currentSong"`
}

type song struct {
	Name   string `json:"name"`
	Artist string `json:"artist"`
}

type errorPayload struct {
	Error      string `json:"error"`
	StatusCode int    `json:"statusCode"`
}

type datastorePayload struct {
	Value []byte
}
type projectInfo struct {
	Project     string `json:"project_id"`
	ClientEmail string `json:"client_email"`
}
