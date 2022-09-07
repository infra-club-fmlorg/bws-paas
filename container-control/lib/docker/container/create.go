package container

import (
	"container-controller/lib/application"
	network_ "container-controller/lib/docker/network"
	"context"

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
func Create(cli *client.Client, app application.ApplicationInfo) (*container.ContainerCreateCreatedBody, error) {
	created, err := cli.ContainerCreate(
		context.Background(),
		// TODO イメージをビルドする
		&container.Config{
			Image:      "ubuntu",
			Entrypoint: []string{app.AssembleActiveAppPath()},
		},
		&container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type: mount.TypeBind,
					// TODO 設定ファイルから読み込み
					// TODO Dockerボリュームからホストのボリュームに変更
					Source: "/tmp/bws-paas-queue/active",
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

func CreateFromImage(cli *client.Client, app *application.ApplicationInfo) (*container.ContainerCreateCreatedBody, error) {
	imageName := app.AssembleContainerName()
	created, err := cli.ContainerCreate(
		context.Background(),
		// TODO イメージをビルドする
		&container.Config{
			Image: imageName,
		},
		&container.HostConfig{
			AutoRemove: true,
		},
		&network.NetworkingConfig{}, nil, imageName)
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
func CreateConnectedNetwork(cli *client.Client, app application.ApplicationInfo, networkID string) (*container.ContainerCreateCreatedBody, error) {
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
