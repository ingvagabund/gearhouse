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
	"strings"

	"github.com/spf13/pflag"
	"k8s.io/klog/v2"
)

var (
	repositories []string
)

func initFlags() {
	pflag.StringSliceVar(&repositories, "repository", repositories, "List of repositories to process")
	pflag.Parse()
}

func validateFlags() {
	if len(repositories) == 0 {
		klog.Error("repository is required")
		os.Exit(1)
		return
	}

	for _, repository := range repositories {
		items := strings.Split(repository, "/")
		if len(items) != 2 {
			klog.Error("repository %q is not in a 'organization/repository' form")
			os.Exit(1)
			return
		}
	}
}
