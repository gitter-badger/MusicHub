package mp3metap

import (
	"fmt"
	"os"
)
//This funtion is the only function in this package.
//This function is used to read the meta-data of any valid mp3 file.
//The concept behind it is that the last 128 bits(characters) in mpeg/mp3 files contain the information about the song.
//So it reads the 128 last bits and reads	the meta-data according to the ID3 standards.
//"import" this package and pass function as :
/*
		var info [4]string
		info = mp3metap.Metaparse(filename)
*/
//where filename is  a string
func Metaparse(filename string)(string,string,string,string,){
	//filename:=  os.Args[1]
	title := ""
	album := ""
	artist := ""
	year := ""
	file, error := os.Open(filename)
	if error != nil {
		fmt.Println("Unable to open file")
	}else{
		info, error := file.Stat()
		if error != nil {
			fmt.Println("Unable to find the statistics of the file")
		}else{
			bs := make([]byte, info.Size())
			_, error = file.Read(bs)
			if error != nil {
				fmt.Println("Unable to read from the file")
			}else{
				str := string(bs)
				str=str[info.Size()-128:]
				fmt.Println("Title :",str[3:33])
				title = str[3:33]
				fmt.Println("Album :",str[63:93])
				album = str[63:93]
				fmt.Println("Artist :",str[33:63])
				artist = str[33:63]
				fmt.Println("Year :",str[93:97])
				year = str[93:97]
			}
		}
	}
	return title , album , artist , year
}
