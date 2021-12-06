package testutils

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/rs/zerolog/log"
	"github.com/yakumo-saki/phantasma-flow/repository"
)

func StartRepository() *repository.Repository {
	_, file, _, _ := runtime.Caller(0)

	dir := path.Dir(file)
	for {
		if path.Base(dir) == "test" {
			break
		}
		dir, _ = path.Split(dir)
	}

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
