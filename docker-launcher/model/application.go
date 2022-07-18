package model

import (
	"fmt"
	"path/filepath"
	"time"
)

// 直接生成禁止
type Application struct {
	UserName        string
	ApplicationName string
	FileName        string
}

func ParseApplicationFromPath(path string) (Application, error) {
	applicationName := filepath.Dir(path)
	userName := filepath.Dir(applicationName)
	if applicationName == "" || userName == "" {
		return Application{
			UserName:        userName,
			ApplicationName: applicationName,
			FileName:        "",
		}, fmt.Errorf("error: invalid path")
	}

	fileName := time.Now().UTC().Format(time.RFC3339Nano)

	return Application{
		UserName:        userName,
		ApplicationName: applicationName,
		FileName:        fileName,
	}, nil
}

// TODO ネームフォーマットの読み込み
func (p *Application) AssembleContainerName() string {
	return fmt.Sprintf("%s-%s", p.UserName, p.ApplicationName)
}

// TODO 設定ファイルから読み込み
func (p *Application) AssembleImcomingPath() string {
	return fmt.Sprintf("/%s/%s/%s", "/queue/imcoiming", p.UserName, p.ApplicationName)
}

// TODO 設定ファイルから読み込み
func (p *Application) AssembleActivePath() string {
	return fmt.Sprintf("/%s/%s/%s", "/queue/imcoiming", p.UserName, p.ApplicationName)
}
