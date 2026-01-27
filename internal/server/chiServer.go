package server

import (
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

/*
this project is a simple web server that uses the chi router to handle HTTP requests.
to install you can use the command: go get -u github.com/go-chi/chi/v5
*/

func CreateChiServer() {

	r := chi.NewRouter()

	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello world from Chi server")
	})

	r.Put("/test/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		body, err := io.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			fmt.Fprint(w, "Error reading body: ", err)
			return
		}
		w.Header().Add("content-type", "application/json")
		res := `{ "message": "Updating user with id", "id": "` + id + `", "body": ` + string(body) + `}`
		fmt.Fprint(w, res)
	})

	r.Post("/test/api", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			fmt.Fprint(w, "Error reading body: ", err)
			return
		}
		w.Header().Add("content-type", "application/json")
		res := `{ "message": "Test post papu from subPath with body", "body": ` + string(body) + `}`
		fmt.Fprint(w, res)
	})

	r.Post("/test", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			fmt.Fprint(w, "Error reading body: ", err)
			return
		}
		w.Header().Add("content-type", "application/json")
		res := `{ "message": "Test post papu with body", "body": ` + string(body) + `}`
		fmt.Fprint(w, res)
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "404 - Page not found fuck you: ", r.URL.Path)
	})

	serverPort := ":8080"

	fmt.Println("Chi server started on ", serverPort)

	http.ListenAndServe(serverPort, r)

}
