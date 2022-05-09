package main

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

func main() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Println(err)
		return
	}

	if len(os.Args) == 1 {
		log.Println("fail to get docker's network name from command line arguments")
		return
	}
	log.Println("network name: " + os.Args[1])

	path, err := filepath.Abs("/queue/incoming")
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("start to walk directory: " + path)
	for {
		err := filepath.WalkDir(path, watchDir(cli))
		if err != nil {
			log.Println(err)
		}
	}
}

func watchDir(cli *client.Client) func(path string, entry fs.DirEntry, err error) error {
	return func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("Error:%s", err)
		}

		if entry.IsDir() {
			return nil
		}

		log.Println(path)
		activation(cli, path)
		return nil
	}
}

func activation(cli *client.Client, path string) error {
	// log.Println(path)
	userName := filepath.Base(filepath.Dir(filepath.Dir(path)))
	// log.Println(userName)
	applicationName := filepath.Base(filepath.Dir(path))
	// log.Println(applicationName)
	applicationFileName := time.Now().UTC().Format(time.RFC3339Nano)

	incomingApplicationFileInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		log.Printf("not found %s\n", path)
		return err
	}
	// log.Printf("file size: %d b\n", fileInfoTarget.Size())

	activeApplicationDirPath, err := filepath.Abs(fmt.Sprintf("/queue/active/%s/%s", userName, applicationName))
	if err != nil {
		log.Println(err)
		return err
	}
	// log.Println(destinationParentDirName)
	activeApplicationPath := filepath.Join(activeApplicationDirPath, applicationFileName)
	// log.Println(destinationPath)

	err = os.MkdirAll(activeApplicationDirPath, 0777)
	if err != nil {
		log.Printf("fail to mkdir %s\n", activeApplicationDirPath)
		return err
	}
	log.Println("success to mkdir: " + activeApplicationDirPath)

	// err = os.Link(path, activeApplicationPath)
	// if err != nil {
	// 	log.Printf("fail to cp from %s to %s\n", path, activeApplicationPath)
	// 	return err
	// }
	newFile, err := os.Create(activeApplicationPath)
	if err != nil {
		log.Printf("fail to cp from %s to %s", path, activeApplicationPath)
		return err
	}

	originalFile, err := os.Open(path)
	if err != nil {
		log.Printf("fail to cp from %s to %s", path, activeApplicationPath)
		return err
	}

	_, err = io.Copy(newFile, originalFile)
	if err != nil {
		log.Printf("fail to cp from %s to %s", path, activeApplicationPath)
		return err
	}
	newFile.Close()
	originalFile.Close()
	log.Printf("success to cp from %s to %s", path, activeApplicationPath)

	activeApplicationFileInfo, err := os.Stat(activeApplicationPath)
	if os.IsNotExist(err) {
		log.Printf("not found %s\n", activeApplicationPath)
		return err
	}
	// log.Printf("file size: %d b\n", fileInfoTarget.Size())

	if activeApplicationFileInfo.Size() != incomingApplicationFileInfo.Size() {
		err = os.RemoveAll(activeApplicationDirPath)
		if err != nil {
			log.Printf("fail to remove %s\n", activeApplicationDirPath)
		}
		log.Println("not equal file size")
		return fmt.Errorf("not equal file size")
	}
	// log.Println("equal file size")

	err = os.Chmod(activeApplicationPath, 0777)
	if err != nil {
		log.Println(err)
		return err
	}

	defer func() {
		err := os.Remove(path)
		if err != nil {
			log.Printf("fail to remove %s\n", path)
		}
		log.Printf("success to remove %s\n", path)
	}()

	err = startContainer(cli, userName, applicationName, applicationFileName)
	if err != nil {
		return err
	}

	return nil
}

func startContainer(cli *client.Client, userName string, applicationName string, fileName string) error {
	containerName := userName + "-" + applicationName
	containerNameFilter := filters.NewArgs()
	containerNameFilter.Add("name", containerName)
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{
		All:     true,
		Filters: containerNameFilter,
	})
	if err != nil {
		return err
	}

	if len(containers) != 0 {
		if containers[0].State == "restarting" || containers[0].State == "running" {
			timeout := time.Duration(500) * time.Millisecond
			err := cli.ContainerStop(
				context.Background(),
				containers[0].ID,
				&timeout,
			)
			if err != nil {
				log.Println("fail to stop container: " + containers[0].ID)
				return err
			}
		}
		log.Println("success to stop container: " + containers[0].ID)

		err := cli.ContainerRemove(
			context.Background(),
			containers[0].ID,
			types.ContainerRemoveOptions{},
		)
		if err != nil {
			log.Println("fail to remove container: " + containers[0].ID)
			return err
		}
		log.Println("success to remove container: " + containers[0].ID)
	}

	executedUser := os.Getenv("EXECUTED_USER")
	entryPoint := fmt.Sprintf("/%s-application-active/%s/%s/%s", executedUser, userName, applicationName, fileName)

	result, err := cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image:      "ubuntu",
			Entrypoint: []string{entryPoint},
		},
		&container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeVolume,
					Source: fmt.Sprintf("%s-application-active", executedUser),
					Target: fmt.Sprintf("/%s-application-active", executedUser),
				},
			},
		},
		&network.NetworkingConfig{}, nil, containerName)
	if err != nil {
		log.Println("fail to create container: " + entryPoint)
		return err
	}
	log.Println("success to create container: " + result.ID)

	networkNameFilter := filters.NewArgs()
	networkNameFilter.Add("name", os.Args[1])
	networks, err := cli.NetworkList(context.Background(), types.NetworkListOptions{
		Filters: networkNameFilter,
	})
	if err != nil {
		return err
	}

	err = cli.NetworkConnect(context.Background(), networks[0].ID, result.ID, &network.EndpointSettings{})
	if err != nil {
		log.Printf("fail to connect network(%s): %s", networks[0].ID, result.ID)
		return err
	}

	err = cli.ContainerStart(context.Background(), result.ID, types.ContainerStartOptions{})
	if err != nil {
		log.Println("fail to launch container: " + result.ID)
		return err
	}

	log.Println("success to launch container: " + result.ID)

	return nil
}
