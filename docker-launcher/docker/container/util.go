package container

import (
	"github.com/docker/docker/api/types"
)

func IsRemoval(container types.Container) bool{
  // ["created", "running", "paused", "restarting", "removing", "exited", "dead"]
  return container.State == "exited" || container.State == "dead"
}
