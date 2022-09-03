package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

const (
	BINARY  = "binary"
	NODE_JS = "nodejs"
	PYTHON  = "python"
)

/*
コンテナで起動するアプリケーションの構造体
*/
type ApplicationInfo struct {
	UserName        string    `json:"user_name"`        // required
	ApplicationName string    `json:"application_name"` // required
	Runtime         string    `json:"runtime"`          // required
	CreatedAt       time.Time `json:"-"`
}

/*
ApplicationInfo構造体用のコンストラクタ関数
*/
func NewApplicationInfo(userName string, applicationName string, runtime string) *ApplicationInfo {

	return &ApplicationInfo{
		UserName:        userName,
		ApplicationName: applicationName,
		Runtime:         runtime,
		CreatedAt:       time.Now(),
	}
}

// TODO 設定ファイルから読み込み
/*
アプリケーションの待機時のパス
*/
func (p *ApplicationInfo) AssembleIncomingDirPath() string {
	return fmt.Sprintf("/queue/incoming/%s/%s", p.UserName, p.ApplicationName)
}

/*
ファイル名を組み立てるメソッド
*/
func (p *ApplicationInfo) AssembleFileName() string {
	return fmt.Sprintf("%s-%s-%s-%s", p.UserName, p.ApplicationName, p.Runtime, p.CreatedAt)
}

func uploadHandler(res http.ResponseWriter, req *http.Request) {
	// 入力のバリデーションと構造体化
	//// 入力の構造体化
	applicationInfo := NewApplicationInfo(
		req.FormValue("user_name"),
		req.FormValue("application_name"),
		req.FormValue("runtime"),
	)

	log.Printf("receive request:%#v", applicationInfo)

	//// 入力のバリデーション
	if len(applicationInfo.ApplicationName) == 0 {
		err := fmt.Errorf("Error:application_name is empty")
		log.Println(err)
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if len(applicationInfo.UserName) == 0 {
		err := fmt.Errorf("Error:application_UserName is empty")
		log.Println(err)
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	switch applicationInfo.Runtime {
	case BINARY, NODE_JS, PYTHON:
	default:
		err := fmt.Errorf("Error:environment is not supported")
		log.Println(err)
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	// 送信されたアプリケーションファイルの検証
	applicationFile, _, err := req.FormFile("application_file")
	if err != nil {
		fmt.Println(err)
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	defer applicationFile.Close()

	log.Printf("success validation:%#v", applicationInfo)

	imcomingDirPath := applicationInfo.AssembleIncomingDirPath()
	fileName := applicationInfo.AssembleFileName()
	filePath := fmt.Sprintf("%s/%s", imcomingDirPath, fileName)

	// 保存用ディレクトリの作成
	err = os.MkdirAll(imcomingDirPath, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	// アプリケーションファイルの保存
	f, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	_, err = io.Copy(f, applicationFile)
	if err != nil {
		fmt.Println(err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("success save file:%s", filePath)

	// Unix Domain Socket経由でContainer起動サブシステムへ通知
	conn, err := net.Dial("unix", "/var/run/docker_launcher.sock")
	if err != nil {
		log.Println(err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
	defer conn.Close()

	json, err := json.Marshal(applicationInfo)
	if err != nil {
		log.Println(err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

	_, err = conn.Write(json)
	if err != nil {
		log.Println(err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

	log.Printf("success notification to container launcher system:%s", filePath)
}

func setupRoutes() {
	mux := http.NewServeMux()
	mux.HandleFunc("/upload", uploadHandler) // 追加

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

func main() {
	fmt.Println("file upload monitor")

	setupRoutes()
}
