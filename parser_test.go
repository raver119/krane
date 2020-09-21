package main

import (
	"gopkg.in/yaml.v2"
	"reflect"
	"testing"
)

var buildConf = BuildConfiguration{
	Images:  []Image{{ContainerName: "Alpha", Dockerpath: "/path/to/Folder"}, {ContainerName: "Beta", Dockerpath: "/path/to/OtherFolder"}},
	Threads: 12,
}

func TestParse_Serde(t *testing.T) {
	bytes, err := yaml.Marshal(buildConf)
	if err != nil {
		t.Fatalf(err.Error())
	}

	var conf BuildConfiguration
	err = yaml.Unmarshal(bytes, &conf)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if !reflect.DeepEqual(buildConf, conf) {
		t.Fatalf("objects are not equal after ser/de")
	}
}

func TestParse_Bytes(t *testing.T) {
	bytes, err := yaml.Marshal(buildConf)
	if err != nil {
		t.Fatalf(err.Error())
	}

	conf, err := ParseBytes(bytes)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if !reflect.DeepEqual(buildConf, conf) {
		t.Fatalf("objects are not equal after ser/de")
	}
}

func TestParse_String(t *testing.T) {
	bytes, err := yaml.Marshal(buildConf)
	if err != nil {
		t.Fatalf(err.Error())
	}

	conf, err := ParseString(string(bytes))
	if err != nil {
		t.Fatalf(err.Error())
	}

	if !reflect.DeepEqual(buildConf, conf) {
		t.Fatalf("objects are not equal after ser/de")
	}
}

func TestParse_File(t *testing.T) {
	conf, err := ParseFile("./resources/test.yaml")
	if err != nil {
		t.Fatalf(err.Error())
	}

	if !reflect.DeepEqual(buildConf, conf) {
		t.Fatalf("objects are not equal after ser/de")
	}
}