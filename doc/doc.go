package doc


import "net/http"
import "fmt"
import "html/template"

func Api(w http.ResponseWriter, r *http.Request) {
	fmt.Println("API DOCS Requested")
	fmt.Println("method : " , r.Method)
	t, _ := template.ParseFiles("apidocs.html")
	t.Execute(w,nil)
}
