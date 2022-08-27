package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// 入力のバリデーション
	file, fileHeader, err := r.FormFile("userfile")
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	userID := r.FormValue("userid")
	appName := r.FormValue("appname")
	if len(userID) == 0 || len(appName) == 0 {
		log.Println("Error:FormValue is Empty")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 保存用ディレクトリの作成
	storagePath := fmt.Sprintf("/uploadfiles/%s/%s", userID, appName)
	err = os.MkdirAll(storagePath, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// アプリケーションファイルの保存
	fileName := fileHeader.Filename
	datetime := time.Now().Format(time.RFC3339Nano)
	dst, err := os.Create(fmt.Sprintf("%s/%s-%s", storagePath, datetime, fileName))
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

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
