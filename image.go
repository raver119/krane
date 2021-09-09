package main

type Image struct {
	ContainerName string `yaml:"containerName"`
	Dockerpath    string `yaml:"dockerpath"`
	ForbidCache   bool   `yaml:"noCache"`
}
