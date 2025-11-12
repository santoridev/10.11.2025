package handlers

import (
	"encoding/json"
	"net/http"
	"santori/linkchecker/models"
)

func Check(w http.ResponseWriter, r *http.Request) {
	var req models.Req
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resultChan := make(chan map[string]string)
	task := models.Task{
		Links:      req.Links,
		ResultChan: resultChan,
	}

	models.TaskQueue <- task

	result := <-resultChan

	models.Mu.Lock()
	index := len(models.LinksStorage)
	models.LinksStorage[index] = result
	linksNum := len(models.LinksStorage)
	models.Mu.Unlock()

	resp := models.RespData{
		Links:    result,
		LinksNum: linksNum,
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
