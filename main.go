package main

import "net/http"
import "server/server/serveraux"
import "server/server/doc"
import "fmt"
import "io/ioutil"
import "encoding/json"


type Config struct{
	Listenaddress string
}

func main() {
	fmt.Println("Server Ready!!!!!!\nListening at localhost:1080")
	http.HandleFunc("/api",doc.Api)
	http.HandleFunc("/api/search",serveraux.Api)
	http.HandleFunc("/hear",serveraux.Hear)
	http.HandleFunc("/search", serveraux.Search)
	http.HandleFunc("/upload", serveraux.Upload)
	http.HandleFunc("/credits",serveraux.Credits)
	http.HandleFunc("/404",serveraux.Fourofour)
	http.HandleFunc("/", serveraux.Handler)
	content, err := ioutil.ReadFile("hubconfig.json")
	if(err != nil){
		fmt.Println(err)
	}
	var conf Config
	err = json.Unmarshal(content , &conf)
	if(err != nil){
		fmt.Println(err)
	}
	fmt.Println(conf.Listenaddress)
	http.ListenAndServe(conf.Listenaddress, nil)
}
