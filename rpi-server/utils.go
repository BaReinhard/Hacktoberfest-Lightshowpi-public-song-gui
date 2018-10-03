package main

import "cloud.google.com/go/datastore"

func getStateKey() *datastore.Key {
	return datastore.NameKey("State", "state", nil)
}
