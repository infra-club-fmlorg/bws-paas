package main

import (
    "time"
    "path/filepath"
    "os"
    "io"
    "fmt"
    "log"
    "net/http"
    "strconv"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("test")
    //省略
    // ------------------------------
    //FormFileの引数はHTML内のform要素のnameと一致している必要があります
    file, fileHeader, err := r.FormFile("userfile")
    if err != nil {
        fmt.Println(err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    defer file.Close()

    // 存在していなければ、保存用のディレクトリを作成します。
    err = os.MkdirAll(fmt.Sprintf("/uploadfiles/%s/%s", r.FormValue("userid"), r.FormValue("appname")), os.ModePerm)
    if err != nil {
        fmt.Println(err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // 保存用ディレクトリ内に新しいファイルを作成します。
    t := time.Now();
    tuki := strconv.Itoa(int(t.Month()))
    fmt.Println(tuki)
    if len(tuki) == 1{
        tuki = "0" + tuki;
    }
    niti := strconv.Itoa(t.Day());
    if len(niti) == 1{
        niti = "0" + niti;
    }
    zi := strconv.Itoa(t.Hour());
    if len(zi) == 1{
        zi = "0" + zi;
    }
    hunn := strconv.Itoa(t.Minute());
    if len(hunn) == 1{
        hunn = "0" + hunn;
    }
    byou := strconv.Itoa(t.Second());
    if len(byou) == 1{
        byou = "0" + byou;
    }
    dst, err := os.Create(fmt.Sprintf("/uploadfiles/%s/%s/", r.FormValue("userid"), r.FormValue("appname")) + fmt.Sprintf("%d%s%s%s%s%s%s", t.Year(), tuki, niti, zi, hunn, byou, filepath.Ext(fileHeader.Filename)))
    if err != nil {
        fmt.Println(err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    defer dst.Close()

    // アップロードされたファイルを先程作ったファイルにコピーします。
    _, err = io.Copy(dst, file)
    if err != nil {
        fmt.Println(err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    fmt.Fprintf(w, "Upload successful")

}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, World")
}

func handler_(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "upload")
}


func setupRoutes() {
    mux := http.NewServeMux()
    mux.HandleFunc("/upload", uploadHandler) // 追加
    // mux.HandleFunc("/upload", handler_)
    mux.HandleFunc("/", handler)

    if err := http.ListenAndServe(":8080", mux); err != nil {
        log.Fatal(err)
    }
}

func main() {
    fmt.Println("FILE UPLOAD MONITOR")

    setupRoutes()
}