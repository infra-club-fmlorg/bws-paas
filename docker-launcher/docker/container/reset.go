package container

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func ResetContainerByName(cli *client.Client, name string) error {
	container, exist, err := FindContainerByName(cli, name)
	if err != nil {
		return err
	}
	if !exist {
		return nil
	}

  if !IsRemoval(*container){
     err = cli.ContainerRemove(context.Background(), container.ID,types.ContainerRemoveOptions{})
    if err != nil{
      return err
    }
  }

  return StopContainer(cli, container.ID)
}
