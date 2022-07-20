package application

import (
	"fmt"
	"path/filepath"
	"time"
)

/*
コンテナで起動するアプリケーションの構造体
*/
type Application struct {
	UserName        string // required
	ApplicationName string // required
	FileName        string // required
}

/*
アプリケーションのパスからアプリケーションの構造体を生成する関数

引数
path - アプリケーションのファイルパス

返り値
app Application - アプリケーションの構造体
error
*/
func ParseApplicationFromPath(path string) (*Application, error) {
	applicationName := filepath.Base( filepath.Dir(path) )
	userName := filepath.Base( filepath.Dir(filepath.Dir(path)) )
	if applicationName == "" || userName == "" {
		return nil, fmt.Errorf("error: invalid path")
	}

	fileName := time.Now().UTC().Format(time.RFC3339Nano)

	return &Application{
		UserName:        userName,
		ApplicationName: applicationName,
		FileName:        fileName,
	}, nil
}

// TODO ネームフォーマットの読み込み
/*
コンテナ名を組み立てるメソッド
*/
func (p *Application) AssembleContainerName() string {
	return fmt.Sprintf("%s-%s", p.UserName, p.ApplicationName)
}

// TODO 設定ファイルから読み込み
/*
アプリケーションの待機時のパス
*/
func (p *Application) AssembleImcomingPath() string {
	return fmt.Sprintf("/queue/imcoming/%s/%s",  p.UserName, p.ApplicationName)
}

// TODO 設定ファイルから読み込み
/*
アプリケーションの保存時のパス
*/
func (p *Application) AssembleActivePath() string {
	return fmt.Sprintf("/queue/active/%s/%s",  p.UserName, p.ApplicationName)
}
