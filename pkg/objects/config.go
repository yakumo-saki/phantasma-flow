package objects

import "fmt"

const KIND_PHFLOW_CFG = "phantasma-flow-config"
const KIND_LOGFILE_EXPORTER_CFG = "logfileexporter-config"

// Config type is base struct of All config structs
// Config structs are serialized to yaml
type Config struct {
	ObjectBase `yaml:",inline"`
	Meta       ObjectMetaBase `yaml:"meta"`
}

func (c Config) String() string {
	return fmt.Sprintf("Kind: %s, Meta: %v",
		c.Kind, c.Meta)
}

// PhantasmaFlowConfig is config object for Phantasma-Flow.
//  kind is KIND_PHFLOW_CFG
type PhantasmaFlowConfig struct {
	Config `yaml:",inline"`

	Security struct {
		Mode string // none, (not impl -> apikey)
	}
}

// LogFileExporterConfig is config object for logfileexporter.
//  kind is KIND_LOGFILE_EXPORTER_CFG
type LogFileExporterConfig struct {
	Config `yaml:",inline"`

	MaxLogFileCount int `yaml:"logFileCount"` // Max logfile count per jobId (If set 0 => default => 30)
}

// GeneralConfig is config object for fallback
//  Use of this is not recommended.
type GeneralConfig struct {
	Config `yaml:",inline"`

	ConfigMap map[string]interface{}
}
