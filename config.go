// Copyright (c) 2021 Kien Nguyen-Tuan <kiennt2609@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config is the top-level configuration
type Config struct {
	PrometheusConfigs map[string]PrometheusConfig `yaml:"prometheus_configs"`
	OutputConfig      OutputConfig                `yaml:"output_config"`
}

// OutputConfig defines output related configurations.
type OutputConfig struct {
	// Format is output format, 'table', 'json', 'csv', 'json'
	// 'table' by default.
	Format string `yaml:"format"`
	// File is the output file path, by default, Prom-summary will
	// return output to stdout. If this field is specified,
	// the output will be written to file instead.
	File string `yaml:"file"`
}

// PrometheusConfig is the Prometheus instance config.
type PrometheusConfig struct {
	Address   string    `yaml:"address"`
	BasicAuth BasicAuth `yaml:"basic_auth"`
}

// BasicAuth stores authentication (username, password) configuration.
type BasicAuth struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

var (
	// DefaultOutputConfig is the default output configuration
	// By default, print the output to stdout/stderr with format table.
	DefaultOutputConfig = OutputConfig{
		Format: "table",
	}

	// DefaultConfig is the default top-level configuration.
	DefaultConfig = Config{
		OutputConfig: DefaultOutputConfig,
	}
)

// UnmarshalYAML implements the yaml.Unmarshaler interface
func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*c = DefaultConfig
	// We want to set c to the defaults and then overwrite it with the input.
	// To make unmarshal fill the plain data struct rather than calling UnmarshalYAML
	// again, we have to hide it using a type indirection.
	type plain Config
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}
	return nil
}

// String represents Configuration instance as string.
func (c *Config) String() string {
	b, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Sprintf("<error creating config string: %s>", err)
	}
	return string(b)
}

// Load parses the YAML input s into a Config.
func Load(s string) (*Config, error) {
	cfg := &Config{}
	err := yaml.UnmarshalStrict([]byte(s), cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// LoadFile parses the given YAML file into a Config.
func LoadFile(filename string) (*Config, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	cfg, err := Load(string(content))
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
