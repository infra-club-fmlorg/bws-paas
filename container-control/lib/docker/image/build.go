package image

import (
	"container-controller/lib/application"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type DockerfileTemplate struct {
	ApplicationPath string
}

const (
	DOCKERFILE_BUILD_CONTEXT_PATH  = "Dockerfile"
	APPLICATION_BUILD_CONTEXT_PATH = "application"
)

func Build(cli *client.Client, app *application.ApplicationInfo) error {
	buildContext, err := ArchiveBuildContext(app)
	if err != nil {
		return err
	}

	// TODO コンテキストに関しては要修整
	name := fmt.Sprintf("%s:%s", app.AssembleContainerName(), "latest")
	res, err := cli.ImageBuild(
		context.Background(),
		buildContext,
		types.ImageBuildOptions{
			Remove:     true,
			Tags:       []string{name},
			Dockerfile: DOCKERFILE_BUILD_CONTEXT_PATH,
		},
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	io.Copy(os.Stdout, res.Body)
	return nil
}
