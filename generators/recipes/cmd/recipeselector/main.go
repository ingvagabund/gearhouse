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
	"io/ioutil"
	"os"

	"github.com/aquilax/cooklang-go"
	"k8s.io/klog/v2"
)

// TODO(ingvagabund):
// - take into account min ingredient size. E.g. couple two recipes that needs a chicken (normally a chicken is sold in 0.5kg packages, yet each recipe consumes only 250g).
// - new constraint: define the min number of meals from repeating the same recipe

func readRecipeFromFile(filename string) (*cooklang.Recipe, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return cooklang.ParseString(string(data))
}

func main() {

	initFlags()
	validateFlags()

	// Current contraints:
	// - each meal can not repeat in two consecutive segments
	// - each meal can not repeat in less than 4 distinct meals

	// Read all recipes

	// Read history of prepared meals

	// Generate recipes selection for the next segment

	shop, err := readShopCategories(shopFilename)
	if err != nil {
		klog.Error(err)
		os.Exit(1)
		return
	}
}
