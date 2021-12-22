package objects

import "fmt"

type Config struct {
	ObjectBase
	Meta ObjectMetaBase `json:"meta"`
}

func (c Config) String() string {
	return fmt.Sprintf("Kind: %s, Meta: %v",
		c.Kind, c.Meta)
}

// LogFileExporterConfig is config object for logfileexporter
type LogFileExporterConfig struct {
	Config

	MaxLogFileCount uint `json:"logFileCount"` // Max logfile count per jobId (If set 0 => default => 30)
}

type JobSchedulerConfig struct {
	Config
}
