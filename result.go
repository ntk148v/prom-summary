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
	"os"
)

// PromSummary is the result format
type PromSummary struct {
	Name                       string     `json:"name" yaml:"name"`
	Address                    string     `json:"address" yaml:"address"`
	Status                     PromStatus `json:"status" yaml:"status"`
	Error                      string     `json:"error" yaml:"error"`
	StorageRetention           string     `json:"storage_retention" yaml:"storage_retention"`
	NumOfActiveTargets         string     `json:"number_of_active_targets" yaml:"number_of_active_targets"`
	NumOfDroppedTargets        string     `json:"number_of_dropped_targets" yaml:"number_of_dropped_targets"`
	NumOfTimeSeries            string     `json:"number_of_time_series" yaml:"number_of_time_series"`
	NumOfChunks                string     `json:"number_of_chunks" yaml:"number_of_chunks"`
	NumOfIngestedSamplesPerSec string     `json:"number_of_ingested_samples_per_seconds" yaml:"number_of_ingested_samples_per_seconds"`
}

// PromStatus is the state of the Prometheus endpoint, if there
// is an error, mark this Prometheus instance as NOT OK.
type PromStatus int

const (
	PromStatusOK PromStatus = iota
	PromStatusNotOK
)

func (s PromStatus) String() string {
	switch s {
	case PromStatusOK:
		return "OK"
	case PromStatusNotOK:
		return "NotOK"
	default:
		return "unknown"
	}
}

func (ps *PromSummary) setStatus(err error) {
	ps.Status = PromStatusOK
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		ps.Status = PromStatusNotOK
		ps.Error = err.Error()
	}
}
