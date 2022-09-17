package image

import (
	"archive/tar"
	"bytes"
	"container-controller/lib/application"
	_ "embed"
	"os"
)

func ArchiveBuildContext(app *application.ApplicationInfo) (*bytes.Reader, error) {
	// archive the Dockerfile
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()

	err := GenerateDockerfile(tw, app)
	if err != nil {
		return nil, err
	}

	applicationFile, err := os.ReadFile(app.AssembleActiveAppPath())
	applicationHeader := &tar.Header{
		Name: APPLICATION_BUILD_CONTEXT_PATH,
		Size: int64(len(applicationFile)),
	}

	err = tw.WriteHeader(applicationHeader)
	if err != nil {
		return nil, err
	}
	_, err = tw.Write(applicationFile)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(buf.Bytes()), nil
}
