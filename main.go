package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/fsouza/go-dockerclient"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type postStruct struct {
	UserID string `json:"userID"`
	Token  string `json:"token"`
	Prefix string `json:"prefix"`
	Space  bool   `json:"space"`
}

type getStruct struct {
	UserID string `json:"userID"`
}

type response struct {
	Message string `json:"message"`
	Yes     string `json:"err"`
}

type getResponse struct {
	Running bool   `json:"running"`
	Version string `json:"version"`
}

func (e *response) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func res(r string, err error) render.Renderer {
	return &response{Message: r, Yes: err.Error()}
}

func (e *getResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func statsResponse(run bool, vers string) render.Renderer {
	return &getResponse{
		Running: run,
		Version: vers,
	}
}

const endpoint = "unix:///var/run/docker.sock"

func main() {
	// creates a chi router
	r := chi.NewRouter()
	fmt.Println("Server is running")
	r.Route("/", func(r chi.Router) {
		// gets data for {userID}'s instance
		r.Get("/", getData)
		// creates an instance for {userID}
		r.Post("/create", createContainer)
		// updates {userID}'s instance
		//r.Post("/update", updateContainer)
		// updates the settings for {userID}'s instance. Settings sent in body
		//.Put("/settings", updateSettings)
	})
	http.ListenAndServe(":8080", r)
}

func getData(w http.ResponseWriter, r *http.Request) {
	var data getStruct
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.WriteHeader(500)
		render.Render(w, r, res("There was an error creating a decoder :(", err))
		return
	}
	client, err := docker.NewClient(endpoint)
	if err != nil {
		w.WriteHeader(500)
		render.Render(w, r, res("There was an error creating an endpoint :(", err))
		return
	}
	contains, err := containerExists(*client, data.UserID)
	if contains {
		stats, err := client.InspectContainer(data.UserID)
		if err != nil {
			w.WriteHeader(500)
			render.Render(w, r, res("There was an error inspecting your container :(", err))
			return
		}

		render.Render(w, r, statsResponse(stats.State.Running, stats.Image))
		return
	}
	if err != nil {
		w.WriteHeader(500)
		render.Render(w, r, res("There was an error trying to find your container :(", err))
		return
	}
	w.WriteHeader(400)
	render.Render(w, r, res("Doesn't look like you have a container. Sorry :(", err))
	return
}

func createContainer(w http.ResponseWriter, r *http.Request) {
	var data postStruct
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.WriteHeader(500)
		render.Render(w, r, res("There was an error getting creating a new decoder :(", err))
		return
	}

	client, err := docker.NewClient(endpoint)
	if err != nil {
		w.WriteHeader(500)
		render.Render(w, r, res("There was an error creating an endpoint :(", err))
		return
	}
	contains, err := containerExists(*client, data.UserID)
	if contains {
		w.WriteHeader(409)
		render.Render(w, r, res("Fuck you greedy scum. Only one container for you", err))
		return
	}
	if err != nil {
		w.WriteHeader(500)
		render.Render(w, r, res("There was an error checking if you are greedy :(", err))
		return
	}
	path := filepath.Join("./users", data.UserID)
	err = os.MkdirAll(path, 0777)
	if err != nil {
		w.WriteHeader(500)
		render.Render(w, r, res("There was an error creating your directory", err))
		return
	}

	_, err = client.CreateContainer(docker.CreateContainerOptions{
		Name: data.UserID,
		Config: &docker.Config{
			AttachStdout: false,
			AttachStdin:  false,
			Image:        "2bot2go:3.0",
			Volumes: map[string]struct{}{
				"../emoji:/go/emoji":                                  {},
				"./users/{UserID}/config.toml:/go/config/config.toml": {},
			},
		},
	})
	if err != nil {
		w.WriteHeader(500)
		render.Render(w, r, res("There was an error creating your container :(", err))
		return
	}
}

func containerExists(client docker.Client, userID string) (bool, error) {
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
