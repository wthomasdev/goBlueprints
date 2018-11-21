package main

import (
	"flag"
	"log"
	"mycode/trace"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"
)

// temp1 represents a single template
type templateHandler struct {
	once     sync.Once
	filename string
	temp1    *template.Template
}

// ServeHTTP handles the HTTP request.
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.temp1 = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.temp1.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080", "the addr of the application.")
	flag.Parse()
	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	http.Handle("/", &templateHandler{filename: "chat.go.html"})
	http.Handle("/room", r)
	// get the room going
	go r.run()
	log.Println("Starting webserver on:", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
