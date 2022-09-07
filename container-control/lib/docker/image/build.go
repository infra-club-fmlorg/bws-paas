package image

import (
	"archive/tar"
	"bytes"
	"container-controller/lib/application"
	"context"
	"embed"
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

//go:embed static/dockerfile/binary.Dockerfile
var dockerfiles embed.FS

type DockerfileTemplate struct {
	ApplicationPath string
}

func Build(cli *client.Client, app *application.ApplicationInfo) error {
	dockerfile, err := getArchivedDockerfile(app)
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

func getArchivedDockerfile(app *application.ApplicationInfo) (*bytes.Reader, error) {
	buf := new(bytes.Buffer)
	t, err := template.ParseFS(dockerfiles, "static/dockerfile/binary.Dockerfile")
	if err != nil {
		return nil, err
	}
	t.Execute(buf, DockerfileTemplate{
		ApplicationPath: app.AssembleActiveAppPath(),
	})

	// archive the Dockerfile
	tarHeader := &tar.Header{
		Name: "binary",
		Size: int64(buf.Len()),
	}
	tw := tar.NewWriter(buf)
	defer tw.Close()
	err = tw.WriteHeader(tarHeader)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(buf.Bytes()), nil
}
