package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"santori/linkchecker/models"
	"strings"
	"sync"
	"time"
)

func Check(w http.ResponseWriter, r *http.Request) {
	var req models.Req

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	result := make(map[string]string)

	for _, url := range req.Links {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			status := checkURL(u, 3*time.Second)

			mu.Lock()
			result[u] = status
			mu.Unlock()
		}(url)
	}

	wg.Wait()

	mu.Lock()
	index := len(models.LinksStorage)
	if models.LinksStorage[index] == nil {
		models.LinksStorage[index] = make(map[string]string)
	}
	for url, status := range result {
		models.LinksStorage[index][url] = status
	}
	mu.Unlock()

	fmt.Println(models.LinksStorage)

	newResponse := models.RespData{
		Links:    result,
		LinksNum: len(models.LinksStorage),
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(newResponse); err != nil {
		return
	}
}

func checkURL(url string, timeout time.Duration) string {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "not_available"
	}

	client := &http.Client{Timeout: timeout}

	resp, err := client.Do(req)
	if err != nil {
		return "not_available"
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return "available"
	}
	return "not_available"
}
