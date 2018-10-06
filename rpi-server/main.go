package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/dhowden/tag"
	"google.golang.org/appengine/log"
)

var shouldExit = make(chan bool)
var dstoreClient *datastore.Client
var err error
var states []lightShowStatePayload

func init() {
	// Create Context
	ctx := context.Background()
	// Get Project name from Credential file
	project, err := getProjectInfo()
	if err != nil {
		fmt.Printf("Error Reading Credential File")
		exitNow()
		panic(1)
	}
	// Init Datastore Client
	dstoreClient, err = datastore.NewClient(ctx, project)
	if err != nil {
		fmt.Printf("Error Getting Client: %v", err)
	}
}
func getSongInfo(path string) (string, string, error) {
	f, err := os.Open(os.ExpandEnv(path))
	if err != nil {
		return "", "", fmt.Errorf("Error Opening Path: %v", err)
	}
	defer f.Close()

	m, err := tag.ReadFrom(f)
	if err != nil {
		return "", "", fmt.Errorf("Error Reading from File: %v", err)
	}
	return m.Title(), m.Artist(), nil
}
func readState() lightShowStatePayload {
	fileContent, err := getPlaylist()
	if err != nil {
		fmt.Printf("%v", err)
	}
	currentSong, err := getCurrentSong()
	if err != nil {
		fmt.Printf("%v", err)
	}
	playlistRunning, err := getPlaylistStatus()
	if err != nil {
		fmt.Printf("%v", err)
	}
	pload := lightShowStatePayload{}
	songNames := []song{}
	rows := strings.Split(fileContent, "\n")
	for _, row := range rows {
		if row != "" {
			title, artist, err := getSongInfo(strings.Split(row, "\t")[1])
			if err != nil {
				fmt.Printf("%v", err)
			}
			songNames = append(songNames, song{Name: title, Artist: artist})
		}
	}
	pload.Songs = songNames
	pload.Running = playlistRunning
	songName, artist := getCurrentSongInfo(currentSong)
	pload.CurrentSong = song{Name: songName, Artist: artist}
	return pload
}
func getCurrentSongInfo(songinfo string) (string, string) {
	infoArray := strings.Split(songinfo, " by ")
	if len(infoArray) == 2 {
		return strings.Replace(infoArray[0], "Now Playing ", "", -1), infoArray[1]
	} else if len(infoArray) == 1 {
		return strings.Replace(infoArray[0], "Now Playing ", "", -1), ""
	}
	return "", ""
}
func getPlaylistStatus() (bool, error) {
	statusString, err := readFromFile("/tmp/show-running")
	if err != nil {
		return false, fmt.Errorf("Error Running Playlist Status: %v", err)
	}
	if statusString == "false" {
		return false, nil
	}
	return true, nil
}
func getCurrentSong() (string, error) {
	return readFromFile("/tmp/current_song")

}
func getPlaylist() (string, error) {
	return readFromFile(os.ExpandEnv("$SYNCHRONIZED_LIGHTS_HOME/music/sample/.playlist"))

}
func readFromFile(path string) (string, error) {
	fileContent, err := ioutil.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("Error Reading Playlist File: %v", err)

	}
	return string(fileContent), nil
}
func exitNow() {
	shouldExit <- true
}
func updateState(state lightShowStatePayload) {
	ctx := context.Background()

	b, err := json.Marshal(state)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	dpayload := datastorePayload{Value: b}
	stateKey := getStateKey()
	commit, err := dstoreClient.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		_, err = tx.Put(stateKey, &dpayload)
		if err != nil && err != datastore.ErrNoSuchEntity {
			return fmt.Errorf("Error Getting State from Datastore: %v", err)
		} else if err != nil {
			b := []byte("{\"songs\":[{\"name\":\"First Song\",\"artist\":\"First Artist\"},{\"name\":\"Second Song\",\"artist\":\"Second Artist\"},{\"name\":\"Third Song\",\"artist\":\"Third Artist\"}],\"running\":true,\"currentSongIndex\":0,\"currentSong\":{\"name\":\"First Song\",\"artist\":\"First Artist\"}}")
			dpayload.Value = b
			pendingKey, err := tx.Put(stateKey, &dpayload)
			if err != nil {
				return fmt.Errorf("Error Putting New Entity: %v", err)
			}
			log.Infof(ctx, "Pending Key: %+v", *pendingKey)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	if commit != nil {
		fmt.Printf("Commit: %v", commit)
	}
}
func readStateInterval() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Fatal Error: %v", r)
			os.Exit(1)
		}
	}()
	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		state := readState()
		updateState(state)
		shouldExit <- false
	}

}
func main() {

	go readStateInterval()
	for exitGo := range shouldExit {
		if exitGo {
			os.Exit(0)
		}
	}
}

func getProjectInfo() (string, error) {
	var projectInformation projectInfo
	plan, err := ioutil.ReadFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(plan, &projectInformation)
	if err != nil {
		return "", err
	}
	return projectInformation.Project, nil
}
