package main

import (
	"container-controller/lib/application"
	"container-controller/lib/docker/container"
	"container-controller/lib/docker/network"
	unixdomainsocket "container-controller/lib/unix-domain-socket"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/docker/docker/client"
)

const (
	DOCKER_LAUNCHER_SOCK = "/socket/docker_launcher.sock"
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
	// incomingDirPath := "/queue/incoming"

	// 実行したアプリケーションの保存先のルート
	// activedDirPath := "/queue/active"

	if _, err := os.Stat(DOCKER_LAUNCHER_SOCK); err == nil {
		if err := os.RemoveAll(DOCKER_LAUNCHER_SOCK); err != nil {
			log.Panic(err)
		}
	}

	listener, err := net.Listen("unix", DOCKER_LAUNCHER_SOCK)
	if err != nil {
		log.Panicln(err)
		return
	}
	defer listener.Close()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		log.Println("container luancher system launched")
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Println(err)
				continue
			}

			appInfo, err := func() (*application.ApplicationInfo, error) {
				defer conn.Close()

				m := &unixdomainsocket.Echo{}
				err := m.Read(conn)
				if err != nil {
					log.Println(err)
					return nil, err
				}

				appInfo := &application.ApplicationInfo{}
				err = json.Unmarshal(m.Data, appInfo)
				if err != nil {
					log.Println(err)
					return nil, err
				}
				return appInfo, nil
			}()

			if err != nil {
				log.Println(err)
				continue
			}

			err = container.Run(cli, dockerNetwork.ID, appInfo)
			if err != nil {
				log.Println(err)
				continue
			}
		}
	}()

	sig := <-quit
	log.Println(sig)
	return

	/*
		アプリケーションの本体
		無限ループの中でWalkDir関数を実行し、新規ファイルの検索を行う
	*/
	// log.Printf("start to walk directory:%s", incomingDirPath)
	// handler := createHandleWalkDir(cli, dockerNetwork.ID)
	// for {
	// 	err := filepath.WalkDir(incomingDirPath, handler)
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// 	// 1秒間隔で全検査
	// 	time.Sleep(1 * time.Second)
	// }
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

		err = container.Run(cli, networkID, app)
		if err != nil {
			return err
		}

		return nil
	}
}
