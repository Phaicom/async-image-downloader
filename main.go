package main

import (
	"flag"

	"github.com/phaicom/async-image-downloader/models"
)

const (
	createImageQueue = "CREATE_IMAGE"
)

func main() {
	numWorkers := 2
	cache := models.Cache{Enable: true}

	flag.StringVar(
		&cache.Address,
		"redis_address",
		"127.0.0.1:6379",
		"Redis Address",
	)

	flag.StringVar(
		&cache.Auth,
		"redis_auth",
		"",
		"Redis Auth",
	)

	flag.StringVar(
		&cache.DB,
		"redis_db_name",
		"0",
		"Redis DB name",
	)

	flag.IntVar(
		&cache.MaxIdle,
		"redis_max_idle",
		10,
		"Redis Max Idle",
	)

	flag.IntVar(
		&cache.MaxActive,
		"redis_max_active",
		100,
		"Redis Max Active",
	)

	flag.IntVar(
		&cache.IdleTimeoutSecs,
		"redis_timeout",
		60,
		"Redis timeout in seconds",
	)

	flag.Parse()

	cache.Pool = cache.NewCachePool()
	go models.CreateImageQueue(numWorkers, cache, createImageQueue)

	app := App{}
	app.Initialize(cache)
	app.DownloadImage("https://jsonplaceholder.typicode.com/photos")
	app.Run(":3000")
}
