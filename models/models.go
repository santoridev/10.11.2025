package models

import (
	"sync"
)

type Req struct {
	Links []string `json:"links"`
}

type RespData struct {
	Links    map[string]string `json:"links"`
	LinksNum int               `json:"links_num"`
}

type PDFReq struct {
	LinksNum []int `json:"links_num"`
}

var (
	LinksStorage = make(map[int]map[string]string) // наш storage можно будет заменить на бд, чтобы хранить данные после перезапуска
	Mu           sync.Mutex
)

type Task struct {
	Links      []string
	ResultChan chan map[string]string
}

var (
	TaskQueue = make(chan Task, 100)
	wg        sync.WaitGroup
	stop      = make(chan struct{})
)

func StartWorkerPool(n int) {
	for i := 0; i < n; i++ {
		wg.Add(1)
		go worker(i)
	}
}

func StopWorkerPool() {
	close(stop)
	wg.Wait()
}

func worker(id int) {
	defer wg.Done()
	for {
		select {
		case task := <-TaskQueue:
			//log.Printf("[%d]  %d links", id, len(task.Links))

			results := make(map[string]string)
			for _, link := range task.Links {
				results[link] = CheckURL(link)
			}

			task.ResultChan <- results

			//log.Printf("[%d] finished", id)

		case <-stop:
			//log.Printf("[%d] stopped  ", id)
			return
		}
	}
}
