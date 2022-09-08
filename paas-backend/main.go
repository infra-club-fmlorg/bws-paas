package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	BINARY               = "binary"
	NODE_JS              = "nodejs"
	HTML                 = "html"
	DOCKER_LAUNCHER_SOCK = "/socket/docker_launcher.sock"
	IMCOMINF_QUEUE_DIR   = "/queue/incoming"
	DATETIME_FORMAT      = time.RFC3339Nano
	FILE_NAME_SEPARATER  = "_"
)

/*
コンテナで起動するアプリケーションの構造体
*/
type ApplicationInfo struct {
	UserName        string    `json:"user_name"`        // required
	ApplicationName string    `json:"application_name"` // required
	Runtime         string    `json:"runtime"`          // required
	CreatedAt       time.Time `json:"created_at"`
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
	return fmt.Sprintf("%s/%s/%s", IMCOMINF_QUEUE_DIR, p.UserName, p.ApplicationName)
}

/*
ファイル名を組み立てるメソッド
*/
func (p *ApplicationInfo) AssembleFileName() string {
	return strings.Join([]string{p.Runtime, p.CreatedAt.Format(DATETIME_FORMAT)}, FILE_NAME_SEPARATER)
}

type Echo struct {
	Length int
	Data   []byte
}

func (e *Echo) String() string {
	return fmt.Sprintf("Length[%02d] Data[%s]", e.Length, e.Data)
}

func (e *Echo) Write(c net.Conn) error {
	data := make([]byte, 0, 4+e.Length)

	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(e.Length))
	data = append(data, buf...)

	w := bytes.Buffer{}
	err := binary.Write(&w, binary.BigEndian, e.Data)
	if err != nil {
		return err
	}

	data = append(data, w.Bytes()...)

	_, err = c.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (e *Echo) Read(c net.Conn) error {
	buf := make([]byte, 4)

	_, err := c.Read(buf)
	if err != nil {
		return err
	}

	byteCount := binary.BigEndian.Uint32(buf)
	e.Length = int(byteCount)
	e.Data = make([]byte, e.Length)

	_, err = c.Read(e.Data)
	if err != nil {
		return err
	}

	return nil
}

func NewEcho(buf []byte) *Echo {
	return &Echo{
		Length: len(buf),
		Data:   buf,
	}
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
	case BINARY, NODE_JS, HTML:
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
	conn, err := net.Dial("unix", DOCKER_LAUNCHER_SOCK)
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
	m := NewEcho(json)

	err = m.Write(conn)
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
