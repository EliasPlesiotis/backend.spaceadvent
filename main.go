package main

import (
	"log"
	"net/http"

	"github.com/EliasPlesiotis/spaceadvent/website/handlers"

	"github.com/gorilla/mux"
)

func errorHandling(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				http.Error(w, "Server Error", 500)
				log.Print(r)
			}
		}()
		f(w, r)
	}
}

func logger(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Print(r.Method, r.URL.Path)
		f(w, r)
	}
}

func main() {
	var r = mux.NewRouter()

	r.HandleFunc("/", logger(errorHandling(handlers.Index)))
	r.HandleFunc("/login", logger(errorHandling(handlers.Login)))
	r.HandleFunc("/register", logger(errorHandling(handlers.Register)))
	r.HandleFunc("/logout", logger(errorHandling(handlers.Logout)))
	r.HandleFunc("/settings", logger(errorHandling(handlers.Settings)))
	r.HandleFunc("/replace", logger(errorHandling(handlers.Replace)))
	r.HandleFunc("/messages", logger(errorHandling(handlers.Messages)))
	r.HandleFunc("/sent", logger(errorHandling(handlers.Sent)))
	r.HandleFunc("/resetmsg", logger(errorHandling(handlers.ResetMsg)))
	r.HandleFunc("/delete", logger(errorHandling(handlers.DeleteMe)))

	http.Handle("/", r)
	http.ListenAndServe(":8080", r)

}
