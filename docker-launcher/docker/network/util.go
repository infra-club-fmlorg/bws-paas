package network

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func FindNetworkByName(cli *client.Client, name string) (*types.NetworkResource, bool, error) {
	networkNameFilter := filters.NewArgs()
	networkNameFilter.Add("name", name)
	networks, err := cli.NetworkList(context.Background(), types.NetworkListOptions{
		Filters: networkNameFilter,
	})
	if err != nil {
		return nil, false, err
	}

	for _, network := range networks {
		if network.Name == name {
			return &network, true, nil
		}
	}
	return nil, false, err
}
