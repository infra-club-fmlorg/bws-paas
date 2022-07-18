package container

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func StartContainer(cli *client.Client, id string) error {
  if err := cli.ContainerStart(context.Background(),id, types.ContainerStartOptions{}); err !=nil{
    return err
  }
  return nil
}
