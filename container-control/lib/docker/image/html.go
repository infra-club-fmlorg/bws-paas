package image

import (
	"archive/tar"
	"bytes"
	"container-controller/lib/application"
	"embed"
	"io/ioutil"
	"text/template"
)

//go:embed dockerfile/html.Dockerfile
var htmlDockerfiles embed.FS

const (
	HTML_DOCKERFILE_PATH = "dockerfile/html.Dockerfile"
)

type htmlDockerfileTemplate struct {
	ApplicationName string
	ApplicationPath string
}

func generateHTMLDockerfile(app *application.ApplicationInfo) ([]byte, *tar.Header, error) {
	template, err := template.ParseFS(htmlDockerfiles, HTML_DOCKERFILE_PATH)
	if err != nil {
		return nil, nil, err
	}

	templateBuf := new(bytes.Buffer)
	template.Execute(templateBuf, htmlDockerfileTemplate{
		ApplicationName: app.ApplicationName,
		ApplicationPath: APPLICATION_BUILD_CONTEXT_PATH,
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
