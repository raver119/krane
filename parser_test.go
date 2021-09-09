package main

import (
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
	"testing"
)

var buildConf = BuildConfiguration{
	Images:  []Image{{ContainerName: "Alpha", Dockerpath: "/path/to/Folder", Folders: []string{}}, {ContainerName: "Beta", Dockerpath: "/path/to/OtherFolder", Folders: []string{}}},
	Threads: 12,
}

func TestParse_Serde(t *testing.T) {
	bytes, err := yaml.Marshal(buildConf)
	require.NoError(t, err)

	var conf BuildConfiguration
	err = yaml.Unmarshal(bytes, &conf)
	require.NoError(t, err)
	require.Equal(t, buildConf, conf)
}

func TestParse_Bytes(t *testing.T) {
	bytes, err := yaml.Marshal(buildConf)
	require.NoError(t, err)

	conf, err := ParseBytes(bytes)
	require.NoError(t, err)
	require.Equal(t, buildConf, conf)
}

func TestParse_String(t *testing.T) {
	bytes, err := yaml.Marshal(buildConf)
	require.NoError(t, err)

	conf, err := ParseString(string(bytes))
	require.NoError(t, err)
	require.Equal(t, buildConf, conf)
}

func TestParse_File(t *testing.T) {
	conf, err := ParseFile("./resources/test.yaml")
	require.NoError(t, err)
	require.Equal(t, buildConf, conf)
}

func TestParse_File_With_Folders(t *testing.T) {
	conf, err := ParseFile("./resources/test_with_folders.yaml")
	require.NoError(t, err)

	require.Len(t, conf.Images, 1)
	require.Len(t, conf.Images[0].Folders, 2)
	require.Equal(t, conf.Images[0].Folders[0], "alpha:ALPHA")
	require.Equal(t, conf.Images[0].Folders[1], "beta:BETA")
}
