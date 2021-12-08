package nodemanager

import (
	"path"
	"runtime"

	"github.com/yakumo-saki/phantasma-flow/util"
)

var nodeManager *NodeManager

func GetInstance() *NodeManager {

	if nodeManager == nil {
		nodeManager = &NodeManager{}
		return nodeManager // without log first time
	}

	traceCaller()
	return nodeManager
}

func traceCaller() {
	log := util.GetLoggerWithSource("NodeManager", "GetInstance")

	_, file, line, ok := runtime.Caller(2)
	if !ok {
		log.Debug().Msgf("Caller unknown")
	}

	log.Debug().Msgf("NodeManager request instance from %s:%v", path.Base(file), line)

}
