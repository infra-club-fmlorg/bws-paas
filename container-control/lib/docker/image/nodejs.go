package image

import (
	"archive/tar"
	"bytes"
	"container-controller/lib/application"
	"embed"
	"io/ioutil"
	"text/template"
)

//go:embed dockerfile/nodejs.Dockerfile
var nodeJSDockerfiles embed.FS

const (
	NODEJS_DOCKERFILE_PATH = "dockerfile/nodejs.Dockerfile"
)

type nodeJSDockerfileTemplate struct {
	ApplicationName string
	ApplicationPath string
}

func generateNodeJSDockerfile(app *application.ApplicationInfo) ([]byte, *tar.Header, error) {
	template, err := template.ParseFS(nodeJSDockerfiles, NODEJS_DOCKERFILE_PATH)
	if err != nil {
		return nil, nil, err
	}

	templateBuf := new(bytes.Buffer)
	template.Execute(templateBuf, nodeJSDockerfileTemplate{
		ApplicationName: app.ApplicationName,
		ApplicationPath: app.AssembleActiveAppPath(),
	})
	dockerfile, err := ioutil.ReadAll(templateBuf)
	if err != nil {
		return nil, nil, err
	}

	dockerfileHeader := &tar.Header{
		Name: DOCKERFILE_BUILD_CONTEXT_PATH,
		Size: int64(len(dockerfile)),
	}
	return dockerfile, dockerfileHeader, nil
}
