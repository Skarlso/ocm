package simplemapmerge

import (
	"ocm.software/ocm/api/ocm/valuemergehandler/handlers/defaultmerge"
	"ocm.software/ocm/api/ocm/valuemergehandler/hpi"
	"ocm.software/ocm/api/utils"
)

type Mode = defaultmerge.Mode

const (
	MODE_DEFAULT = defaultmerge.MODE_DEFAULT
	MODE_NONE    = defaultmerge.MODE_NONE
	MODE_LOCAL   = defaultmerge.MODE_LOCAL
	MODE_INBOUND = defaultmerge.MODE_INBOUND
)

func NewConfig(overwrite Mode, entries ...*hpi.Specification) *Config {
	return &Config{
		Config:  *defaultmerge.NewConfig(overwrite),
		Entries: utils.Optional(entries...),
	}
}

type Config struct {
	defaultmerge.Config
	Entries *hpi.Specification `json:"entries,omitempty"`
}
