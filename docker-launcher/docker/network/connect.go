package network

import (
	"context"

	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

func ConnectNetwork(cli *client.Client, networkID string, containerID string, option *network.EndpointSettings) error {
	return cli.NetworkConnect(context.Background(), networkID, containerID, option)
}
