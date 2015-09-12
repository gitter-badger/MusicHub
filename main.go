package main

import "net/http"
import "server/server/serveraux"
import "fmt"

func main() {
	fmt.Println("Server Ready!!!!!!")
	http.HandleFunc("/search", serveraux.Search)
	http.HandleFunc("/upload", serveraux.Upload)
	http.HandleFunc("/hear",serveraux.Hear)
	http.HandleFunc("/credits",serveraux.Credits)
	http.HandleFunc("/404",serveraux.Fourofour)
	http.HandleFunc("/", serveraux.Handler)
	http.ListenAndServe(":1080", nil)
}
