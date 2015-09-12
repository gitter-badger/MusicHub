package serveraux

import 	"fmt"
import	"net/http"
import "os"
import "io"
import "html/template"
import "server/server/mp3metap"
import "server/server/returnMD5"
import _ "server/go-sqlite3"
import "database/sql"
import "strings"

//This function is used to compare different strings....
//As 'Ko' and 'KO' are not the same according to go (any other damn language), this function considers
// ('Ko' == 'KO') to be true
func isComparable(main string , compared string) (bool){
	value := false
	if compared == " " || compared == "  " || compared == "  " {
	}else{
		if(strings.Contains(main,compared) || strings.Contains(main,strings.ToUpper(compared))){
			value = true
		}
		if (strings.Contains(main,strings.ToLower(compared))){
			value = true
		}
		if(strings.Contains(main,strings.ToUpper(string(compared[0])) + compared[1:])){
			value = true
		}
	}
	return value
}



//Easter Egg !!!!!
/*
NO OFFENCE PLEASE
*/
func Fourofour(w http.ResponseWriter , r *http.Request)  {
	w.Header().Set("Content-Type","text/html")
	fmt.Fprintf(w,"<center><h1><strike><pre> 404 </pre></strike></h1></center><br>Wow you have found this page....")
}
//The funtion handler for the /credits URL
func Credits(w http.ResponseWriter , r *http.Request){
	fmt.Println("URL Requested")
	//baseurl := r.URL.Path[1:]
	fmt.Println("method : " , r.Method)
	t, _ := template.ParseFiles("credits.html")
	// fmt.Println(t)
	t.Execute(w,nil)
}

//The function handler for the search page
func Search(w http.ResponseWriter, r *http.Request){
	if r.Method == "GET" {
		fmt.Println("Search URL Requested")
		w.Header().Set("Content-Type","text/html",)
		fmt.Fprintf(w, "<body><center><font size=40>Write the search query : <br><form action=\"/search\" name=\"searchbox\" method=\"POST\"><input type=\"text\" name=\"keyword\" value=\"SAMPLEQUERY\"/><br><input type=\"submit\" value=\"Search\"/></form></font></center></body> ")
	}else{
		mdb,err := sql.Open("sqlite3","song/mdb.db")
		if err != nil {
				fmt.Println("Unable to access db")
		}
		rows, err := mdb.Query("select * from mdatabase")
		if err != nil {
			fmt.Println(err)
		}
		htmlContent := "<style>body{font-family:monospace;background-color:#FF265A;color:#000022;}td{border:solid;}a,a:hover,a:active,a:visited{text-decoration:none;}</style>"
		titleContent := "<body><h2>Results based on song title search</h2><br><br><table>"+"<tr><td><strong>Title</strong></td><td><strong>Album</strong></td><td><strong>Artist</strong></td><td><strong>Year</strong></td><td><strong>Download URL</strong></td></tr>"
		albumContent := "<h2>Results based on song album search</h2><br><br><table>"+"<tr><td><strong>Title</strong></td><td><strong>Album</strong></td><td><strong>Artist</strong></td><td><strong>Year</strong></td><td><strong>Download URL</strong></td></tr>"
		artistContent := "<h2>Results based on song artist search</h2><br><br><table>"+"<tr><td><strong>Title</strong></td><td><strong>Album</strong></td><td><strong>Artist</strong></td><td><strong>Year</strong></td><td><strong>Download URL</strong></td></tr>"
		yearContent := "<h2>Results based on song release year</h2><br><br><table>"+"<tr><td><strong>Title</strong></td><td><strong>Album</strong></td><td><strong>Artist</strong></td><td><strong>Year</strong></td><td><strong>Download URL</strong></td></tr>"
		r.ParseForm()
		keyword := r.FormValue("keyword")
		var title,album,artist,year,md5sum,songURL string
		var songNO int
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&songNO,&title,&album,&artist,&year,&md5sum,&songURL)
			songURL = strings.Replace(songURL, "./song/","",1)
			if err != nil {
				fmt.Println(err)
			}
			if (isComparable(title,keyword)){
				titleContent = titleContent + "<tr><td>"+title+"</td><td>"+album+"</td><td> "+artist+"</td><td> "+year+"</td><td> <a href=\"hear?="+songURL+"\">Download</a></td></tr>"
			}
			if (isComparable(album,keyword)){
				albumContent = albumContent + "<tr><td>"+title+"</td><td>"+album+"</td><td>"+artist+"</td><td>"+year+"</td><td><a href=\"hear?="+songURL+"\">Download</a></td></tr>"
			}
			if (isComparable(artist,keyword)){
				artistContent = artistContent + "<td>"+title+"</td><td>"+album+"</td><td>"+artist+"</td><td>"+year+"</td><td><a href=\"hear?="+songURL+"\">Download</a></td></tr>"
			}
			if (isComparable(year,keyword)){
				yearContent = yearContent + "<td>"+title+"</td><td>"+album+"</td><td>"+artist+"</td><td>"+year+"</td><td><a href=\"hear?="+songURL+"\">Download</a></td></tr>"
			}
			//fmt.Println(album, songURL)
		}
		titleContent = titleContent + "</table>"
		albumContent = albumContent + "</table>"
		artistContent = artistContent + "</table>"
		yearContent = yearContent + "</table></body>"
		htmlContent = htmlContent +  titleContent+albumContent+artistContent+yearContent
		// fmt.Println(htmlContent)
		w.Header().Set("Content-Type","text/html",)
		fmt.Println("Post Method")
		fmt.Fprintf(w, "<center><h1>Results</h1></center> \n %s",htmlContent)
	}
}


//The function handler for the uploader page
func Upload(w http.ResponseWriter , r *http.Request){
	// fmt.Println("method : " , r[0:13],"}")
	if r.Method == "GET" {
		t, _ := template.ParseFiles("upload.gtpl")
		t.Execute(w,nil)
	}else{
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("upload_file")
		if err !=nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Println(handler.Header["Content-Type"])
		//fmt.Println(returnMD5.ReturnMD5(handler.Filename))
		if(handler.Header["Content-Type"][0] == "audio/mpeg" || handler.Header["Content-Type"][0] == "audio/mp3"){
			if(strings.Contains(handler.Filename," ")){
				handler.Filename = strings.Replace(handler.Filename," ","_",-1)
				fmt.Println(handler.Filename, "indn")
			}
			f, err := os.OpenFile("./song/"+handler.Filename,os.O_WRONLY|os.O_CREATE,0666)
			if err !=nil {
				fmt.Println(err)
				return
			}
			defer f.Close()
			io.Copy(f, file)
			title,album,artist,year := mp3metap.Metaparse("./song/"+handler.Filename)
			fmt.Println("Songfile : ","./song/"+handler.Filename)
			md5sum := returnMD5.ReturnMD5("./song/"+handler.Filename)
			fmt.Println(title,album,artist,year,md5sum)
			mdb,err := sql.Open("sqlite3","song/mdb.db")
			if err != nil {
					fmt.Println("Unable to access db")
			}
			rows, err := mdb.Query("select * from mdatabase where  MD5Sum= ?", md5sum)
			if err != nil{
				fmt.Println(err)
			}
			defer rows.Close()
			if !rows.Next(){
				stmt, err := mdb.Prepare("INSERT INTO mdatabase(Title,Album,Artist,Year,MD5Sum,SongURL) VALUES(?,?,?,?,?,?)")
				if err != nil{
					fmt.Println("Unable to prepare database")
				}
				_,err = stmt.Exec(title,album,artist,year,md5sum,"./song/"+handler.Filename)
				if err != nil {
					fmt.Println("Unable to process meta-data\n",err)
					w.Header().Set("Content-Type","text/html",)
					fmt.Fprintf(w,"<center><strong>Unable to process the metadata : May be the song name clashes with another song<br>Try saving it with a different name<strong></center>")
				}
			}else{
				fmt.Println("Already the song is in the database")
				w.Header().Set("Content-Type","text/html",)
				fmt.Fprintf(w,"<center><strong>The given song is already in the database<strong></center>")
			}
			defer mdb.Close()
		}else{
			w.Header().Set("Content-Type","text/html",)
			fmt.Fprintf(w,"<script type=\"text/javascript\">alert(\"Please upload Valid mp3 files.....\")</script>")
		}
	}
}


//The http://<address>:<port>/ display funtion handler
func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("URL Requested")
	//baseurl := r.URL.Path[1:]
	fmt.Println("method : " , r.Method)
	t, _ := template.ParseFiles("home.html")
	// fmt.Println(t)
	t.Execute(w,nil)
}

//The http://<address>:<port>/ display funtion handler
func Hear(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HEAR_URL Requested")
	if r.Method=="GET"{
		urlQuery := string(r.URL.RawQuery)[1:]
		urlQuery = "./song/"+ urlQuery
		if strings.Contains(urlQuery,"%20"){
			strings.Replace(urlQuery,"%20","\\ ",-1)
		}
		fmt.Println(urlQuery)
		http.ServeFile(w,r,urlQuery)
	}
}
