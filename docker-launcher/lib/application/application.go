package application

import (
	"fmt"
	"path/filepath"
	"time"
)

/*
コンテナで起動するアプリケーションの構造体
*/
type ApplicationInfo struct {
	UserName        string `json:"user_name"`        // required
	ApplicationName string `json:"application_name"` // required
	FileName        string `json:"file_name"`        // required
	Runtime         string `json:"runtime"`          // required
}

/*
ApplicationInfo構造体用のコンストラクタ関数
*/
func NewApplicationInfo(userName string, applicationName string, fileName string, runtime string) *ApplicationInfo {
	return &ApplicationInfo{
		UserName:        userName,
		ApplicationName: applicationName,
		FileName:        fileName,
		Runtime:         runtime,
	}
}

/*
アプリケーションのパスからアプリケーションの構造体を生成する関数

引数
path - アプリケーションのファイルパス

返り値
app ApplicationInfo - アプリケーションの構造体
error
*/
func ParseApplicationInfoFromPath(path string) (*ApplicationInfo, error) {
	applicationName := filepath.Base(filepath.Dir(path))
	userName := filepath.Base(filepath.Dir(filepath.Dir(path)))
	if applicationName == "" || userName == "" {
		return nil, fmt.Errorf("error: invalid path")
	}

	fileName := time.Now().UTC().Format(time.RFC3339Nano)

	return &ApplicationInfo{
		UserName:        userName,
		ApplicationName: applicationName,
		FileName:        fileName,
	}, nil
}

// TODO ネームフォーマットの読み込み
/*
コンテナ名を組み立てるメソッド
*/
func (p *ApplicationInfo) AssembleContainerName() string {
	return fmt.Sprintf("%s-%s", p.UserName, p.ApplicationName)
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
func (p *ApplicationInfo) AssembleActivePath() string {
	return fmt.Sprintf("/queue/active/%s/%s/%s", p.UserName, p.ApplicationName, p.FileName)
}
