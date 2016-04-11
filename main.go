package main

//page 34

import (
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
	"flag"
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
	t.templ.Execute(w, r)
}

func main() {

	// returns *string
	// in Go, returning *type return the address of the value.
	// To get the value itself must use *addr
	var addr = flag.String("addr", "8080", "The addr of the application")
	flag.Parse();

	r := newRoom()

	//why do we pass the reference?
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)

	// run the room
	go r.run()

	log.Println("starting server on,", *addr);


	//start the Webserver
	if err := http.ListenAndServe(*addr, nil); err != nil {
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
