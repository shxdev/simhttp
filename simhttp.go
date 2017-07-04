package main

import(
    "net/http"
    "fmt"
    "flag"
    "os"
    "log"
    "io/ioutil"
    "mime"
    "strings"
)

var rootpath string
var indexfile string
var port string

func FileExist(path string) bool {
    stat, err := os.Stat(path)
    if err != nil &&  os.IsNotExist(err) {
        return false
    }
    return !stat.IsDir()
}

func webHandle(w http.ResponseWriter, r *http.Request){
    path:=rootpath+r.URL.Path[1:]
    if(path[len(path)-1]=='/'){
        path+=indexfile
    }
    fmt.Printf("path1:%s\n",path)
    if(!FileExist(path)){
        fmt.Printf("    (Not Found)\n")
        http.NotFound(w,r)
    }else{
        fileext:=path[strings.LastIndex(path,"."):]
        // fmt.Printf("    (%s)\n",fileext)
        filecontent,err:=ioutil.ReadFile(path)
        contenttype:=mime.TypeByExtension(fileext)
        w.Header().Set("Content-Type", contenttype)
        fmt.Printf("    (%s)\n",contenttype)
        if err==nil{
            _,err:=w.Write(filecontent)
            if err!=nil{
                http.Error(w,"",http.StatusInternalServerError)
            }
        }else{
            log.Fatal(err)
            http.Error(w,"",http.StatusInternalServerError)
        }
    }
}

func main() {
    flag.StringVar(&rootpath,"d","./","")
    flag.StringVar(&port,"p","80","")
    flag.StringVar(&indexfile,"i","index.html","")
    flag.Parse()

    http.HandleFunc("/",webHandle)
    fmt.Printf("Server will start at \"%s\" using port %s\n",rootpath,port);
    err:=http.ListenAndServe(":"+port,nil)
    if(err!=nil){
        log.Fatal("ListenAndServe: ", err)
    }
}