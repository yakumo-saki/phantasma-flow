package testutils

import (
	"fmt"
	"os"
	"path"

	"github.com/rs/zerolog/log"
	"github.com/yakumo-saki/phantasma-flow/messagehub"
	"github.com/yakumo-saki/phantasma-flow/messagehub/messagehub_impl"
	"github.com/yakumo-saki/phantasma-flow/procman"
	"github.com/yakumo-saki/phantasma-flow/repository"
)

func StartBaseModules() (*messagehub_impl.MessageHub, *procman.ProcessManager) {
	StartRepository()
	hub := messagehub_impl.MessageHub{}
	hub.Initialize()
	messagehub.SetMessageHub(&hub)

	pman := procman.NewProcessManager(make(chan string, 1))
	pman.Add(&TestLogListener{})

	return &hub, &pman

}

func StartRepository() *repository.Repository {

	dir := GetTestBaseDir()
	home := path.Join(dir, "phantasma-flow")
	fmt.Printf("SET PHFLOW_HOME = %s\n", home)
	os.Setenv("PHFLOW_HOME", home)

	repo := repository.GetRepository()
	err := repo.Initialize()
	if err != nil {
		log.Error().Err(err).Msg("Error occured at reading initialize data")
		log.Error().Msg("Maybe data is corrupted or misseditted.")
		return nil
	}

	return repo

}
