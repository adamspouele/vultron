package reader

import (
	"fmt"
	"log"

	"gopkg.in/yaml.v2"
)

// this struct is the structure of the configuration blueprint
type ConfigBlueprint struct {
	Organization string
	Docker       struct {
		Address string
	}
	Cosmoverse struct {
		Pidpath   string
		Blueprint struct {
			Default string
		}
	}
}

// Read blueprint content from string
func ReadConfigBlueprintFromPath(path string) ConfigBlueprint {

	content := ReadFileContent(path)

	fmt.Println(content)

	cprint := ConfigBlueprint{}

	err := yaml.Unmarshal([]byte(content), &cprint)
	checkError(err)

	// if mandatory parameters are not present, throw exception
	if cprint.Docker.Address == "" {
		log.Fatalf("docker.address parameter is mandatory in the config file.")
	}

	if cprint.Cosmoverse.Pidpath == "" {
		log.Fatalf("cosmoverse.pidpath parameter is mandatory in the config file.")
	}
	if cprint.Cosmoverse.Blueprint.Default == "" {
		log.Fatalf("cosmoverse.blueprint.default parameter is mandatory in the config file.")
	}

	return cprint
}
