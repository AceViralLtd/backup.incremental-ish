package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// Load config from file
func Load(path string, target interface{}) error {
	fh, err := os.Open(path)
	if nil != err {
		return err
	}

	data, err := ioutil.ReadAll(fh)
	if nil != err {
		return err
	}

	return yaml.Unmarshal(data, target)
}

// Save config to file
func Save(path string, target interface{}) error {
	data, err := yaml.Marshal(target)

	if nil != err {
		return err
	}

	return ioutil.WriteFile(path, data, 0744)
}
