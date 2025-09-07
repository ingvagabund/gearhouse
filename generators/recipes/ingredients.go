package main

import (
	"fmt"
	"strings"

	"github.com/aquilax/cooklang-go"
)

type ingredientAmountList []cooklang.IngredientAmount

type ingredientsSet struct {
	ingredients map[string]ingredientAmountList
}

func newIngredientsSet() *ingredientsSet {
	return &ingredientsSet{
		ingredients: make(map[string]ingredientAmountList),
	}
}

func (ingSet *ingredientsSet) mergeIngredients(ingredients []cooklang.Ingredient) {
	for _, ing := range ingredients {
		if _, ok := ingSet.ingredients[ing.Name]; !ok {
			ingSet.ingredients[ing.Name] = ingredientAmountList{ing.Amount}
		} else {
			ingSet.ingredients[ing.Name] = append(ingSet.ingredients[ing.Name], ing.Amount)
		}
	}
}

func (ingSet *ingredientsSet) consolidate() {
	for ing := range ingSet.ingredients {
		units := make(map[string]ingredientAmountList)
		for _, amount := range ingSet.ingredients[ing] {
			if _, exists := units[amount.Unit]; !exists {
				units[amount.Unit] = ingredientAmountList{}
			}
			units[amount.Unit] = append(units[amount.Unit], amount)
		}
		ingSet.ingredients[ing] = nil
		for unit := range units {
			totalQuantity := float64(0)
			isNumeric := false
			for _, amount := range units[unit] {
				totalQuantity += amount.Quantity
				isNumeric = amount.IsNumeric
			}
			ingSet.ingredients[ing] = append(ingSet.ingredients[ing], cooklang.IngredientAmount{
				IsNumeric: isNumeric,
				Quantity:  totalQuantity,
				Unit:      unit,
			})
		}
	}
}

func (list ingredientAmountList) toString() string {
	strList := []string{}
	for _, amount := range list {
		if amount.IsNumeric {
			if amount.Unit == "" {
				strList = append(strList, fmt.Sprintf("%v ks", amount.Quantity))
			} else {
				strList = append(strList, fmt.Sprintf("%v %v", amount.Quantity, amount.Unit))
			}
		}
	}
	if len(strList) == 0 {
		return ""
	}
	return strings.Join(strList, " + ")
}

func (ingSet *ingredientsSet) toString() {
	for ing := range ingSet.ingredients {
		str := ingSet.ingredients[ing].toString()
		if str == "" {
			fmt.Printf("%v\n", ing)
		} else {
			fmt.Printf("%v: %v\n", ing, str)
		}
	}
}
