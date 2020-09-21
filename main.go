package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var configFile string
	var dryRun bool

	// parse configuration flags from command line
	flag.StringVar(&configFile, "f", "", "Path to build configuration file")
	flag.BoolVar(&dryRun, "d", false, "Don't run docker, only build and print sorted map")
	flag.Parse()

	// Exit if something is off
	_ = ValidatePath(configFile, true)

	// get configuration
	buildConfiguration, err := ParseFile(configFile)
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		os.Exit(1)
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
		os.Exit(0)
	} else {
		executable, err := buildExecutableMap(buildConfiguration)
		if err != nil {
			fmt.Printf("%v\n", err.Error())
			os.Exit(1)
		}

		fmt.Printf("%v\n", executable)
		os.Exit(0)
	}


}
