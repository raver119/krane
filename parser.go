package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

type NamesMap map[string]Image

/*
	This method returns specified Dockerfile as string
*/
func (i Image) Dockerfile() (content string, err error) {
	// FIXME: remove double // here
	bytes, err := ioutil.ReadFile(i.Dockerpath + "/Dockerfile")
	if err == nil {
		content = string(bytes)
	}

	return
}

/*
	This method returns number of images to be built
*/
func (bc BuildConfiguration) NumJobs() int {
	return len(bc.Images)
}

/*
	This method returns container names organized into map
*/
func (bc BuildConfiguration) NamesMap() (NamesMap, error) {
	result := make(NamesMap)

	for _, v := range bc.Images {
		if _, has := result[v.ContainerName]; has {
			return result, fmt.Errorf("image [%v] is declared more than once", v.ContainerName)
		}

		if strings.Contains(v.ContainerName, ":") {
			result[v.ContainerName] = v
		} else {
			result[v.ContainerName+":latest"] = v
		}
	}

	return result, nil
}

/*main
This method returns slice of image names
*/
func (bc BuildConfiguration) Names() (result []string) {
	SortImages(&bc)
	for _, v := range bc.Images {

		if strings.Contains(v.ContainerName, ":") {
			result = append(result, v.ContainerName)
		} else {
			result = append(result, v.ContainerName+":latest")
		}
	}

	return
}

/*
	This function provides YAML deserialization of given byte slice
*/
func ParseBytes(conf []byte) (bc BuildConfiguration, err error) {
	err = yaml.Unmarshal(conf, &bc)
	if err == nil {
		SortImages(&bc)
		for i, _ := range bc.Images {
			if len(bc.Images[i].Folders) == 0 {
				bc.Images[i].Folders = []string{}
			}
		}
	}
	return
}

/*
	This function provides deserialization of a given string
*/
func ParseString(conf string) (BuildConfiguration, error) {
	return ParseBytes([]byte(conf))
}

/*
	This function provides deserialization of a given YAML file
*/
func ParseFile(fileName string) (bc BuildConfiguration, err error) {
	conf, err := ioutil.ReadFile(fileName)
	if err != nil {
		return
	}

	return ParseBytes(conf)
}
