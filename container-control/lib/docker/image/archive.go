package image

import (
	"archive/tar"
	"bytes"
	"container-controller/lib/application"
	"os"
)

func ArchiveBuildContext(app *application.ApplicationInfo) (*bytes.Reader, error) {
	dockerfile, dockerfileHeader, err := GenerateDockerfile(app)
	if err != nil {
		return nil, err
	}

	application, err := os.ReadFile(app.AssembleActiveAppPath())
	applicationHeader := &tar.Header{
		Name: APPLICATION_BUILD_CONTEXT_PATH,
		Size: int64(len(application)),
	}

	// archive the Dockerfile
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()

	err = tw.WriteHeader(dockerfileHeader)
	if err != nil {
		return nil, err
	}
	_, err = tw.Write(dockerfile)
	if err != nil {
		return nil, err
	}

	err = tw.WriteHeader(applicationHeader)
	if err != nil {
		return nil, err
	}
	_, err = tw.Write(application)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(buf.Bytes()), nil
}
