package main

import (
	"fmt"
	"path"
	"strings"
)

type Folder struct {
	Source string
	Target string
}

func NewFolder(str string) (f Folder, err error) {
	split := strings.Split(str, ":")
	if len(split) > 2 {
		return Folder{}, fmt.Errorf("wrong folder format: [%v]", str)
	} else if len(split) == 2 {
		return Folder{Source: split[0], Target: split[1]}, nil
	} else {
		// last path component becomes a target
		_, target := path.Split(split[0])
		return Folder{Source: split[0], Target: target}, nil
	}
}
