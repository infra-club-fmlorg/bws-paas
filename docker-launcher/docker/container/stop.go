package container

import (
	"context"
	"time"

	"github.com/docker/docker/client"
)

/*
Dockerコンテナを停止する関数

引数
cli - Dockerクライアント
containerID - DockerコンテナのID

返り値
error
*/
func Stop(cli *client.Client, containerID string) error {
  // TODO 設定ファイルから読み込み
	timeout := time.Duration(500) * time.Millisecond
	return cli.ContainerStop(
		context.Background(),
		containerID,
		&timeout,
	)
}
