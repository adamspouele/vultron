package reader

import (
	"fmt"
	"log"

	"gopkg.in/yaml.v2"
)

// Blueprint is all instructions to build an environment
type Blueprint struct {
	Version  int16 // mandatory
	Services []Service
}

// Service is the atomic unit in an environment, it contain containers and volumes
type Service struct {
	Name       string
	Containers []Container
}

// Container is a docker container representing an application
type Container struct {
	Name  string
	Image string // mandatory
	Ports []Port
}

// Port represent the ports of the host and of the container which are binded
type Port struct {
	Extern int16
	Intern int16
}

func checkError(e error) {
	if e != nil {
		log.Fatalf("error: %v", e)
	}
}

// ReadBlueprintContent Read blueprint content from string
func ReadBlueprintContent(data string) Blueprint {
	bprint := Blueprint{}

	err := yaml.Unmarshal([]byte(data), &bprint)
	checkError(err)
	return bprint
}

// ReadBlueprintFromPath Read blueprint content from path
func ReadBlueprintFromPath(path string) Blueprint {

	content := ReadFileContent(path)

	fmt.Println(content)

	bprint := Blueprint{}

	err := yaml.Unmarshal([]byte(content), &bprint)
	checkError(err)

	checkBlueprintStructure(bprint)

	fmt.Println(bprint.Services[0].Containers[0].Ports)

	return bprint
}

// Check the structure is well writed
func checkBlueprintStructure(bprint Blueprint) {
	// if Version parameter is not present, throw exception
	if bprint.Version == 0 {
		log.Fatalf("version parameter is mandatory in the blueprint.")
	}

	// check that Services are well writed
	for _, serviceItem := range bprint.Services {
		for _, containerItem := range serviceItem.Containers {
			if containerItem.Image == "" {
				log.Fatalf("image parameter is mandatory on a container.")
			}
			// element is the element from someSlice for where we are
		}
	}

}
