package decorator

import (
	"fmt"
	"log"
	"net/http"
)

func autoAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("Auth")
		if err != nil || cookie.Value != "Authentic" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		h(w, r)
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Println(w, "Hello world!"+r.URL.Path)
}

func main() {
	http.HandleFunc("/hello", autoAuth(hello))
	err := http.ListenAndServe(":5666", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
