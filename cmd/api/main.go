package main

import (
	"fmt"
	"net/http"
	"santori/linkchecker/handlers"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlers.Check(w, r)
		default:
			fmt.Println("use different method to send links")
		}
	})
	mux.HandleFunc("/report", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlers.GenerateReport(w, r)
		default:
			fmt.Println("use different method to send links")
		}
	})

	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Println(err)
	}
}
