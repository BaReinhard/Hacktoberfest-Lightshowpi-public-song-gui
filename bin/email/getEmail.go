package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type projectInfo struct {
	Project     string `json:"project_id"`
	ClientEmail string `json:"client_email"`
}

func main() {
	_, email, err := getProjectInfo()
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
	fmt.Printf("%s", email)
}
func getProjectInfo() (string, string, error) {
	var projectInformation projectInfo
	plan, err := ioutil.ReadFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	if err != nil {
		return "", "", err
	}
	err = json.Unmarshal(plan, &projectInformation)
	if err != nil {
		return "", "", err
	}
	return projectInformation.Project, projectInformation.ClientEmail, nil
}
