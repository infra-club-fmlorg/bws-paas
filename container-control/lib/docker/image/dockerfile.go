package image

import (
	"archive/tar"
	"container-controller/lib/application"
	"fmt"
)

func GenerateDockerfile(tw *tar.Writer, app *application.ApplicationInfo) error {
	var (
		err error
	)

	switch app.Runtime {
	case application.BINARY:
		err = generateBinaryDockerfile(tw, app)
	case application.NODE_JS:
		err = generateNodeJSDockerfile(tw, app)
	case application.HTML:
		err = generateHTMLDockerfile(tw, app)
	default:
		return fmt.Errorf("error: runtime not supported")
	}

	if err != nil {
		return err
	}

	return nil

}
