package image

import (
	"archive/tar"
	"bytes"
	"container-controller/lib/application"
	"embed"
	"io/ioutil"
	"text/template"
)

//go:embed dockerfile/binary.Dockerfile
var binaryDockerfiles embed.FS

const (
	BINARY_DOCKERFILE_PATH = "dockerfile/binary.Dockerfile"
)

type binaryDockerfileTemplate struct {
	ApplicationPath string
}

func generateBinaryDockerfile(tw *tar.Writer, app *application.ApplicationInfo) error {
	template, err := template.ParseFS(binaryDockerfiles, BINARY_DOCKERFILE_PATH)
	if err != nil {
		return err
	}

	templateBuf := new(bytes.Buffer)
	template.Execute(templateBuf, binaryDockerfileTemplate{
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

	return nil
}
