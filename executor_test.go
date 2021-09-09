package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"
)

func Test_findDockerDependencies(t *testing.T) {

	tests := []struct {
		name       string
		dockerfile string
		wantDeps   []string
		wantErr    bool
	}{
		{"test_0", "FROM ubuntu:20.04\n#do something", []string{"ubuntu:20.04"}, false},
		{"test_1", "FROM ubuntu:20.04\n#do something\nFROM alpine:latest\n#do something else", []string{"ubuntu:20.04", "alpine:latest"}, false},
		{"test_2", "FROM ubuntu:20.04\n#do something\nFROM alpine\n#do something else", []string{"ubuntu:20.04", "alpine:latest"}, false},
		{"test_10", "some random file content", []string{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDeps, err := findDockerDependencies(tt.dockerfile)
			if (err != nil) != tt.wantErr {
				t.Errorf("findDockerDependencies() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(tt.wantDeps) > 0 && !reflect.DeepEqual(gotDeps, tt.wantDeps) {
				t.Errorf("findDockerDependencies() gotDeps = %v, want %v", gotDeps, tt.wantDeps)
			}
		})
	}
}

func Test_scanDependencies(t *testing.T) {
	tests := []struct {
		name    string
		config  BuildConfiguration
		wantExt Dependencies
		wantInt Dependencies
		wantBwd Dependencies
		wantErr bool
	}{
		{"test_0", BuildConfiguration{
			Images:  []Image{{ContainerName: "image1", Dockerpath: "./resources/setup_nodeps/Image1"}, {ContainerName: "image2", Dockerpath: "./resources/setup_nodeps/Image2"}, {ContainerName: "image3", Dockerpath: "./resources/setup_nodeps/Image3"}},
			Threads: 0,
		}, Dependencies{"image1:latest": []string{"ubuntu:20.04"}, "image2:latest": []string{"ubuntu:latest", "alpine:latest"}, "image3:latest": []string{"nginx:latest"}},
			Dependencies{"image1:latest": []string{}, "image2:latest": []string{}, "image3:latest": []string{}},
			Dependencies{"image1:latest": []string{}, "image2:latest": []string{}, "image3:latest": []string{}}, false},

		{"test_1", BuildConfiguration{
			Images:  []Image{{ContainerName: "image1", Dockerpath: "./resources/setup_onedep/Image1"}, {ContainerName: "image2", Dockerpath: "./resources/setup_onedep/Image2"}, {ContainerName: "image3", Dockerpath: "./resources/setup_onedep/Image3"}},
			Threads: 0,
		}, Dependencies{"image1:latest": []string{"ubuntu:20.04"}, "image2:latest": []string{"ubuntu:latest"}, "image3:latest": []string{"nginx:latest"}},
			Dependencies{"image1:latest": []string{}, "image2:latest": []string{"image1:latest"}, "image3:latest": []string{}},
			Dependencies{"image1:latest": []string{"image2:latest"}, "image2:latest": []string{}, "image3:latest": []string{}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotExt, gotInt, gotBwd, err := scanDependencies(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("scanDependencies() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotExt, tt.wantExt) {
				t.Errorf("scanDependencies() gotExt = %v, want %v", gotExt, tt.wantExt)
			}
			if !reflect.DeepEqual(gotInt, tt.wantInt) {
				t.Errorf("scanDependencies() gotInt = %v, want %v", gotInt, tt.wantInt)
			}
			if !reflect.DeepEqual(gotBwd, tt.wantBwd) {
				t.Errorf("scanDependencies() gotBwd = %v, want %v", gotBwd, tt.wantBwd)
			}
		})
	}
}

func Test_buildExecutableMap(t *testing.T) {

	configNoDeps := BuildConfiguration{
		Images:  []Image{{ContainerName: "image1", Dockerpath: "./resources/setup_nodeps/Image1"}, {ContainerName: "image2", Dockerpath: "./resources/setup_nodeps/Image2"}, {ContainerName: "image3", Dockerpath: "./resources/setup_nodeps/Image3"}},
		Threads: 0,
	}

	configSingleDep := BuildConfiguration{
		Images:  []Image{{ContainerName: "image1", Dockerpath: "./resources/setup_onedep/Image1"}, {ContainerName: "image2", Dockerpath: "./resources/setup_onedep/Image2"}, {ContainerName: "image3", Dockerpath: "./resources/setup_onedep/Image3"}},
		Threads: 0,
	}

	configSingleRoot := BuildConfiguration{
		Images:  []Image{{ContainerName: "image1", Dockerpath: "./resources/setup_oneroot/Image1"}, {ContainerName: "image2", Dockerpath: "./resources/setup_oneroot/Image2"}, {ContainerName: "image3", Dockerpath: "./resources/setup_oneroot/Image3"}},
		Threads: 0,
	}

	tests := []struct {
		name       string
		config     BuildConfiguration
		wantResult ExecutableMap
		wantErr    bool
	}{
		{"test_0", configNoDeps, ExecutableMap{0: configNoDeps.Images}, false},
		{"test_1", configSingleDep, ExecutableMap{0: []Image{{ContainerName: "image1", Dockerpath: "./resources/setup_onedep/Image1"}, {ContainerName: "image3", Dockerpath: "./resources/setup_onedep/Image3"}},
			1: []Image{{ContainerName: "image2", Dockerpath: "./resources/setup_onedep/Image2"}}}, false},
		{"test_2", configSingleRoot, ExecutableMap{1: []Image{{ContainerName: "image1", Dockerpath: "./resources/setup_oneroot/Image1"}, {ContainerName: "image3", Dockerpath: "./resources/setup_oneroot/Image3"}},
			0: []Image{{ContainerName: "image2", Dockerpath: "./resources/setup_oneroot/Image2"}}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := buildExecutableMap(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildExecutableMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("buildExecutableMap() gotResult = \n%v\nvs\n%v", gotResult, tt.wantResult)
			}
		})
	}
}

func Test_prepareFolders_1(t *testing.T) {
	target, err := os.MkdirTemp(os.TempDir(), "*-target")
	require.NoError(t, err)

	srcA, err := os.MkdirTemp(os.TempDir(), "*-src")
	require.NoError(t, err)

	srcB, err := os.MkdirTemp(os.TempDir(), "*-src")
	require.NoError(t, err)

	require.NoError(t, prepareFolders(target, srcA, srcB))

	f, err := ioutil.ReadDir(target)
	require.NoError(t, err)
	require.Len(t, f, 2)

	_, sA := path.Split(srcA)
	_, sB := path.Split(srcB)

	for _, v := range f {
		assert.True(t, v.Name() == sA || v.Name() == sB, "failed name: %v", v.Name())
	}
}

func Test_checkFoldersExistence(t *testing.T) {
	tests := []struct {
		name    string
		folders []string
		wantErr bool
	}{
		{"test_0", []string{os.TempDir()}, false},
		{"test_1", []string{os.TempDir(), os.TempDir()}, false},
		{"test_2", []string{os.TempDir(), path.Join(os.TempDir(), "some non-existent folder")}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkFoldersExistence(tt.folders...); (err != nil) != tt.wantErr {
				t.Errorf("checkFoldersExistence() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
