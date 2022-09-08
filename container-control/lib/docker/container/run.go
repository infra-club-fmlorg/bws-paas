package container

import (
	"container-controller/lib/application"
	"container-controller/lib/docker/image"
	"container-controller/lib/docker/network"
	"container-controller/lib/file"
	"fmt"
	"log"
	"os"

	"github.com/docker/docker/client"
)

/*
Docker Containerを起動する処理をまとめた関数
*/
func Run(cli *client.Client, networkID string, app *application.ApplicationInfo) error {
	containerName := app.AssembleContainerName()
	incomingAppPath := fmt.Sprintf("%s/%s", app.AssembleIncomingDirPath(), app.AssembleFileName())
	activeAppPath := app.AssembleActiveAppPath()
	log.Printf("%#v\n", app)

	err := ResetByName(cli, containerName)
	if err != nil {
		return err
	}
	log.Printf("reset container: %s\n", containerName)

	err = file.Copy(incomingAppPath, activeAppPath)
	if err != nil {
		return err
	}
	log.Printf("copy to %s from %s\n", activeAppPath, incomingAppPath)
	err = os.Chmod(activeAppPath, 0100)
	if err != nil {
		return err
	}

	err = image.Build(cli, app)
	if err != nil {
		return err
	}

	created, err := CreateFromImage(cli, app)
	if err != nil {
		return err
	}

	err = network.ConnectContainer(cli, networkID, created.ID, nil)
	if err != nil {
		return err
	}
	log.Printf("create container connected network(%s): %s(%s)\n", networkID, containerName, created.ID)

	err = Start(cli, created.ID)
	if err != nil {
		return err
	}
	log.Printf("start container: %s(%s)\n", containerName, created.ID)

	// err = os.Remove(incomingAppPath)
	// if err != nil {
	// 	return err
	// }
	// log.Printf("remove %s\n", incomingAppPath)
	//
	return nil

}
