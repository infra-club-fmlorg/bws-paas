package network

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

/*
名前からDockerネットワークを検索する関数

引数
cli - Dockerクライアント
networkID - Dockerネットワークの名前

返り値
network - 検索したネットワークのポインタ
exist - 対象のネットワークが存在しているか
error
*/
func FindByName(cli *client.Client, networkID string) (*types.NetworkResource, bool, error) {
	networkNameFilter := filters.NewArgs()
	networkNameFilter.Add("name", networkID)
	networks, err := cli.NetworkList(context.Background(), types.NetworkListOptions{
		Filters: networkNameFilter,
	})
	if err != nil {
		return nil, false, err
	}

	for _, network := range networks {
		if network.Name == networkID {
			return &network, true, nil
		}
	}
	return nil, false, err
}
