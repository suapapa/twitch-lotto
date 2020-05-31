package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var (
	port string
	from string
)

func init() {
	flag.StringVar(&port, "-p", "8000", "listen port")
	flag.StringVar(&from, "-f", "meetup", "collect member from")
}

func main() {
	flag.Parse()

	http.HandleFunc("/", IndexView)
	http.HandleFunc("/event", GetEventHandler)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Printf("Listening on %s.\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// IndexView render the index template
func IndexView(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html; charset=utf-8")
	f, err := os.Open("index.html")
	chk(err)
	io.Copy(w, f)
}

// GetEventHandler return the latest event and attendance info as json
func GetEventHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	writeErrAsJSON := func(w io.Writer, err error) {
		chk(
			json.NewEncoder(w).Encode(
				map[string]string{
					"message": err.Error(),
				},
			),
		)
	}

	users, err := MeetupResvUsersOfLastEvent()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeErrAsJSON(w, err)
		return
	}
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeErrAsJSON(w, err)
		return
	}
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
