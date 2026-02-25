package server

import (
	"fmt"
	"io"
	"log"
	"main/internal/utils"
	"net/http"
	"time"
)

/*
This program is a simple web server that can be used to test the implementation of the server package.
This don't require any external packages.
*/

func Server() { // check the implementation usign "github.com/go-chi/chi/v5"
	mux := http.NewServeMux()

	// the {$} demands and exact match, be sure to add the / before the {$}

	mux.HandleFunc("GET /test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Test With get method endpoint papu")
	})

	mux.HandleFunc("POST /test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Test With post method endpoint papu")
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close() // good practice to close the body
		if err != nil {
			fmt.Fprintf(w, "Error reading body: %s", err)
			return
		}
		fmt.Fprintf(w, "Test With post for sub-path method endpoint papu\n %s", body)
	})

	mux.HandleFunc("POST /test/api", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close() // good practice to close the body
		if err != nil {
			fmt.Fprintf(w, "Error reading body: %s", err)
			return
		}
		fmt.Fprintf(w, "Test With post for sub-path method endpoint papu\n %s", body)
	})

	mux.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method { // Ugly way to handle methods
		case http.MethodGet:
			// Handle GET request
			fmt.Fprintf(w, "Getting users")
		case http.MethodPost:
			// Handle POST request
			fmt.Fprintf(w, "Post endpoint papu")
		default:
			// Method not allowed
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "404 - Page not found: %s", r.URL.Path)
	})

	port := ":" + utils.LoadPort()

	server := &http.Server{
		Addr:         port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	log.Println("Server started on ", port)
	log.Fatal(server.ListenAndServe())

}
