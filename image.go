package main

type Image struct {
	Folders       []string `yaml:"folders"`
	ContainerName string   `yaml:"containerName"`
	Dockerpath    string   `yaml:"dockerpath"`
	ForbidCache   bool     `yaml:"noCache"`
}
