package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"cloud.google.com/go/datastore"
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
		panic(1)
	}
	// Init Datastore Client
	dstoreClient, err = datastore.NewClient(ctx, project)
	if err != nil {
		fmt.Printf("Error Getting Client: %v", err)
	}
}
func readState() lightShowStatePayload {
	pload := lightShowStatePayload{}
	pload.Running = true
	pload.CurrentSong = song{Name: "Holla", Artist: "Back"}
	pload.Songs = []song{pload.CurrentSong}
	return pload
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