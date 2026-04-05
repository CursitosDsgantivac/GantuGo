package server

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"main/internal/utils"
	"net/http"
	"os"

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

	r.Get("/test/html", func(w http.ResponseWriter, r *http.Request) {
		_, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		// create the html template
		templateString := `
		<!DOCTYPE html>
		<html>
			<head>
				<meta charset="UTF-8">
				<title>{{.Title}}</title>
			</head>
			<body>
				{{range .Items}}<div>{{ . }}</div>{{else}}<div><strong>no rows</strong></div>{{end}}
			</body>
		</html>`

		check := func(err error) {
			if err != nil {
				log.Fatal(err)
			}
		}
		// build the template
		t, err := template.New("webpage").Parse(templateString)
		check(err)

		data := struct {
			Title string
			Items []string
		}{
			Title: "My page",
			Items: []string{
				"My photos",
				"My blog",
			},
		}

		// return the templated the first argument is the writer and the second argument is the data to be rendered
		// the data is writed in the first argument
		err = t.Execute(w, data)
		check(err)

		noItems := struct {
			Title string
			Items []string
		}{
			Title: "My another page",
			Items: []string{},
		}

		err = t.Execute(os.Stdout, noItems)
		check(err)

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

	serverPort := ":" + utils.LoadPort()

	fmt.Println("Chi server started on ", serverPort)

	http.ListenAndServe(serverPort, r)

}
