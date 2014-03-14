package main

import (
	"code.google.com/p/go-uuid/uuid"
	"code.google.com/p/go.net/websocket"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.Seed(time.Now().UTC().UnixNano())
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
	s := &http.Server{
		Addr:    ":8080",
		Handler: NewHandler(),
	}
	log.Printf("starting server...")
	err := s.ListenAndServe()
	if err != nil {
		log.Printf("oops, exited with error: %v\n", err)
	}
}

type Handler struct {
	locker      sync.Locker
	currentPost string
	clients     map[string]chan string
}

func NewHandler() *Handler {
	return &Handler{
		locker:  new(sync.Mutex),
		clients: make(map[string]chan string),
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	log.Printf("%s %q\n", r.Method, r.RequestURI)

	serve := func(f string) {
		http.ServeFile(w, r, "web/"+f)
	}

	switch r.URL.Path {
	case "/":
		serve("display.html")
	case "/kill":
		if r.Method == "POST" {
			os.Exit(0)
		} else {
			http.Error(w, "need to POST", http.StatusMethodNotAllowed)
		}

	case "/display.js":
		serve("display.js")
	case "/ws.js":
		serve("reconnecting-websocket.js")
	case "/jquery.js":
		serve("jquery.js")
	case "/jquery-2.0.3.min.map":
		serve("jquery-2.0.3.min.map")
	case "/display.css":
		serve("display.css")
	case "/ws":
		websocket.Handler(h.handle).ServeHTTP(w, r)
	case "/post":
		if r.Method == "PUT" {
			h.post(r.URL.Query().Get("url"))
			fmt.Fprintln(w, "ok")
		} else {
			http.Error(w, "need to PUT url's", http.StatusMethodNotAllowed)
		}
	case "/proxy":
		h.ServeProxy(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h *Handler) ServeProxy(w http.ResponseWriter, r *http.Request) {
	u := r.URL.Query().Get("url")
	resp, err := http.Get(u)
	if err != nil {
		http.Error(w, "can't get it", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	io.Copy(w, resp.Body)
}

func (d *Handler) post(u string) {
	d.locker.Lock()
	defer d.locker.Unlock()
	d.currentPost = u
	for _, ch := range d.clients {
		ch <- u
	}
}

func (d *Handler) register(id string, ch chan string) {
	d.locker.Lock()
	defer d.locker.Unlock()
	d.clients[id] = ch
}

func (d *Handler) deregister(id string) {
	d.locker.Lock()
	defer d.locker.Unlock()
	close(d.clients[id])
	delete(d.clients, id)
}

type Client struct {
	Id      string
	Href    string
	Channel chan string
}

func (d *Handler) handle(ws *websocket.Conn) {

	id := uuid.New()
	ch := make(chan string, 10)

	d.register(id, ch)
	defer d.deregister(id)

	if len(d.currentPost) > 0 {
		err := websocket.JSON.Send(ws, d.currentPost)
		if err != nil {
			return
		}
	}

	for c := range ch {
		err := websocket.JSON.Send(ws, c)
		if err != nil {
			return
		}
	}
}
