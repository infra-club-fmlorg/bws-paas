package main

import (
	"docker-launcher/lib/application"
	"docker-launcher/lib/docker/container"
	"docker-launcher/lib/docker/network"
	"docker-launcher/lib/file"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/docker/client"
)

// コマンドライン引数の構造体
type Flag struct {
	network string // 生成したコンテナを所属させるDocker Network Name
}

func main() {
	// コマンドライン引数のパース
	var myFlag Flag
	flag.StringVar(&myFlag.network, "network", "", "Docker Network Name")
	flag.Parse()

	//コマンドライン引数のバリデーション
	if len(myFlag.network) == 0 {
		log.Panicln(fmt.Errorf("error: the following required arguments were not provided: --network={Docker Network Name}"))
		return
	}

	// Dockerクライアントの生成
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Panicln(err)
		return
	}

	// 利用するDockerネットワークの取得
	dockerNetwork, exist, err := network.FindByName(cli, myFlag.network)
	if err != nil {
		log.Panicln(err)
		return
	}
	if !exist {
		log.Panicf("not found network: %s\n", myFlag.network)
		return
	}

	// アプリケーションの一時保存先のルート
	incomingDirPath := "/queue/incoming"

	// 実行したアプリケーションの保存先のルート
	// activedDirPath := "/queue/active"

	/*
		アプリケーションの本体
		無限ループの中でWalkDir関数を実行し、新規ファイルの検索を行う
	*/
	log.Printf("start to walk directory:%s", incomingDirPath)
	handler := createHandleWalkDir(cli, dockerNetwork.ID)
	for {
		err := filepath.WalkDir(incomingDirPath, handler)
		if err != nil {
			log.Println(err)
		}
		// 1秒間隔で全検査
		time.Sleep(1 * time.Second)
	}
}

/*
Dockerクライアントを受け取って、WalkDir用の関数を返す高階関数

引数
cli *client.Client -- Dockerクライアント
networkID - DockerネットワークのID

返り値
handleWalkDir func(path string, entry fs.DirEntry, err error) error -- WalkDir関数ハンドラ
	引数
	path string -- 現在のディレクトリ
	entry fs.DirEntry -- ファイル及びディレクトリの情報
	err error -- 実行時エラー
	返り値
	error -- 実行時例外
*/
func createHandleWalkDir(cli *client.Client, networkID string) func(path string, entry fs.DirEntry, err error) error {
	// 無名関数を返す
	return func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("Error:%s", err)
		}

		if entry.IsDir() {
			return nil
		}

		// 処理対象がファイルだった場合
		if !entry.Type().IsRegular() {
			return nil
		}

		log.Println(path)

		app, err := application.ParseApplicationInfoFromPath(path)
		if err != nil {
			return err
		}

		containerName := app.AssembleContainerName()
		incomingAppPath := fmt.Sprintf("%s/%s", app.AssembleIncomingDirPath(), app.AssembleFileName())
		activeAppPath := app.AssembleActiveAppPath()
		log.Printf("%#v\n", app)

		container.ResetByName(cli, containerName)
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
			return nil
		}

		created, err := container.CreateConnectedNetwork(cli, *app, networkID)
		if err != nil {
			return err
		}
		log.Printf("create container connected network(%s): %s(%s)\n", networkID, containerName, created.ID)

		err = container.Start(cli, created.ID)
		if err != nil {
			return err
		}
		log.Printf("start container: %s(%s)\n", containerName, created.ID)

		err = os.Remove(incomingAppPath)
		if err != nil {
			return err
		}
		log.Printf("remove %s\n", incomingAppPath)

		return nil
	}
}
