package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	var err error
	var configFile string
	var dockerfile string
	var dryRun bool
	var name string
	var folder string

	var buildConfiguration BuildConfiguration

	// parse configuration flags from command line
	flag.StringVar(&folder, "folders", "", "Folders to include in docker")
	flag.StringVar(&name, "name", "", "Image name")
	flag.StringVar(&dockerfile, "dockerfile", "", "Full path to the dockerfile")
	flag.StringVar(&configFile, "f", "", "Path to build configuration file")
	flag.BoolVar(&dryRun, "d", false, "Don't run docker, only build and print sorted map")
	flag.Parse()

	// if configFile is specified - deserialize it
	if len(configFile) > 0 {
		// Exit if something is off
		_ = ValidatePath(configFile, true)

		// get configuration
		buildConfiguration, err = ParseFile(configFile)
		if err != nil {
			log.Fatal(err)
		}
	} else if len(dockerfile) > 0 {
		if len(name) == 0 {
			log.Fatal("-name must be specified")
		}

		var folders []string
		if len(folder) > 0 {
			folders = strings.Split(folder, ",")
			err = checkFoldersExistence(folders...)
			if err != nil {
				log.Fatal(err)
			}
		}

		// Dockerfile mode, just create config with 1 image
		buildConfiguration = BuildConfiguration{
			Images: []Image{{
				Folders:       folders,
				ContainerName: name,
				Dockerpath:    dockerfile,
				ForbidCache:   false,
			}},
		}
	} else {
		// show error & exit
		log.Fatalf("Neither configFile or Dockerfile was specified")
	}

	// build images
	if !dryRun {
		err = BuildImages(buildConfiguration)
		if err != nil {
			fmt.Printf("%v\n", err.Error())
			os.Exit(1)
		}

		// if everything is ok - exit gracefully
		fmt.Printf("Successfully built %v images\n", buildConfiguration.NumJobs())
	} else {
		executable, err := buildExecutableMap(buildConfiguration)
		if err != nil {
			fmt.Printf("%v\n", err.Error())
			os.Exit(1)
		}

		fmt.Printf("%v\n", executable)
	}

	os.Exit(0)
}
