package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/aquilax/cooklang-go"
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

func main() {
	ingSet := newIngredientsSet()

	recipeFiles := []string{
		"data/gnocchi-with-chicken-and-spinach.cook",
		"data/houbove-rizotto.cook",
		"data/sekana-z-cervene-repy.cook",
	}
	for _, recipeFile := range recipeFiles {
		r, err := readRecipeFromFile(recipeFile)
		if err != nil {
			log.Fatal(err)
			return
		}

		for _, step := range r.Steps {
			ingSet.mergeIngredients(step.Ingredients)
		}
	}

	ingSet.toString()
	ingSet.consolidate()
	fmt.Printf("\n\n")
	ingSet.toString()
}
