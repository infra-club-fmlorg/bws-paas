package application

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

/*
コンテナで起動するアプリケーションの構造体
*/
type ApplicationInfo struct {
	UserName        string    `json:"user_name"`        // required
	ApplicationName string    `json:"application_name"` // required
	Runtime         Runtime   `json:"runtime"`          // required
	CreatedAt       time.Time `json:"created_at"`
}

/*
ApplicationInfo構造体用のコンストラクタ関数
*/
func NewApplicationInfo(userName string, applicationName string, runtime Runtime) *ApplicationInfo {
	return &ApplicationInfo{
		UserName:        userName,
		ApplicationName: applicationName,
		Runtime:         runtime,
		CreatedAt:       time.Now(),
	}
}

/*
アプリケーションのパスからアプリケーションの構造体を生成する関数
{userName}/{appName}/{runtime}-{createdAt}

引数
path - アプリケーションのファイルパス

返り値
app ApplicationInfo - アプリケーションの構造体
error
*/
func ParseApplicationInfoFromPath(path string) (*ApplicationInfo, error) {
	fileName := filepath.Base(path)
	applicationName := filepath.Base(filepath.Dir(path))
	userName := filepath.Base(filepath.Dir(filepath.Dir(path)))
	if applicationName == "" || userName == "" {
		return nil, fmt.Errorf("error: invalid path")
	}

	parsedFileName := strings.Split(fileName, "_")
	createdAt, err := time.Parse(time.RFC3339Nano, parsedFileName[1])
	if err != nil {
		return nil, fmt.Errorf("error:fail parse to RFC3339Nano:%s", parsedFileName[1])
	}

	return &ApplicationInfo{
		UserName:        userName,
		ApplicationName: applicationName,
		Runtime:         Runtime(parsedFileName[0]),
		CreatedAt:       createdAt,
	}, nil
}

// TODO ネームフォーマットの読み込み
/*
コンテナ名を組み立てるメソッド
*/
func (p *ApplicationInfo) AssembleContainerName() string {
	return fmt.Sprintf("%s-%s", p.UserName, p.ApplicationName)
}

/*
ファイル名を組み立てるメソッド
*/
func (p *ApplicationInfo) AssembleFileName() string {
	return fmt.Sprintf("%s_%s", p.Runtime, p.CreatedAt.Format(time.RFC3339Nano))
}

// TODO 設定ファイルから読み込み
/*
アプリケーションの待機時のパス
*/
func (p *ApplicationInfo) AssembleIncomingDirPath() string {
	return fmt.Sprintf("/queue/incoming/%s/%s", p.UserName, p.ApplicationName)
}

// TODO 設定ファイルから読み込み
/*
アプリケーションの保存時のパス
*/
func (p *ApplicationInfo) AssembleActiveAppPath() string {
	return fmt.Sprintf("/queue/active/%s/%s/%s", p.UserName, p.ApplicationName, p.AssembleFileName())
}
