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
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/aquilax/cooklang-go"
	"k8s.io/klog/v2"
)

// TODO(ingvagabund):
// - take into account min ingredient size. E.g. couple two recipes that needs a chicken (normally a chicken is sold in 0.5kg packages, yet each recipe consumes only 250g).
// - new constraint:
//   - min number of meals before repeating the same recipe
//   - max number of meals before repeating the same recipe
//   - going vegetarian
//   - beef meat only once in a while
//   - chicken more often (max items before the repetition)
//   - at least that much fibers per N items

const (
	minimalMealGap     int = 3
	minimalBeefGap     int = 2
	maximalNonFiberGap int = 2
	maximalNonVeganGap int = 2

	dietaryOptionKey        = "dietary option"
	dietaryOptionVegetarian = "vegetarian"
)

func readRecipeFromFile(filename string) (*cooklang.Recipe, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return cooklang.ParseString(string(data))
}

type meal struct {
	date   time.Time
	recipe *cooklang.Recipe
}

type meals []meal

func (a meals) Len() int { return len(a) }

// (the most recent first)
func (a meals) Less(i, j int) bool { return a[i].date.After(a[j].date) }
func (a meals) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func readHistoryFromFile(filename string) ([]meal, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	mealsList := []meal{}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Each line is expected to have three parts:
		// [date] [type] value
		// In here type is expected to be set to "meal"
		re := regexp.MustCompile(`^\[(.*?)\]\s+\[(.*?)\](\s+(.*))*$`)
		matches := re.FindStringSubmatch(line)
		if len(matches) < 5 {
			return nil, fmt.Errorf("unable to parse line %q", line)
		}

		if matches[2] != "meal" {
			continue
		}

		if matches[4] == "" {
			return nil, fmt.Errorf("line %q is missing a recipe", line)
		}

		text := strings.Replace(matches[4], "\\n", "\n", -1)

		recipe, err := cooklang.ParseString(text)
		if err != nil {
			return nil, fmt.Errorf("line %q: error parsing the recipe: %v", line, err)
		}

		// const layout = "2025-01-02 10:15"
		t, err := time.Parse("2006-01-02 15:04", matches[1])
		if err != nil {
			return nil, fmt.Errorf("line %q: unable to parse date %q: %v", line, matches[1], err)
		}

		mealsList = append(mealsList, meal{
			date:   t,
			recipe: recipe,
		})
	}

	sort.Sort(meals(mealsList))

	return mealsList, nil
}

func main() {

	initFlags()
	validateFlags()

	// Current contraints:
	// - each meal can not repeat in two consecutive segments
	// - each meal can not repeat in less than 4 distinct meals

	recipes := []*cooklang.Recipe{}
	// Read all recipes
	for _, recipeFile := range recipeFilenames {
		r, err := readRecipeFromFile(recipeFile)
		if err != nil {
			klog.Error(err)
			os.Exit(1)
			return
		}

		recipes = append(recipes, r)
	}

	keys := map[string]struct{}{}
	for _, recipe := range recipes {
		// TODO(ingvagabund): filter out by category (only lunch)
		title := recipe.Metadata["title"].(string)
		// fmt.Printf("recipe: %v\n", title)
		if _, exists := keys[title]; exists {
			klog.Error(fmt.Errorf("recipe %q is duplicated in the list of recipes", title))
			os.Exit(1)
			return
		}
		keys[title] = struct{}{}
	}

	// Read history of prepared meals
	meals, err := readHistoryFromFile(historyFilename)
	if err != nil {
		klog.Error(err)
		os.Exit(1)
		return
	}

	// titles are the keys
	for _, meal := range meals {
		title := meal.recipe.Metadata["title"].(string)
		if _, exists := keys[title]; !exists {
			klog.Error(fmt.Errorf("recipe %q not found in the list of recipes", title))
			os.Exit(1)
			return
		}
		// fmt.Printf("[%v] #%v#%v\n", meal.date, meal.recipe.Metadata["title"], meal.recipe.Steps)
	}

	initDfsExploration(meals, recipes).startExploration()
}

type dfsExploration struct {
	mealNames  []string
	mealsList  []meal
	recipesMap map[string]*cooklang.Recipe

	validPaths    [][]string
	validPathsMap map[string]struct{}
}

func initDfsExploration(meals []meal, recipes []*cooklang.Recipe) *dfsExploration {
	// generate already existing meals
	mealNames := []string{}
	mealsList := []meal{}
	for _, meal := range meals {
		mealNames = append(mealNames, meal.recipe.Metadata["title"].(string))
		mealsList = append(mealsList, meal)
	}
	recipesMap := map[string]*cooklang.Recipe{}
	for _, recipe := range recipes {
		recipesMap[recipe.Metadata["title"].(string)] = recipe
	}
	return &dfsExploration{
		mealNames:     mealNames,
		mealsList:     mealsList,
		recipesMap:    recipesMap,
		validPaths:    make([][]string, 0),
		validPathsMap: make(map[string]struct{}),
	}
}

func (dfs *dfsExploration) startExploration() {
	// for idx, meal := range dfs.mealNames {
	// 	fmt.Printf("[%v] %v\n", idx, meal)
	// }
	// generate the first step
	dfs.processNode(0, []string{})
	// process the valid paths
	for _, path := range dfs.validPaths {
		fmt.Printf("VP: %#v [%v]\n", path, dfs.validateFull(path))
	}
	fmt.Printf("Valid paths in total: %v\n", len(dfs.validPaths))
}

func mealsList2Key(path []string) string {
	return strings.Join(path, "#")
}

func (dfs *dfsExploration) processNode(level uint, path []string) {
	// End the search and store the list
	if level == segmentLength+3 {
		newPath := path[:segmentLength]
		key := mealsList2Key(newPath)
		if _, exists := dfs.validPathsMap[key]; !exists {
			if !dfs.validateFull(newPath) {
				return
			}
			dfs.validPaths = append(dfs.validPaths, append([]string{}, newPath...))
			dfs.validPathsMap[key] = struct{}{}
		}
		return
	}
	// generate the next step
	for recipe := range dfs.recipesMap {
		newPath := append(path, dfs.recipesMap[recipe].Metadata["title"].(string))
		// Validate
		// fmt.Printf("%v: %#v\n", strings.Repeat("\t", int(level)), newPath)
		if !dfs.validate(newPath) {
			continue
		}
		dfs.processNode(level+1, newPath)
	}
}

func (dfs *dfsExploration) mealAtIdx(idx int, path []string) string {
	if idx < 0 {
		// explore the historical list of meals
		j := -idx - 1
		// fmt.Printf("j=%v %v\n", j, dfs.mealNames[j])
		return dfs.mealNames[j]
	}
	// fmt.Printf("i=%v %v\n", idx, path[idx])
	return path[idx]
}

func (dfs *dfsExploration) recipeAtIdx(idx int, path []string) *cooklang.Recipe {
	if idx < 0 {
		// explore the historical list of meals
		j := -idx - 1
		// fmt.Printf("j=%v %v\n", j, dfs.mealsList[j].recipe)
		return dfs.mealsList[j].recipe
	}
	// fmt.Printf("i=%v %v\n", idx, dfs.recipesMap[path[idx]])
	return dfs.recipesMap[path[idx]]
}

func hasBeef(recipe *cooklang.Recipe) bool {
	return recipe.Metadata[dietaryOptionKey] == "beef"
}

func hasFiber(recipe *cooklang.Recipe) bool {
	for _, step := range recipe.Steps {
		for _, ingr := range step.Ingredients {
			// fmt.Printf("ig: %v\n", ingr.Name)
			// TODO(ingvagabund): make this configurable
			if ingr.Name == "mlete hovezi maso" {
				return true
			}
		}
	}
	return false
}

func (dfs *dfsExploration) validateMinimalMealGap(path []string) bool {
	hIndex := len(path) - 1
	for i := hIndex - 1; i >= hIndex-minimalMealGap; i-- {
		value := dfs.mealAtIdx(i, path)
		if path[hIndex] == value {
			return false
		}
	}
	return true
}

func (dfs *dfsExploration) validate(path []string) bool {
	hIndex := len(path) - 1

	if !dfs.validateMinimalMealGap(path) {
		return false
	}

	// beef meat only once in a while
	firstBeefIdx := hIndex + 1
	for i := hIndex; i >= hIndex-minimalBeefGap; i-- {
		value := dfs.recipeAtIdx(i, path)
		// fmt.Printf("[%v]: %v\n", value.Metadata["title"].(string), hasBeef(value))
		if hasBeef(value) {
			// the first occurence
			if firstBeefIdx == hIndex+1 {
				firstBeefIdx = i
				// the second occurence
			} else {
				idxDiff := firstBeefIdx - i
				// fmt.Printf("idxDiff: %v\n", idxDiff)
				if idxDiff <= minimalBeefGap {
					return false
				}
			}
		}
	}

	return true
}

func (dfs *dfsExploration) validateFull(path []string) bool {
	// check the history for any occurence
	veganIdx := -1
	for i := 0; i < maximalNonVeganGap && i < len(dfs.mealsList); i++ {
		// fmt.Printf("do[%v]: %v\n", dfs.mealsList[i].recipe.Metadata["title"], dfs.mealsList[i].recipe.Metadata[dietaryOptionKey])
		if dfs.mealsList[i].recipe.Metadata[dietaryOptionKey] == dietaryOptionVegetarian {
			veganIdx = i
			break
		}
	}
	// fmt.Printf("vegetarian at %v\n", veganIdx)
	if veganIdx == -1 {
		// if not found the first suggested meal has to be vegetarian
		if len(path) == 0 {
			return false
		}
		if dfs.recipesMap[path[0]].Metadata[dietaryOptionKey] != dietaryOptionVegetarian {
			return false
		}
	}
	// there's at least one vegetarian meal

	hIndex := len(path) - 1
	lastHitIdx := len(path)
	// meal with fibers not repeating at most maximalNonVeganGap times
	for i := hIndex; i >= -maximalNonVeganGap; i-- {
		value := dfs.recipeAtIdx(i, path)
		// fmt.Printf("[%v]: %v\n", value.Metadata["title"].(string), value.Metadata[dietaryOptionKey].(string))
		if value.Metadata[dietaryOptionKey].(string) == dietaryOptionVegetarian {
			// the first occurrence
			if lastHitIdx == hIndex+1 {
				idxDiff := hIndex - i
				// fmt.Printf("F(%v, %v) idxDiff: %v, %v\n", lastHitIdx, i, idxDiff, maximalNonVeganGap)
				if idxDiff > maximalNonVeganGap {
					return false
				}
			}
			lastHitIdx = i
		}
		if lastHitIdx != hIndex+1 {
			idxDiff := lastHitIdx - i
			// fmt.Printf("O(%v, %v) idxDiff: %v, %v\n", lastHitIdx, i, idxDiff, maximalNonVeganGap)
			if idxDiff > maximalNonVeganGap {
				return false
			}
		}
	}
	// fmt.Printf("lastHitIdx: %v, idx: %v\n", lastHitIdx, -maximalNonVeganGap)
	idxDiff := lastHitIdx + maximalNonVeganGap
	if idxDiff > maximalNonVeganGap {
		return false
	}
	return true
}
