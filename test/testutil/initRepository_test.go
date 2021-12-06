package testutil

import (
	"fmt"
	"runtime"

	"github.com/rs/zerolog/log"
	"github.com/yakumo-saki/phantasma-flow/repository"
)

func StartRepository() *repository.Repository {

	fmt.Println(runtime.Caller(0))
	return nil

	repo := repository.GetRepository()
	err := repo.Initialize()
	if err != nil {
		log.Error().Err(err).Msg("Error occured at reading initialize data")
		log.Error().Msg("Maybe data is corrupted or misseditted.")
		return nil
	}

	return repo
}
