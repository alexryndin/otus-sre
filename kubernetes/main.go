package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

const version = "1.0.0"

type BasicAuth struct {
	User     string
	Password string
}

func NewBasicAuth() (*BasicAuth, error) {
	user := os.Getenv("APP_USER")
	pass := os.Getenv("APP_PASSWORD")

	if user == "" || pass == "" {
		return nil, fmt.Errorf("APP_USER and APP_PASSWORD required")
	}

	return &BasicAuth{user, pass}, nil
}

func (b *BasicAuth) check(user string, pass string) bool {
	return b.User == user && b.Password == pass
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Println("PORT is required")
		os.Exit(1)
	}

	bauth, err := NewBasicAuth()
	if err != nil {
		log.Printf("%v\n", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.HandleFunc(`/`, app)
	mux.HandleFunc("/admin", auth(admin, bauth))
	mux.HandleFunc("/healthz", healthz)

	log.Printf("Server listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func app(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving request: %s", r.URL.Path)
	host, _ := os.Hostname()

	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.Encode(map[string]string{"app": "Go App", "zone": "public", "version": version, "hostname": host})
}

func admin(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving request: %s", r.URL.Path)

	host, _ := os.Hostname()

	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.Encode(map[string]string{"app": "Go app", "zone": "private", "version": version, "hostname": host})
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func auth(fn http.HandlerFunc, auth *BasicAuth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || !auth.check(user, pass) {
			http.Error(w, "Unauthorized.", 401)
			return
		}
		fn(w, r)
	}
}
