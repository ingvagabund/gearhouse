package main

import (
	"fmt"
	"io/ioutil"
	"maps"
	"os"
	"slices"
	"sort"

	"github.com/aquilax/cooklang-go"
	"github.com/ingvagabund/gearhouse/generators/recipes/pkg"
	"gopkg.in/yaml.v3"
	"k8s.io/klog/v2"
)

// TODO(ingvagabund):
// - validate the yaml data
// - sort the ingredients based on where they are in a shop (predefined ordering for each shop, growing ad-hoc)
// - divide the list into two parts: always buy and once in a while (e.g. sul, pepr, koreni, ...)
// - ingredient expansion: e.g. "horky zeleninovy vyvar" -> list of ingredients to make it instead
// - unit sticking: stick basic units (e.g. g, l, ...) into the quantity, while keep a longer "unit" with a blank (e.g. lzice, spetka, ...)
// - extend the recipes with a number of servings

func readRecipeFromFile(filename string) (*cooklang.Recipe, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return cooklang.ParseString(string(data))
}

func readShopCategories(filename string) (*pkg.Shop, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var shop pkg.Shop
	if err := yaml.Unmarshal(data, &shop); err != nil {
		return nil, err
	}

	return &shop, nil
}

func printIngredientsByShopIndex(ingSet *pkg.IngredientsSet, shop *pkg.Shop) {
	mapping := pkg.ConstructIngredientIndex(shop)

	// index -> []ingredient
	listing := make(map[int][]string)
	for ing := range ingSet.Ingredients {
		idx, exists := mapping[ing]
		if !exists {
			klog.Infof("ingredient %q not found in %q shop", ing, shop.Name)
			continue
		}
		listing[idx] = append(listing[idx], ing)
	}

	listingIndices := slices.Sorted(maps.Keys(listing))

	for _, idx := range listingIndices {
		ings := listing[idx]
		sort.Strings(ings)
		for _, ing := range ings {
			str := ingSet.Ingredients[ing].ToString()
			if str == "" {
				fmt.Printf("%v\n", ing)
			} else {
				fmt.Printf("%v: %v\n", ing, str)
			}
		}
	}
}

func main() {

	initFlags()
	validateFlags()

	shop, err := readShopCategories(shopFilename)
	if err != nil {
		klog.Error(err)
		os.Exit(1)
		return
	}

	ingSet := pkg.NewIngredientsSet()

	for _, recipeFile := range recipeFilenames {
		r, err := readRecipeFromFile(recipeFile)
		if err != nil {
			klog.Error(err)
			os.Exit(1)
			return
		}

		for _, step := range r.Steps {
			ingSet.MergeIngredients(step.Ingredients)
		}
	}

	ingSet.Consolidate()

	printIngredientsByShopIndex(ingSet, shop)
}
