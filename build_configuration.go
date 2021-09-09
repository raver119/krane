package main

type BuildConfiguration struct {
	Images  []Image `yaml:"build"`
	Threads int     `yaml:"threads"`
}
