package executer

import (
	"context"

	"github.com/yakumo-saki/phantasma-flow/pkg/message"
	"github.com/yakumo-saki/phantasma-flow/util"
)

func (e *Executer) Run(parentCtx context.Context) {
	log := util.GetLoggerWithSource("Executer " + e.Job.Id)
	log.Info().Msg("Run job step (dummy)")

}

func (e *Executer) RunInitialize(parentCtx context.Context) {
	// SSH Connection
}

func (e *Executer) RunStep(parentCtx context.Context) {
	execMsg := message.ExecuterMsg{}
	execMsg.JobId = "hoge"
}
