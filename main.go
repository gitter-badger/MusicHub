package main

import "net/http"
import "server/server/serveraux"
import "fmt"

func main() {
	fmt.Println("Server Ready!!!!!!\nListening at localhost:1080")
	http.HandleFunc("/api/search",serveraux.Api)
	http.HandleFunc("/hear",serveraux.Hear)
	http.HandleFunc("/search", serveraux.Search)
	http.HandleFunc("/upload", serveraux.Upload)
	http.HandleFunc("/credits",serveraux.Credits)
	http.HandleFunc("/404",serveraux.Fourofour)
	http.HandleFunc("/", serveraux.Handler)
	http.ListenAndServe("localhost:1080", nil)
}
