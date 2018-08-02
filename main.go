package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/fsouza/go-dockerclient"
	"github.com/go-chi/render"
)

type postStruct struct {
	UserID string `json:"userID"`
	Token  string `json:"token"`
	Prefix string `json:"prefix"`
	Space  bool   `json:"space"`
}

type response struct {
	Message string `json:"message"`
}

func (e *response) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func res(r string) render.Renderer {

	return &response{Message: r}
}

const endpoint = "unix:///var/run/docker.sock"

func main() {
	// creates a chi router
	//r := chi.NewRouter()
	fmt.Println("Server is running")
	http.HandleFunc("/", createContainer)
	http.ListenAndServe(":8080", nil)
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
	//})
}

func getData() {}

func createContainer(w http.ResponseWriter, r *http.Request) {
	var data postStruct
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.WriteHeader(500)
		render.Render(w, r, res("Oopsiedoopsie our server had a little fuckywucky"))
		return
	}

	client, err := docker.NewClient(endpoint)
	if err != nil {
		w.WriteHeader(500)
		render.Render(w, r, res("Oopsiedoopsie our server had a little fuckywucky"))
		return
	}
	_, err = containerExists(data.UserID)
	if err != nil {
		w.WriteHeader(409)
		render.Render(w, r, res("Fuck you greedy scum. Only one container for you"))
		return
	}
	path := filepath.Join("./users", data.UserID)
	err = os.MkdirAll(path, 0777)
	if err != nil {
		w.WriteHeader(500)
		render.Render(w, r, res("Oopsiedoopsie our server had a little fuckywucky"))
		return
	}

	_, err = client.CreateContainer(docker.CreateContainerOptions{
		Name: data.UserID,
		Config: &docker.Config{
			AttachStdout: true,
			AttachStdin:  false,
			Cmd:          []string{"/bin/bash"},
			Image:        "library/hello-world",
		},
	})
	if err != nil {
		w.WriteHeader(500)
		render.Render(w, r, res("Oopsiedoopsie our server had a little fuckywucky"))
		return
	}
}

func containerExists(userID string) (bool, error) {
	client, err := docker.NewClient(endpoint)
	if err != nil {
		return false, err
	}
	containers, err := client.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		return false, err
	}

	if containers != nil {
		for i := 0; i < len(containers); i++ {
			if containers[i].Names[0] == "/"+userID {
				return true, err
			}
		}
	}
	return false, err
}
