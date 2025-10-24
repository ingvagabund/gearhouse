package main

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

import (
	"fmt"
)

type RecipeSelector struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Spec       Spec   `yaml:"spec"`
}

// Spec contains the main configuration for recipes and constraints.
type Spec struct {
	Recipes     []string    `yaml:"recipes"`
	Constraints Constraints `yaml:"constraints"`
}

// Constraints holds the various gap and limitation settings.
type Constraints struct {
	MinimalMealGap     int `yaml:"minimalMealGap"`
	MinimalBeefGap     int `yaml:"minimalBeefGap"`
	MaximalNonFiberGap int `yaml:"maximalNonFiberGap"`
	MaximalNonVeganGap int `yaml:"maximalNonVeganGap"`
	MaximalNonFishGap  int `yaml:"maximalNonFishGap"`
}

func validateConfig(config *RecipeSelector) error {
	// one recipe => all the meals are the same
	// two recipes => meals alternate
	if len(config.Spec.Recipes) < 3 {
		return fmt.Errorf("At least three recipes need to be provided, got %v instead", len(config.Spec.Recipes))
	}
  return nil
}
