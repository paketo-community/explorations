package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type ExampleRouter struct {
	*mux.Router
	tmpl *template.Template
}

func NewExampleRouter() (*ExampleRouter, error) {
	r := mux.NewRouter()

	tmpl, err := template.ParseGlob("./web/templates/*.tmpl")
	if err != nil {
		return nil, err
	}

	router := &ExampleRouter{
		Router: r,
		tmpl:   tmpl,
	}

	fs := http.FileServer(http.Dir("./web"))
	r.HandleFunc("/", router.index)
	r.PathPrefix("/").Handler(fs)

	return router, nil
}

func (r *ExampleRouter) index(w http.ResponseWriter, req *http.Request) {
	err := r.tmpl.ExecuteTemplate(w, "index.tmpl", map[string]interface{}{
		"Greeting": "Hello",
		"Congrats": false, // Try switching this to true!
	})
	if err != nil {
		log.Printf("index: %v", err)
	}
}

func main() {
	router, err := NewExampleRouter()
	if err != nil {
		log.Fatalf("Router creation failed: %v", err)
	}
	http.Handle("/", router)

	log.Println("Serving on port 8000")
	log.Println("TEST CHANGE")
	err = http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatalf("Server exited with: %v", err)
	}
}
