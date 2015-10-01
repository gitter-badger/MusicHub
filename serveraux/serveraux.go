package serveraux

import 	"fmt"
import	"net/http"
import "os"
import "os/exec"
import "io"
import "html/template"
import "server/server/mp3metap"
import "server/server/returnMD5"
import _ "server/go-sqlite3"
import "database/sql"
import "strings"

func times(str string, ch byte)(int) {
	i := 0
	count :=0
	for(i < len(str)){
		if (str[i]==ch){
			count+=1
		}
		i++
	}
	return count
}

func log(printa string) {
	fmt.Println(printa);
}

//Function for API requests
//For More details refer api url of your server.
//(http://<address>:<port>/api)
func Api(w http.ResponseWriter , r *http.Request){
	fmt.Println("API_URL Requested")
	if r.Method=="GET"{
		urlQuery := string(r.URL.RawQuery)[:]
		fmt.Println(urlQuery)
		if (urlQuery == ""){
			fmt.Fprintf(w,"Bad Request")
			return
		}

		urlParts := strings.Split(urlQuery,"&")
		fmt.Println(urlParts)
		if(times(urlQuery,'=')<3){
			urlQuery += "="
		}
		var queryString string ="sample"
		var searchBase string ="album"
		var mode string ="xml"
		i := 0
		for (i<len(urlParts)){
				if(strings.Contains(urlParts[i],"query")){
					fmt.Println("Part %d contains query",i)
					if (!strings.Contains(urlParts[i],"=")){
						urlParts[i]+="="
					}
					queryString = strings.Split(urlParts[i],"=")[1]
				}
				if(strings.Contains(urlParts[i],"based")){
					fmt.Println("Part %d contains based",i)
					if (!strings.Contains(urlParts[i],"=")){
						urlParts[i]+="="
					}
					searchBase = strings.Split(urlParts[i],"=")[1]
				}
				if(strings.Contains(urlParts[i],"mode")){
					fmt.Println("Part %d contains mode",i)
					if (!strings.Contains(urlParts[i],"=")){
						urlParts[i]+="="
					}
					mode = strings.Split(urlParts[i],"=")[1]
				}
				i++;
				fmt.Println(i,len(urlParts))
		}

		if (len(searchBase)==0) {
			log("Empty")
			searchBase ="album"
		}
		if (len(queryString)==0) {
			log("Empty")
			queryString ="sample"
		}
		if(len(mode)==0){
			log("EMPTY")
			mode="xml"
		}
		fmt.Println(queryString,searchBase,mode)
		mdb,err := sql.Open("sqlite3","mdb.db")
		var apiresult string = "<songs>\n"
		if(mode == "json"){
			apiresult  = "{\"songs\": {\"song\":["
		}
		if err != nil {
				fmt.Println("Unable to access db")
		}
		rows, err := mdb.Query("select * from mdatabase")
		if err != nil {
			fmt.Println(err)
		}
		var title,album,artist,year,md5sum,songURL string
		var songNO int
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&songNO,&title,&album,&artist,&year,&md5sum,&songURL)
			songURL = strings.Replace(songURL, "./song/","",1)
			songURL = "hear?=" + songURL
			if err != nil {
				fmt.Println(err)
			}
			switch searchBase{
			case "title":
				if (isComparable(title,queryString)){
					if (mode == "json"){
						apiresult = apiresult + "{\"title\":\""+title+"\",\"album\":\""+album+"\",\"artist\":\""+artist+"\",\"year\":\""+year+"\",\"url\":\""+songURL+"\"},"
					}else{
						apiresult = apiresult + "\t<song title=\""+title+"\" album=\""+album +"\" artist=\""+artist +"\" year=\"" + year + "\" url=\""+songURL+"\" />\n"
					}
				}
			case "album":
				if (isComparable(album,queryString)){
					if (mode == "json"){
						apiresult = apiresult + "{\"title\":\""+title+"\",\"album\":\""+album+"\",\"artist\":\""+artist+"\",\"year\":\""+year+"\",\"url\":\""+songURL+"\"},"
					}else{
						apiresult = apiresult + "\t<song title=\""+title+"\" album=\""+album +"\" artist=\""+artist +"\" year=\"" + year + "\" url=\""+songURL+"\" />\n"
					}
				}
			case "artist":
				if (isComparable(artist,queryString)){
					if (mode == "json"){
						apiresult = apiresult + "{\"title\":\""+title+"\",\"album\":\""+album+"\",\"artist\":\""+artist+"\",\"year\":\""+year+"\",\"url\":\""+songURL+"\"},"
					}else{
						apiresult = apiresult + "\t<song title=\""+title+"\" album=\""+album +"\" artist=\""+artist +"\" year=\"" + year + "\" url=\""+songURL+"\" />\n"
					}
				}
			case "year":
				if (isComparable(year,queryString)){
					if (mode == "json"){
						apiresult = apiresult + "{\"title\":\""+title+"\",\"album\":\""+album+"\",\"artist\":\""+artist+"\",\"year\":\""+year+"\",\"url\":\""+songURL+"\"},"
					}else{
						apiresult = apiresult + "\t<song title=\""+title+"\" album=\""+album +"\" artist=\""+artist +"\" year=\"" + year + "\" url=\""+songURL+"\" />\n"
					}
				}
			default:
				if (isComparable(album,queryString)){
					if (mode == "json"){
						apiresult = apiresult + "{\"title\":\""+title+"\",\"album\":\""+album+"\",\"artist\":\""+artist+"\",\"year\":\""+year+"\",\"url\":\""+songURL+"\"},"
					}else{
						apiresult = apiresult + "\t<song title=\""+title+"\" album=\""+album +"\" artist=\""+artist +"\" year=\"" + year + "\" url=\""+songURL+"\" />\n"
					}
				}
			}
		}
		if(mode == "json"){
			apiresult = apiresult[:len(apiresult)-1]
			apiresult = apiresult + "]}}"
			w.Header().Set("Content-Type","application/json; charset=utf-8")
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}else{
			apiresult = apiresult + "</songs>"
			w.Header().Set("Content-Type","application/xml")
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		fmt.Fprintf(w,apiresult)
	}
}


//This function is used to compare different strings....
//As 'Ko' and 'KO' are not the same according to go (any other damn language), this function considers
// ('Ko' == 'KO') to be true
func isComparable(main string , compared string) (bool){
	value := false
	if compared == " " || compared == "  " || compared == "  " {
	}else{
		if(compared == "*ALL^*"){
			value = true
		}
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
	fmt.Println("URL Requested")
	//baseurl := r.URL.Path[1:]
	fmt.Println("method : " , r.Method)
	t, _ := template.ParseFiles("hubJS/audio.html")
	// fmt.Println(t)
	t.Execute(w,nil)
}


//The function handler for the uploader page
func Upload(w http.ResponseWriter , r *http.Request){
	// fmt.Println("method : " , r[0:13],"}")
	if r.Method == "GET" {
		t, _ := template.ParseFiles("upload.html")
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
		// if(handler.Header["Content-Type"][0] == "audio/mpeg" || handler.Header["Content-Type"][0] == "audio/mp3") {
		switch (handler.Header["Content-Type"][0]){
		case "audio/mpeg":
		case "audio/mp3":
			if(strings.Contains(handler.Filename," ")){
				handler.Filename = strings.Replace(handler.Filename," ","_",-1)
				fmt.Println(handler.Filename, "indn")
			}
			if (strings.Contains(handler.Filename,"'")) {
				handler.Filename = strings.Replace(handler.Filename,"'","_",-1)				
			}
			if (strings.Contains(handler.Filename,"(")) {
				handler.Filename = strings.Replace(handler.Filename,"(","_",-1)				
			}
			if (strings.Contains(handler.Filename,")")) {
				handler.Filename = strings.Replace(handler.Filename,")","_",-1)				
			}
			if (strings.Contains(handler.Filename,"%")) {
				handler.Filename = strings.Replace(handler.Filename,"%","_",-1)				
			}
			if (strings.Contains(handler.Filename,"[")) {
				handler.Filename = strings.Replace(handler.Filename,"[","_",-1)				
			}
			if (strings.Contains(handler.Filename,"]")) {
				handler.Filename = strings.Replace(handler.Filename,"]","_",-1)				
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
			cmd := exec.Command("alcov","/home/aki237/gospace/src/server/server/song/"+handler.Filename)
			fmt.Println(cmd)
			err = cmd.Run()
			if err != nil{
				fmt.Println("Unable to generate album art for "+"./song/"+handler.Filename+"\n",err)
			}
			fmt.Println(title,album,artist,year,md5sum)
			mdb,err := sql.Open("sqlite3","mdb.db")
			if err != nil {
					fmt.Println("Unable to access db")
			}
			rows, err := mdb.Query("select * from mdatabase where  MD5Sum= ?", md5sum)
			if err != nil{
				fmt.Println(err)
			}
			defer rows.Close()
			if (!rows.Next()){
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
				fmt.Fprintf(w,"<script>alert(\"Song successfully uploaded and added to the database\");window.location=\"/upload\";</script>")
			}else{
				fmt.Println("Already the song is in the database")
				w.Header().Set("Content-Type","text/html",)
				fmt.Fprintf(w,"<center><strong>The given song is already in the database<strong></center>")
			}
			defer mdb.Close()
		default:
			w.Header().Set("Content-Type","text/html",)
			fmt.Fprintf(w,"<script type=\"text/javascript\">alert(\"Please upload Valid mp3 files.....\")</script>")
		}
	}
}


//The http://<address>:<port>/ display funtion handler\
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
		if(strings.Contains(urlQuery,"../") || strings.Contains(urlQuery,"..")){
			fmt.Fprintf(w,"Oh please I dont know about that.... You sucker")
			return
		}
		urlQuery = "./song/"+ urlQuery
		if strings.Contains(urlQuery,"%20"){
			strings.Replace(urlQuery,"%20","\\ ",-1)
		}
		fmt.Println(urlQuery)
		http.ServeFile(w,r,urlQuery)
	}
}

//The http://<address>:<port>/ display funtion handler
func Module(w http.ResponseWriter, r *http.Request) {
	fmt.Println("MODULE_URL Requested")
	if r.Method=="GET"{
		urlQuery := string(r.URL.RawQuery)[1:]
		if(strings.Contains(urlQuery,"../") || strings.Contains(urlQuery,"..")){
			fmt.Fprintf(w,"Oh please I dont know about that.... You sucker")
			return
		}
		urlQuery = "./hubJS/"+ urlQuery
		if strings.Contains(urlQuery,"%20"){
			strings.Replace(urlQuery,"%20","\\ ",-1)
		}
		fmt.Println(urlQuery)
		http.ServeFile(w,r,urlQuery)
	}
}
