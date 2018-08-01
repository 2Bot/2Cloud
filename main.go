package main

import (
	"os"

	"github.com/fsouza/go-dockerclient"
)

func main() {
	// creates a chi router
	//r := chi.NewRouter()

	//r.Route("/{userID}", func(r chi.Router) {
	// gets data for {userID}'s instance
	//r.Get("/", getData)
	// creates an instance for {userID}
	//r.Post("/create", createContainer)
	// updates {userID}'s instance
	//r.Post("/update", updateContainer)
	// updates the settings for {userID}'s instance. Settings sent in body
	//.Put("/settings", updateSettings)
	//})
	createContainer("100")
}

func getData() {}

func createContainer(userID string) {
	endpoint := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(endpoint)
	if err != nil {
		panic(err)
	}
	containers, err := client.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		panic(err)
	}
	if containers != nil {
		for i := 0; i < len(containers); i++ {
			if containers[i].Names[0] == "/"+userID {
				return
			}
		}
	}
	err = os.Mkdir("./users/"+userID, 0777)
	if err != nil {
		panic(err)
	}
	_, err = client.CreateContainer(docker.CreateContainerOptions{
		Name: userID,
		Config: &docker.Config{
			AttachStdout: true,
			AttachStdin:  false,
			Cmd:          []string{"/bin/bash"},
			Image:        "library/hello-world",
		},
	})
	if err != nil {
		panic(err)
	}
}

func updateContainer() {}

func updateSettings() {}
