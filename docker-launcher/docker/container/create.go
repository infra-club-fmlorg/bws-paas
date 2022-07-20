package container

import (
	"context"
	network_ "docker-launcher/docker/network"
	"docker-launcher/model/application"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

/*
コンテナを生成する関数

引数
cli - Dockerクライアント
app - アプリケーションの構造体

返り値
created - 作成したコンテナのID及び警告のポインタ
error
*/
func Create(cli *client.Client, app application.Application) (*container.ContainerCreateCreatedBody, error) {
	created, err := cli.ContainerCreate(
		context.Background(),
		// TODO イメージをビルドする
		&container.Config{
			Image:      "ubuntu",
			Entrypoint: []string{app.AssembleActivePath()},
		},
		&container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type: mount.TypeVolume,
					// TODO 設定ファイルから読み込み
					// TODO Dockerボリュームからホストのボリュームに変更
					Source: fmt.Sprintf("application-active"),
					// TODO 設定ファイルから読み込み
					Target: "/queue/active",
				},
			},
		},
		&network.NetworkingConfig{}, nil, app.AssembleContainerName())
	if err != nil {
		return nil, err
	}
	return &created, nil
}

/*
ネットワークに接続済みのコンテナを生成する関数

docker run --network に相当

引数
cli - Dockerクライアント
app - アプリケーションの構造体のポインタ
networkID - DockerネットワークのID

返り値
*/
func CreateConnectedNetwork(cli *client.Client, app application.Application, networkID string) (*container.ContainerCreateCreatedBody, error) {
	created, err := Create(cli, app)
	if err != nil {
		return nil, err
	}

	err = network_.ConnectContainer(cli, networkID, created.ID, &network.EndpointSettings{})
	if err != nil {
		return nil, err
	}
	return created, nil
}
