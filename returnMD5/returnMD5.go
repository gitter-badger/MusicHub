package returnMD5

import(
  "fmt"
  "crypto/md5"
  "os"
  "encoding/hex"
)

//Returns the md5 Hash for the given string
func ReturnMD5forString(str []byte) (string){
  bit := md5.Sum(str)
  var stri string = hex.EncodeToString(bit[:])
  return stri
}

//Returns the md5 Hash for a file
func ReturnMD5(filename string)(string){
  f,err := os.Open(filename)
  if err != nil {
    fmt.Println("Error")
    var stro string = ""
    return stro
  }
  defer f.Close()
  info,_ := f.Stat()
  content := make([]byte,info.Size())
  _,err = f.Read(content)
  bit := md5.Sum(content)
  var str string = hex.EncodeToString(bit[:])
  return str
}
