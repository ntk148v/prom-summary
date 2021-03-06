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
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	promclient "github.com/prometheus/client_golang/api"
	prometheus "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
)

var headers = []string{
	"name", "address", "status", "error", "version",
	"storage retention", "number of active targets",
	"number of dropped targets", "number of time series",
	"number of chunks", "number of ingested samples per seconds",
}

func initClient(address, username, password string) (prometheus.API, error) {
	promCfg := promclient.Config{Address: address}
	if username != "" && password != "" {
		promCfg.RoundTripper = &BasicAuthTransport{
			Username: username,
			Password: password,
		}
	}
	client, err := promclient.NewClient(promCfg)
	if err != nil {
		return nil, err
	}
	api := prometheus.NewAPI(client)
	return api, nil
}

func main() {

	a := kingpin.New(filepath.Base(os.Args[0]), "A lazy tool written by Golang to export Prometheus summary.")
	var (
		cfgFile string
		cfg     *Config
		results []*PromSummary
		wg      sync.WaitGroup
	)
	a.Flag("config.file", "Prom-summary configuration file path.").
		Default("etc/config.yml").StringVar(&cfgFile)

	_, err := a.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrapf(err, "Error parsing commandline arguments"))
		a.Usage(os.Args[1:])
		os.Exit(2)
	}

	cfg, err = LoadFile(cfgFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrapf(err, "Error loading configuration file"))
		os.Exit(2)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for k, v := range cfg.PrometheusConfigs {
		wg.Add(1)
		go func(ctx context.Context, promName string, promCfg PrometheusConfig) {
			record := &PromSummary{
				Name:    promName,
				Address: promCfg.Address,
				Status:  PromStatusOK,
			}
			defer func() {
				results = append(results, record)
				wg.Done()
			}()
			promAPI, err := initClient(promCfg.Address, promCfg.BasicAuth.Username,
				promCfg.BasicAuth.Password)
			if err != nil {
				record.setStatus(errors.Wrapf(err, "Error initializing Prometheus API client"))
				return
			}
			// Get version
			buildInfo, err := promAPI.Buildinfo(ctx)
			if err != nil {
				record.setStatus(errors.Wrapf(err, "Error getting build info"))
				return
			}
			record.Version = buildInfo.Version
			// Ger number of targets
			targets, err := promAPI.Targets(context.Background())
			if err != nil {
				record.setStatus(errors.Wrapf(err, "Error getting targets"))
				return
			}
			record.NumOfActiveTargets = strconv.Itoa(len(targets.Active))
			record.NumOfDroppedTargets = strconv.Itoa(len(targets.Dropped))
			// Get storage retention
			runtimeInfo, err := promAPI.Runtimeinfo(ctx)
			if err != nil {
				record.setStatus(errors.Wrapf(err, "Error getting runtime info"))
				return
			}
			record.StorageRetention = runtimeInfo.StorageRetention
			// Get number of time series
			record.NumOfTimeSeries = strconv.Itoa(runtimeInfo.TimeSeriesCount)
			// Get number of chunks
			record.NumOfChunks = strconv.Itoa(runtimeInfo.ChunkCount)
			// Get number of ingested samples per second
			val, _, err := promAPI.Query(ctx, "rate(prometheus_tsdb_head_samples_appended_total[5m])", time.Now())
			if err != nil {
				record.setStatus(errors.Wrapf(err, "Error querying metrics"))
				return
			}
			switch v := val.(type) {
			case model.Vector:
				total := 0.0
				for _, s := range v {
					total += float64(s.Value)
				}
				record.NumOfIngestedSamplesPerSec = strconv.FormatFloat(total/float64(len(v)), 'E', -1, 64)
			default:
				record.setStatus(errors.Errorf("unsupported type: '%q'", v))
				return
			}
		}(ctx, k, v)
	}

	wg.Wait()

	// Write the result
	if cfg.OutputConfig.File != "" {
		fmt.Println("You can found the report here ", cfg.OutputConfig.File)
	}
	switch strings.ToLower(cfg.OutputConfig.Format) {
	case "table":
		writer := os.Stdout
		if cfg.OutputConfig.File != "" {
			writer, err = os.Create(cfg.OutputConfig.File)
			defer writer.Close()
			if err != nil {
				fmt.Fprintln(os.Stderr, errors.Wrapf(err, "Error printing result"))
			}
		}
		table := tablewriter.NewWriter(writer)
		table.SetHeader(headers)
		table.SetAlignment(tablewriter.ALIGN_RIGHT)
		for _, record := range results {
			table.Append([]string{
				record.Name, record.Address, record.Status.String(),
				record.Error, record.Version, record.StorageRetention,
				record.NumOfActiveTargets, record.NumOfDroppedTargets,
				record.NumOfTimeSeries, record.NumOfChunks,
				record.NumOfIngestedSamplesPerSec,
			})
		}
		table.Render()
	case "json":
		content, _ := json.MarshalIndent(results, "", "")
		if cfg.OutputConfig.File != "" {
			_ = ioutil.WriteFile(cfg.OutputConfig.File, content, 0644)
		} else {
			os.Stdout.Write(content)
		}
	case "yaml":
		content, _ := yaml.Marshal(results)
		if cfg.OutputConfig.File != "" {
			_ = ioutil.WriteFile(cfg.OutputConfig.File, content, 0644)
		} else {
			os.Stdout.Write(content)
		}
	case "csv":
		writer := os.Stdout
		if cfg.OutputConfig.File != "" {
			writer, err = os.Create(cfg.OutputConfig.File)
			defer writer.Close()
			if err != nil {
				fmt.Fprintln(os.Stderr, errors.Wrapf(err, "Error printing result"))
			}
		}
		w := csv.NewWriter(writer)
		defer w.Flush()

		w.Write(headers)
		for _, record := range results {
			w.Write([]string{
				record.Name, record.Address, record.Status.String(),
				record.Error, record.Version, record.StorageRetention,
				record.NumOfActiveTargets, record.NumOfDroppedTargets,
				record.NumOfTimeSeries, record.NumOfChunks,
				record.NumOfIngestedSamplesPerSec,
			})
		}
	}
}
