package main

import (
	"log"
	"net/http"
	"text/template"
)

type Data struct {
	Title string
}

func main() {
	homeHandler := func(w http.ResponseWriter, req *http.Request) {
		tmpl := template.Must(template.ParseFiles("index..gohtml"))

		data := Data{
			"Go setup",
		}

		tmpl.Execute(w, data)
	}

	http.HandleFunc("/", homeHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
