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
	recipeFilenames     []string
	shopFilename        string
	printAllIngredients bool
)

func initFlags() {
	pflag.StringSliceVar(&recipeFilenames, "recipe", recipeFilenames, "List of recipes")
	pflag.StringVar(&shopFilename, "shop", shopFilename, "Shop to visit")
	pflag.BoolVar(&printAllIngredients, "printAllIngredients", printAllIngredients, "Print all ingredients. E.g. water")
	pflag.Parse()
}

func validateFlags() {
	if len(recipeFilenames) == 0 {
		klog.Error("At least one recipe needs to be provided")
		os.Exit(1)
		return
	}

	if len(shopFilename) == 0 {
		klog.Error("A shop to visit needs to be provided")
		os.Exit(1)
		return
	}
}
