package server

import (
	"github.com/adamspouele/vultron/reader"
)

// load configurations from config.yml
func LoadConfig() reader.ConfigBlueprint {
	return reader.ReadConfigBlueprintFromPath("config/config.yml")
}
