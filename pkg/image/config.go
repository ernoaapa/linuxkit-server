package image

import (
	"github.com/c2h5oh/datasize"
	yaml "gopkg.in/yaml.v2"
)

// Config includes extra fields what linuxkit-server supports in the Linuxkit config
type Config struct {
	Image Image
}

// Image output configuration
type Image struct {
	Partitions []Partition
}

// Partition contains information about disk partition
type Partition struct {
	Boot   bool   `yaml:"boot"`
	Start  uint64 `yaml:"start"`
	Size   uint64 `yaml:"size"`
	FsType string `yaml:"type"`
}

// NewConfig parses a config file
func NewConfig(raw []byte) (*Config, error) {
	var config = &Config{}

	err := yaml.Unmarshal(raw, config)
	if err != nil {
		return config, err
	}

	if len(config.Image.Partitions) == 0 {
		config.Image.Partitions = []Partition{
			{Start: startOffset, Size: uint64(datasize.MB * 100), Boot: true, FsType: "fat32"},
		}
	}

	return config, nil
}
