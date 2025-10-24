/*
Copyright 2025 The Gearhouse Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"os"

	"github.com/spf13/pflag"
	"k8s.io/klog/v2"
)

var (
	segmentLength   uint
	historyFilename string
	config          string
)

func initFlags() {
	pflag.UintVar(&segmentLength, "segment-length", 3, "Length of a generated segment. E.g. a working week can consist of 3 meal selections")
	pflag.StringVar(&historyFilename, "history-filename", historyFilename, "Path to a file that contains already entered meals")
	pflag.StringVar(&config, "config", historyFilename, "Generator configuration")
	pflag.Parse()
}

func validateFlags() {
	if segmentLength == 0 {
		klog.Error("The --segment-length must be positive")
		os.Exit(1)
		return
	}

	if len(historyFilename) == 0 {
		klog.Error("The --history-filename flag must be set")
		os.Exit(1)
		return
	}
}
