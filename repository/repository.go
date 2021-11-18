package repository

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/goccy/go-yaml"
	"github.com/yakumo-saki/phantasma-flow/objects"
	"github.com/yakumo-saki/phantasma-flow/util"
)

type Repository struct {
	nodes   []objects.NodeDefinition
	jobs    []objects.JobDefinition
	configs []objects.Config
}

func (r *Repository) Initialize(path string) error {
	log := util.GetLogger()
	log.Debug().Msg("Repository initialize")

	dirType := map[objectType][]string{
		NODE:   {"node"},
		CONFIG: {"config"},
		JOB:    {"job", "definitions"},
	}

	for typ, pt := range dirType {
		readDirPath := util.JoinPath(path, pt)
		log.Debug().Msgf("Reading %s from %s", typ, readDirPath)
		err := r.readAllYaml(readDirPath, typ)
		if err != nil {
			return err
		}
	}

	if r.Dump() {
		log.Debug().Msg("Repository initialized")
		return errors.New("not impremented")
	}

	return nil
}

func (repo *Repository) Dump() bool {
	fmt.Println("Jobs")
	for _, v := range repo.jobs {
		fmt.Println(v)
		fmt.Println("")
	}
	fmt.Println("-------------------------------------")
	fmt.Println("Nodes")
	for _, v := range repo.nodes {
		fmt.Println(v)
		fmt.Println("")
	}
	fmt.Println("-------------------------------------")
	fmt.Println("Configs")
	for _, v := range repo.configs {
		fmt.Println(v)
		fmt.Println("")
	}
	fmt.Println("-------------------------------------")
	return false
}

func (repo *Repository) readAllYaml(path string, objType objectType) error {
	log := util.GetLogger()
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Err(err)
	}

	for _, fileinfo := range files {
		if fileinfo.IsDir() {
			continue
		}

		fp := filepath.Join(path, fileinfo.Name())
		bytes, err := ioutil.ReadFile(fp)
		if err != nil {
			return err
		}

		log.Debug().Msgf("Reading %s", fp)

		switch objType {
		case NODE:
			obj := objects.NodeDefinition{}
			err := yaml.Unmarshal(bytes, &obj)
			if err != nil {
				return err
			}
			repo.nodes = append(repo.nodes, obj)
		case JOB:
			obj := objects.JobDefinition{}
			err := yaml.Unmarshal(bytes, &obj)
			if err != nil {
				return err
			}
			repo.jobs = append(repo.jobs, obj)
		case CONFIG:
			obj := objects.Config{}
			err := yaml.Unmarshal(bytes, &obj)
			if err != nil {
				return err
			}
			repo.configs = append(repo.configs, obj)
		}

	}

	return nil
}
