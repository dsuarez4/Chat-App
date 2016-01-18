package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServeHTTP handles the HTTP request
// sync.once guarantees that the function passed as an argument will only
// be executed once, regardless of how many goroutines called ServeHTTP
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ =
			template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, nil)
}

func main() {
	fmt.Println("starting server on http://localhost:8080/")
	//http.HandleFunc("/", IndexHandler)
	//why do we pass the reference?
	http.Handle("/", &templateHandler{filename: "chat.html"})
	//start the Webserver
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

} //Can I add () to function call?

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`
		<html>
			<head>
				<title>Chat</title>
			</head>
			<body>
				Let's chat!
			</body>
		</html>

		`))
}
