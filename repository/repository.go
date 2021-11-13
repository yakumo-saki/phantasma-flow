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

type repository struct {
	Nodes   []objects.NodeDefinition
	Jobs    []objects.JobDefinition
	Configs []objects.Config
}

var repo repository

func Initialize(path string) error {
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
		err := readAllYaml(readDirPath, typ)
		if err != nil {
			return err
		}
	}

	if dump() {
		log.Debug().Msg("Repository initialized")
		return errors.New("not impremented")
	}

	return nil
}

func dump() bool {
	fmt.Println("Jobs")
	for _, v := range repo.Jobs {
		fmt.Println(v)
		fmt.Println("")
	}
	fmt.Println("-------------------------------------")
	fmt.Println("Nodes")
	for _, v := range repo.Nodes {
		fmt.Println(v)
		fmt.Println("")
	}
	fmt.Println("-------------------------------------")
	fmt.Println("Configs")
	for _, v := range repo.Configs {
		fmt.Println(v)
		fmt.Println("")
	}
	fmt.Println("-------------------------------------")
	return false
}

func readAllYaml(path string, objType objectType) error {
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
			repo.Nodes = append(repo.Nodes, obj)
		case JOB:
			obj := objects.JobDefinition{}
			err := yaml.Unmarshal(bytes, &obj)
			if err != nil {
				return err
			}
			repo.Jobs = append(repo.Jobs, obj)
		case CONFIG:
			obj := objects.Config{}
			err := yaml.Unmarshal(bytes, &obj)
			if err != nil {
				return err
			}
			repo.Configs = append(repo.Configs, obj)
		}

	}

	return nil
}

func GetConfig() {}

func ApplyNode(nodeDef objects.NodeDefinition) {

}

func ApplyJob(nodeDef objects.JobDefinition) {

}

func ApplyConfig(nodeDef objects.Config) {

}
