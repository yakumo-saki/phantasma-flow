package repository

import (
	"errors"
	"io/ioutil"
	"path/filepath"

	"github.com/goccy/go-yaml"
	"github.com/yakumo-saki/phantasma-flow/messagehub"
	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
	"github.com/yakumo-saki/phantasma-flow/util"
)

const myname = "Repository"

type Repository struct {
	nodes   []objects.NodeDefinition
	jobs    []objects.JobDefinition
	configs []objects.Config
}

func (r *Repository) Initialize(path string) error {
	log := util.GetLoggerWithSource(myname, "initialize")
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
	log := util.GetLoggerWithSource(myname, "dump")
	log.Debug().Msg("Jobs")
	for _, v := range repo.jobs {
		log.Debug().Msgf("%s", v)
		log.Debug().Msg("")
	}
	log.Debug().Msg("-------------------------------------")
	log.Debug().Msg("Nodes")
	for _, v := range repo.nodes {
		log.Debug().Msgf("%s", v)
		log.Debug().Msg("")
	}
	log.Debug().Msg("-------------------------------------")
	log.Debug().Msg("Configs")
	for _, v := range repo.configs {
		log.Debug().Msgf("%s", v)
		log.Debug().Msg("")
	}
	log.Debug().Msg("-------------------------------------")
	return false
}

func (repo *Repository) readAllYaml(path string, objType objectType) error {
	log := util.GetLoggerWithSource(myname, "readYaml")
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

func (repo *Repository) SendAllNodes() {
	for _, v := range repo.nodes {
		messagehub.Post(messagehub.TOPIC_NODE_DEFINITION, v)

	}
}

func (repo *Repository) SendAllJobs() {
	for _, v := range repo.jobs {
		messagehub.Post(messagehub.TOPIC_JOB_DEFINITION, v)
	}
}
