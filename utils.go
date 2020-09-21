package main

import (
	"fmt"
	"os"
	"sort"
)

func ValidatePath(path string, fatal bool) (err error) {
	if len(path) == 0 {
		fmt.Printf("Please specify build configuration file\n")
		if fatal {
			os.Exit(1)
		} else {
			return fmt.Errorf("missing the filename argument\n")
		}
	}

	// validate the file
	if d, err := os.Stat(path); err != nil {
		fmt.Printf("got error [%v] when tried to stat file [%v]\n", err.Error(), path)
		if fatal {
			os.Exit(1)
		} else {
			return err
		}
	} else {
		if d.IsDir() {
			fmt.Printf("build configuration must be a file, but got directory instead\n")
			if fatal {
				os.Exit(1)
			} else {
				return fmt.Errorf("build configuration must be a file, but got directory instead\n")
			}
		}
	}

	return
}

type sortBy func(p1, p2 *Image) bool

type imageSorter struct {
	images		[]Image
	sorter		func(p1, p2 *Image) bool
}

func (s *imageSorter) Len() int {
	return len(s.images)
}

func (s *imageSorter) Swap(i, j int) {
	s.images[i], s.images[j] = s.images[j], s.images[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *imageSorter) Less(i, j int) bool {
	return s.sorter(&s.images[i], &s.images[j])
}

func (s sortBy) Sort(images []Image) {
	is := &imageSorter{
		images: images,
		sorter: s,
	}

	sort.Sort(is)
}

func SortImages(conf *BuildConfiguration) {
	sorter := func(p1, p2 *Image) bool {
		return p1.ContainerName < p2.ContainerName
	}

	sortBy(sorter).Sort(conf.Images)
}