package network

import (
	"context"

	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

/*
コンテナをネットワークに接続する関数

引数
cli - Dockerクライアント
networkID - DockerネットワークのID
containerID - Docker ContainerのID
option - 接続時のオプション

返り値
error
*/
func ConnectContainer(cli *client.Client, networkID string, containerID string, option *network.EndpointSettings) error {
	return cli.NetworkConnect(context.Background(), networkID, containerID, option)
}
