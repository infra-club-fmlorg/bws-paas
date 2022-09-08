package image

import (
	"archive/tar"
	"bytes"
	"container-controller/lib/application"
	_ "embed"
	"os"
)

//go:embed static/default.conf
var defaultConf []byte

func ArchiveBuildContext(app *application.ApplicationInfo) (*bytes.Reader, error) {
	dockerfile, dockerfileHeader, err := GenerateDockerfile(app)
	if err != nil {
		return nil, err
	}

	applicationFile, err := os.ReadFile(app.AssembleActiveAppPath())
	applicationHeader := &tar.Header{
		Name: APPLICATION_BUILD_CONTEXT_PATH,
		Size: int64(len(applicationFile)),
	}

	defaultConfHeader := &tar.Header{
		Name: "default.conf",
		Size: int64(len(defaultConf)),
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
	_, err = tw.Write(applicationFile)
	if err != nil {
		return nil, err
	}

	if app.Runtime == application.HTML {
		err = tw.WriteHeader(defaultConfHeader)
		if err != nil {
			return nil, err
		}
		_, err = tw.Write(defaultConf)
		if err != nil {
			return nil, err
		}
	}

	return bytes.NewReader(buf.Bytes()), nil
}
