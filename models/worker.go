package models

import (
	"fmt"
	"log"
	"sync"

	redigo "github.com/garyburd/redigo/redis"
)

type Worker struct {
	cache Cache
	id    int
	queue string
}

func newWorker(id int, cache Cache, queue string) Worker {
	return Worker{cache: cache, id: id, queue: queue}
}

func (w Worker) process(id int) {
	for {
		conn := w.cache.Pool.Get()
		var channel string
		var imageURL string
		if reply, err := redigo.Values(conn.Do("BLPOP", w.queue, 30+id)); err == nil {

			if _, err := redigo.Scan(reply, &channel, &imageURL); err != nil {
				w.cache.EnqueueValue(w.queue, imageURL)
				continue
			}

			file, err := createFile(imageURL)
			if err != nil {
				w.cache.EnqueueValue(w.queue, imageURL)
				continue
			}

			fmt.Printf("id: %d Downloaded a file %v\n", id, file.Name)

		} else if err != redigo.ErrNil {
			log.Fatal(err)
		}
		conn.Close()
	}
}

func UsersToDB(numWorkers int, cache Cache, queue string) {
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(id int, cache Cache, queue string) {
			worker := newWorker(id, cache, queue)
			worker.process(id)
			defer wg.Done()
		}(i, cache, queue)
	}
	wg.Wait()
}
