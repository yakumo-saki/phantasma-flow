package executer

import (
	"context"

	"github.com/yakumo-saki/phantasma-flow/util"
)

func (e *Executer) Run(parentCtx context.Context) {
	log := util.GetLoggerWithSource("Executer " + e.Job.Id)
	log.Info().Msg("Run job (dummy)")
}

func (e *Executer) RunInitialize(parentCtx context.Context) {
	// SSH Connection
}

func (e *Executer) RunStep(parentCtx context.Context) {

}
