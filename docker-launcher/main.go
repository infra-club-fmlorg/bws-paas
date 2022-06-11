package main

import (
	"context"
	"flag"
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

/*
コマンドライン引数のオプションの構造体

network string -- 生成したコンテナを所属させるDocker Network名
*/
type Flag struct {
	network string
}

func main() {
	var myFlag Flag
	// コマンドライン引数のパース
	flag.StringVar(&myFlag.network, "network", "", "Docker Network Name")
	flag.Parse()

	//コマンドライン引数のバリデーション
	if len(myFlag.network) == 0 {
		log.Fatalln(fmt.Errorf("error: The following required arguments were not provided: \"network name\""))
		return
	}

	// Dockerクライアントの生成
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Panicln(err)
		return
	}

	// アプリケーションの一時保管先
	incomingPath := "/queue/incoming"
	// activedPath := "queue/active"

	/*
		アプリケーションの本体

		無限ループの中でWalkDir関数を実行し、新規ファイルの検索を行う

		TODO 無限ループに実行間隔を追加する
	*/
	log.Println("start to walk directory: " + incomingPath)
	for {
		err := filepath.WalkDir(incomingPath, createHandleWatchDir(cli))
		if err != nil {
			log.Println(err)
		}
	}
}

/*
Dockerクライアントを受け取って、WalkDir用の関数を返す高階関数

引数
cli *client.Client -- Dockerクライアント

返り値
handlerWalkDir func(path string, entry fs.DirEntry, err error) error -- WalkDir関数ハンドラ
	引数
	path string -- 現在のディレクトリ
	entry fs.DirEntry -- ファイル及びディレクトリの情報
	err error -- 実行時エラー
	返り値
	error -- 実行時エラー
*/
func createHandleWatchDir(cli *client.Client) func(path string, entry fs.DirEntry, err error) error {
	// 無名関数を返す
	return func(path string, entry fs.DirEntry, err error) error {
		// 実行時エラーが発生した場合
		if err != nil {
			return fmt.Errorf("Error:%s", err)
		}

		// 処理対象がディレクトリだった場合
		if entry.IsDir() {
			return nil
		}

		// 処理対象がファイルだった場合
		if entry.Type().IsRegular() {
			log.Println(path)

			// Docker Containerの起動
			activation(cli, path)
			return nil
		}

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
