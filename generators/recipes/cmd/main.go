package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"maps"
	"os"
	"slices"
	"sort"

	"github.com/aquilax/cooklang-go"
	"github.com/ingvagabund/gearhouse/generators/recipes/pkg"
	"github.com/spf13/pflag"
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

// Example struct matching the YAML structure
type Shop struct {
	Name       string     `yaml:"name"`
	Categories []Category `yaml:"categories"`
}

type Category struct {
	Name string   `yaml: "name"`
	List []string `yaml: "list"`
}

func readShopCategories(filename string) (*Shop, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var shop Shop
	if err := yaml.Unmarshal(data, &shop); err != nil {
		return nil, err
	}

	return &shop, nil
}

func constructIngredientIndex(shop *Shop) map[string]int {
	mapping := make(map[string]int)
	for idx, category := range shop.Categories {
		for _, item := range category.List {
			mapping[item] = idx
		}
	}
	return mapping
}

func printIngredientsByShopIndex(ingSet *pkg.IngredientsSet, shop *Shop) {
	mapping := constructIngredientIndex(shop)

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

	var recipeFilenames []string
	var shopFilename string
	pflag.StringSliceVar(&recipeFilenames, "recipe", []string{}, "List of recipes")
	pflag.StringVar(&shopFilename, "shop", "", "Shop to visit")
	pflag.Parse()

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
			log.Fatal(err)
			return
		}

		for _, step := range r.Steps {
			ingSet.MergeIngredients(step.Ingredients)
		}
	}

	ingSet.Consolidate()

	printIngredientsByShopIndex(ingSet, shop)
}
