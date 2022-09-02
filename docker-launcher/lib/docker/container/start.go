package container

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

/*
コンテナを開始する関数

引数
cli - Dockerクライアント
containerID - DockerコンテナのID

返り値
error
*/
func Start(cli *client.Client, containerID string) error {
	if err := cli.ContainerStart(context.Background(), containerID, types.ContainerStartOptions{}); err != nil {
		return err
	}
	return nil
}
