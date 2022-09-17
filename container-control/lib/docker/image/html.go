package image

import (
	"archive/tar"
	"bytes"
	"container-controller/lib/application"
	"embed"
	"io/ioutil"
	"text/template"
)

var (
	//go:embed static/default.conf
	htmlDefaultConf []byte

	//go:embed dockerfile/html.Dockerfile
	htmlDockerfiles embed.FS
)

const (
	HTML_DOCKERFILE_PATH   = "dockerfile/html.Dockerfile"
	HTML_DEFAULT_CONF_PATH = "default.conf"
)

type htmlDockerfileTemplate struct {
	ApplicationName string
	ApplicationPath string
}

func generateHTMLDockerfile(tw *tar.Writer, app *application.ApplicationInfo) error {
	template, err := template.ParseFS(htmlDockerfiles, HTML_DOCKERFILE_PATH)
	if err != nil {
		return err
	}

	templateBuf := new(bytes.Buffer)
	template.Execute(templateBuf, htmlDockerfileTemplate{
		ApplicationName: app.ApplicationName,
		ApplicationPath: APPLICATION_BUILD_CONTEXT_PATH,
	})
	dockerfile, err := ioutil.ReadAll(templateBuf)
	if err != nil {
		return err
	}

	dockerfileHeader := &tar.Header{
		Name: DOCKERFILE_BUILD_CONTEXT_PATH,
		Size: int64(len(dockerfile)),
	}
	err = tw.WriteHeader(dockerfileHeader)
	if err != nil {
		return err
	}
	_, err = tw.Write(dockerfile)
	if err != nil {
		return err
	}

	defaultConfHeader := &tar.Header{
		Name: HTML_DEFAULT_CONF_PATH,
		Size: int64(len(htmlDefaultConf)),
	}
	err = tw.WriteHeader(defaultConfHeader)
	if err != nil {
		return err
	}
	_, err = tw.Write(htmlDefaultConf)
	if err != nil {
		return err
	}

	return nil
}
