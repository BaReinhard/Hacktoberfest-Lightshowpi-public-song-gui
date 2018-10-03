package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/datastore"
	"github.com/gorilla/mux"
	"google.golang.org/appengine"
	// "google.golang.org/appengine/datastore"

	"google.golang.org/appengine/log"
)

var dstoreClient *datastore.Client
var err error

func main() {
	r := mux.NewRouter()
	fs := http.FileServer(http.Dir("static"))
	// r.HandleFunc("/",redirectToProfile)
	r.Handle("/", fs)
	r.HandleFunc("/api/getState", getStateHandler)
	http.Handle("/", r)
	appengine.Main()
}

func getStateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	dstoreClient, err = datastore.NewClient(ctx, os.Getenv("PROJECT_ID"))
	if err != nil {
		log.Errorf(ctx, "Error Getting New Datastore Client: %v", err)
		json.NewEncoder(w).Encode(errorPayload{Error: err.Error(), StatusCode: http.StatusInternalServerError})
		return
	}
	statePayload := &lightShowStatePayload{}
	dpayload := datastorePayload{}
	stateKey := getStateKey()
	commit, err := dstoreClient.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		err = tx.Get(stateKey, &dpayload)
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
	if err != nil && err == datastore.ErrConcurrentTransaction {
		log.Errorf(ctx, "Transaction Error: %v", err)
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(errorPayload{Error: err.Error(), StatusCode: http.StatusTooManyRequests})
		return
	} else if err != nil {
		log.Errorf(ctx, "Transaction Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorPayload{Error: err.Error(), StatusCode: http.StatusInternalServerError})
		return
	}
	if commit != nil {
		log.Infof(ctx, "Commit: %+v", *commit)
	}
	err = json.Unmarshal(dpayload.Value, statePayload)
	if err != nil {
		log.Errorf(ctx, "Unable to Unmarshall Datastore Payload: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorPayload{Error: err.Error(), StatusCode: http.StatusInternalServerError})
		return
	}
	json.NewEncoder(w).Encode(statePayload)
	return
}
