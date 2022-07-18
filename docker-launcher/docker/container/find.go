package container

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func FindContainerByName(cli *client.Client, name string) (*types.Container, bool, error) {
	nameFilter := filters.NewArgs()
	nameFilter.Add("name", name)
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{
		All:     true,
		Filters: nameFilter,
	})
	if err != nil {
		return nil, false, err
	}

	for _, container := range containers {
		for _, containerName := range container.Names {
			if containerName == name {
				return &container, true, nil
			}
		}

	}
	return nil, false, err
}
