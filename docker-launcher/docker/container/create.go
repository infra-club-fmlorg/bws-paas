package container

import (
	"context"
	"docker-launcher/model"
	network_ "docker-launcher/docker/network"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

func CreateContainer(cli *client.Client, app model.Application) (*container.ContainerCreateCreatedBody, error) {
	created, err := cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image:      "ubuntu",
			Entrypoint: []string{app.AssembleActivePath()},
		},
		&container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type: mount.TypeVolume,
					// TODO 設定ファイルから読み込み
					Source: fmt.Sprintf("application-active"),
					// TODO 設定ファイルから読み込み
					Target: "/queue/active",
				},
			},
		},
		&network.NetworkingConfig{}, nil, app.FileName)
	if err != nil {
		return nil, err
	}
	return &created, nil
}

func CreateAndConnectContainer(cli *client.Client, app model.Application, networkID string) (*container.ContainerCreateCreatedBody, error) {
	created, err := CreateContainer(cli, app)
	if err != nil {
		return nil, err
	}

  err = network_.ConnectNetwork(cli, networkID, created.ID,&network.EndpointSettings{})
  if err != nil{
    return nil, err
  }
  return created, nil
}
