package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	key, base := os.Getenv("INFRAI_API_KEY"), os.Getenv("APP_BASE_URL")
	if key == "" || base == "" {
		log.Fatal("set INFRAI_API_KEY and APP_BASE_URL")
	}
	client := oauthClient{HTTP: &http.Client{Timeout: 10 * time.Second}, Key: key}
	http.HandleFunc("GET /login", func(w http.ResponseWriter, r *http.Request) {
		input := loginRequest{Provider: r.URL.Query().Get("provider"), Workflow: r.URL.Query().Get("workflow"), Campaign: r.URL.Query().Get("campaign")}
		destination, err := returnTo(base, input)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		link, err := client.authorizeURL(r.Context(), input.Provider, destination)
		if err != nil {
			var apiErr *apiError
			if errors.As(err, &apiErr) && apiErr.Status >= 400 && apiErr.Status < 500 {
				http.Error(w, apiErr.Error(), apiErr.Status)
				return
			}
			http.Error(w, "authorization request failed", http.StatusBadGateway)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(struct {
			AuthorizeURL string `json:"authorize_url"`
			ReturnTo     string `json:"return_to"`
		}{link, destination})
	})
	addr := os.Getenv("LISTEN_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	log.Fatal(http.ListenAndServe(addr, nil))
}
