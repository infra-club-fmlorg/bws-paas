package container

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

/*
コンテナ名から初期化する関数

同名のコンテナが存在すれば停止させた上で削除する。

引数
cli - Dockerクライアント
containerName - Dockerコンテナの名前

返り値
error
*/
func ResetByName(cli *client.Client, containerName string) error {
	container, exist, err := FindByName(cli, containerName)
	if err != nil {
		return err
	}
	if !exist {
		return nil
	}

  if !IsRemoval(*container){
    err = Stop(cli, container.ID)
    if err != nil{
      return err
    }
  }

  return  cli.ContainerRemove(context.Background(), container.ID,types.ContainerRemoveOptions{})
}
