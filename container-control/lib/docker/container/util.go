package container

import (
	"github.com/docker/docker/api/types"
)

/*
コンテナが削除可かチェックする関数

引数
container - コンテナの構造体

返り値
isRemoval - コンテナが削除可能か
*/
func IsRemoval(container types.Container) bool {
	// ["created", "running", "paused", "restarting", "removing", "exited", "dead"]
	return container.State == "exited" || container.State == "dead"
}
