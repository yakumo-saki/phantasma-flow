package repository

import (
	"errors"
	"io/ioutil"
	"path/filepath"

	"github.com/goccy/go-yaml"
	"github.com/yakumo-saki/phantasma-flow/messagehub"
	"github.com/yakumo-saki/phantasma-flow/pkg/messagehubObjects"
	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
	"github.com/yakumo-saki/phantasma-flow/util"
)

const myname = "Repository"

type Repository struct {
	nodes   []objects.NodeDefinition
	jobs    []objects.JobDefinition
	configs []objects.Config

	paths phflowPath
}

func (r *Repository) Initialize() error {
	log := util.GetLoggerWithSource(myname, "initialize")
	log.Debug().Msg("Repository initialize")

	r.paths = aquirePhflowPath()

	dirType := map[objectType]string{
		NODE:   r.paths.NodeDef,
		CONFIG: r.paths.ConfigDef,
		JOB:    r.paths.JobDef,
	}

	for typ, dirPath := range dirType {
		log.Debug().Msgf("Reading %s from %s", typ, dirPath)
		err := r.readAllYaml(dirPath, typ)
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

func (repo *Repository) SendAllNodes() int {
	sent := 0
	for _, v := range repo.nodes {
		nodeMsg := messagehubObjects.NodeDefinitionMsg{}
		nodeMsg.Reason = messagehubObjects.DEF_REASON_INITIAL
		nodeMsg.NodeDefinition = v
		messagehub.Post(messagehub.TOPIC_NODE_DEFINITION, nodeMsg)
		sent++
	}
	return sent
}

func (repo *Repository) SendAllJobs() int {
	sent := 0
	for _, v := range repo.jobs {
		jobMsg := messagehubObjects.JobDefinitionMsg{}
		jobMsg.Reason = messagehubObjects.DEF_REASON_INITIAL
		jobMsg.JobDefinition = v
		messagehub.Post(messagehub.TOPIC_JOB_DEFINITION, jobMsg)
		sent++
	}
	return sent
}
