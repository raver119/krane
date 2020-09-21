package main

import (
	"golang.org/x/exp/errors/fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

type Dependencies map[string][]string

type Report struct {
	ContainerName	string
	Log				string
	Error			error
	Success			bool
}

type ExecutableMap map[int][]Image

/*
	This function scans Dockerfile, given as string with commands, and extracts image names it depends
 */
func findDockerDependencies(dockerfile string) (deps []string, err error) {
	re := regexp.MustCompile(`(?m)FROM (.*?[|$| |\n])`)
	substrings := re.FindAllStringSubmatch(dockerfile, -1)
	for _, v := range substrings {
		for i, dep := range v {
			// skip first match, since it's full match
			if i == 0 {
				continue
			}

			dep = strings.TrimSpace(dep)
			if !strings.Contains(dep, ":") {
				// if no tag given, assume we're on the latest tag then
				dep += ":latest"
			}

			deps = append(deps, dep)
		}
	}

	if len(deps) == 0 {
		err = fmt.Errorf("no docker dependencies found. wrong Dockerfile was passed in?")
	}

	return
}

func scanDependencies(config BuildConfiguration) (ext, int, bwd Dependencies, err error) {
	// create empty maps first
	ext = make(Dependencies)
	bwd = make(Dependencies)
	int = make(Dependencies)

	// build map of image names first
	namesMap, err := config.NamesMap()
	if err != nil {
		return
	}

	// fill backward deps map at least
	for k, _ := range namesMap {
		bwd[k] = []string{}
		ext[k] = []string{}
		int[k] = []string{}
	}

	// for each image build dependencies
	for k, v := range namesMap {
		dockerfile, err := v.Dockerfile()
		if err != nil {
			return ext, int, bwd, err
		}

		deps, err := findDockerDependencies(dockerfile)
		if err != nil {
			return ext, int, bwd, err
		}

		// store forward deps as either external or internal dependency
		// and update backward deps
		for _, v := range deps {
			// store internal backward dependency
			if _, has := bwd[v]; has {
				bwd[v] = append(bwd[v], k)
			}

			// store internal forward dependency
			if _, has := namesMap[v]; has {
				int[k] = append(int[k], v)
			} else {
				ext[k] = append(ext[k], v)
			}
		}
	}

	return
}

/*
	This function scans mapped map in search of highest integer value from deps
 */
func findDeepestLayer(mapped map[string]int, deps []string) int {
	max := 0
	for _, v := range deps {
		if layer, has := mapped[v]; has && layer > max {
			max = layer
		}
	}

	return max
}

/*
	This function builds topologically sorted graph of images, and returns it as map
 */
func buildExecutableMap(config BuildConfiguration) (result ExecutableMap, err error) {
	namesMap, _ := config.NamesMap()
	names := config.Names()

	_, inDeps, _, err := scanDependencies(config)
	if err != nil {
		return
	}

	result = make(ExecutableMap)

	currentLayer := 0
	result[currentLayer] = []Image{}

	// this structure stores information about mapping
	mapped := make(map[string]int)

	// roll through map, probably multiple times
	// first time we add all images without deps to the root image
	for _, k := range names {
		 deps, _ := inDeps[k]

		 if len(deps) == 0 {
		 	v, _ := namesMap[k]
		 	result[currentLayer] = append(result[currentLayer], v)
		 	mapped[k] = currentLayer
		 }
	}

	// shortcut: if all results were mapped to the first layer - just return then
	if len(result[currentLayer]) == config.NumJobs() {
		return
	}

	currentLayer++

	// now roll n^2 times through remaining elements
	for i := 0; i < config.NumJobs() - len(result[0]); i++ {
		for k, v := range namesMap {

			// skip this layer if it was already mapped
			if _, wasMapped := mapped[k]; wasMapped {
				continue
			}

			// now, check if all deps were already mapped
			// if yes - map this layer
			// if no - skip this layer
			deps, _ := inDeps[k]
			allMapped := true
			for _, dep := range deps {
				_, wasMapped := mapped[dep]
				allMapped = allMapped && wasMapped
			}

			// if all layers were mapped we should see, which layer they were mapped to, and map to the next one
			if allMapped {
				deepestLayer := findDeepestLayer(mapped, deps)
				result[deepestLayer+1] = append(result[deepestLayer+1], v)
				mapped[k] = deepestLayer + 1
			}
		}
	}

	// if these two values are not equal - it means I've failed to sort the graph of dependencies. panic time
	if len(mapped) != config.NumJobs() {
		err = fmt.Errorf("wasn't able to sort the graph")
		return
	}

	return
}

/*
	This function builds Docker images
 */
func BuildImages(config BuildConfiguration) (err error) {
	// make sure we use some threads
	if config.Threads < 1 {
		config.Threads = runtime.NumCPU()
	}

	executableMap, err := buildExecutableMap(config)
	if err != nil {
		return
	}

	// reports queue. so we'll know the outcome of every build
	requeue := make(chan Report, config.NumJobs())

	// now, let's build some workers which will do the actual job
	workers := make(map[int]chan Image)
	for i := 0; i < runtime.NumCPU(); i++ {
		workers[i] = make(chan Image)
		go worker(workers[i], requeue)
	}

	// storage for the reports
	var failed []Report
	var succeed []Report

	// dispatch all jobs one by one
	jobsCounter := 0
	for i := 0; i < len(executableMap); i++ {
		dispatched := 0

		// each layer is an array of images
		layer, _ := executableMap[i]
		for _, image := range layer {
			workers[jobsCounter % config.Threads] <- image
			jobsCounter += 1
			dispatched++
		}

		// now, when all jobs on this layer were dispatched - wait for them to finish
		for i := 0; i < dispatched; i++ {
			report := <- requeue
			if !report.Success {
				failed = append(failed, report)
			} else {
				succeed = append(succeed, report)
			}
		}

		// do something better here?
		if len(failed) > 0 {
			return fmt.Errorf("At least %v out of %v jobs failed", len(failed), len(config.Images))
		}
	}

	// looks like we're all good
	return
}

/*
	This function runs in an endless loop, building all images that come from input channel
 */
func worker(input chan Image, output chan<- Report) {
	for true {
		image := <- input
		builder(image, output)
	}
}

/*
	This function executes docker build
 */
func builder(image Image, reporting chan<- Report) {
	fmt.Printf("docker build -t %v %v\n", image.ContainerName, image.Dockerpath)
	cmd := exec.Command("docker", "build", "-t", image.ContainerName, image.Dockerpath)
	output, err := cmd.Output()
	if err != nil {
		reporting <- Report{ContainerName: image.ContainerName, Log: string(output), Error: err, Success: false}
	} else {
		reporting <- Report{ContainerName: image.ContainerName, Log: string(output), Error: err, Success: true}
	}

	return
}