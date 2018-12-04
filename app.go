package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/phaicom/async-image-downloader/models"
	"github.com/urfave/negroni"
)

var (
	fullURL string
	images  []models.Image
)

type App struct {
	Router *mux.Router
	Cache  models.Cache
}

func (a *App) Initialize(cache models.Cache) {
	a.Cache = cache
	a.Router = mux.NewRouter()
}

func (a *App) Run(addr string) {
	n := negroni.Classic()
	n.UseHandler(a.Router)
	log.Fatal(http.ListenAndServe(addr, n))
}

func (a *App) DownloadImage(URL string) {
	fullURL = URL
	client := httpClient()
	res, err := client.Get(fullURL)
	checkError(err)

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	checkError(err)

	err = json.Unmarshal([]byte(body), &images)
	checkError(err)

	for _, image := range images {
		a.Cache.EnqueueValue(createImageQueue, image.ThumbnailURL)
	}
}

func httpClient() *http.Client {
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	return &client
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
