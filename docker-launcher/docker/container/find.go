package container

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

/*
名前からDockerコンテナを検索する関数

引数
cli - Dockerクライアント
containerName - Dockerコンテナの名前

返り値
container - 検索したコンテナのポインタ
exist - 対象のコンテナが存在しているか
error
*/
func FindByName(cli *client.Client, containerName string) (*types.Container, bool, error) {
	nameFilter := filters.NewArgs()
	nameFilter.Add("name", containerName)
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{
		All:     true,
		Filters: nameFilter,
	})
	if err != nil {
		return nil, false, err
	}

  // コンテナは複数の名前を持てるため、繰り返しで一致するまで探索
	for _, container := range containers {
		for _, containerName := range container.Names {
			if containerName == containerName {
				return &container, true, nil
			}
		}

	}
	return nil, false, err
}
