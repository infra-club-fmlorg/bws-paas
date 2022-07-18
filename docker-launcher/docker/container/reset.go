package container

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

/*
コンテナ名からコンテナ作成のため初期化する関数

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
     err = cli.ContainerRemove(context.Background(), container.ID,types.ContainerRemoveOptions{})
    if err != nil{
      return err
    }
  }

  return Stop(cli, container.ID)
}
