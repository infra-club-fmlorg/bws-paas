package image

import (
	"archive/tar"
	"container-controller/lib/application"
	"fmt"
)

func GenerateDockerfile(app *application.ApplicationInfo) ([]byte, *tar.Header, error) {
	var (
		dockerfile       []byte
		dockerfileHeader *tar.Header
		err              error
	)

	switch app.Runtime {
	case application.BINARY:
		dockerfile, dockerfileHeader, err = generateBinaryDockerfile(app)
	case application.NODE_JS:
		dockerfile, dockerfileHeader, err = generateNodeJSDockerfile(app)
	default:
		return nil, nil, fmt.Errorf("error: runtime not supported")
	}

	if err != nil {
		return nil, nil, err
	}

	return dockerfile, dockerfileHeader, nil

}
