package container

import (
	"context"
	"time"

	"github.com/docker/docker/client"
)

func StopContainer(cli *client.Client, id string) error {
	timeout := time.Duration(500) * time.Millisecond
	return cli.ContainerStop(
		context.Background(),
		id,
		&timeout,
	)
}
