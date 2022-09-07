package image

import (
	"archive/tar"
	"bytes"
	"container-controller/lib/application"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// go:embed static/dockerfile/binary
var binaryDockerfile []byte

func Build(cli *client.Client, app *application.ApplicationInfo) error {
	dockerfile, err := getArchivedDockerfile()
	if err != nil {
		return err
	}

	name := fmt.Sprintf("%s:%s", app.AssembleContainerName(), "latest")
	res, err := cli.ImageBuild(
		context.Background(),
		dockerfile,
		types.ImageBuildOptions{
			Remove:     true,
			Tags:       []string{name},
			Dockerfile: "binary",
		},
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	io.Copy(os.Stdout, res.Body)
	return nil
}

func getArchivedDockerfile() (*bytes.Reader, error) {
	// archive the Dockerfile
	tarHeader := &tar.Header{
		Name: "binary",
		Size: int64(len(binaryDockerfile)),
	}
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()
	err := tw.WriteHeader(tarHeader)
	if err != nil {
		return nil, err
	}
	_, err = tw.Write(binaryDockerfile)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(buf.Bytes()), nil
}
