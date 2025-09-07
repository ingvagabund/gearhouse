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

package pkg

import (
	"fmt"
	"strings"

	"github.com/aquilax/cooklang-go"
)

type IngredientAmountList []cooklang.IngredientAmount

type IngredientsSet struct {
	Ingredients map[string]IngredientAmountList
}

func NewIngredientsSet() *IngredientsSet {
	return &IngredientsSet{
		Ingredients: make(map[string]IngredientAmountList),
	}
}

func (ingSet *IngredientsSet) MergeIngredients(ingredients []cooklang.Ingredient) {
	for _, ing := range ingredients {
		if _, ok := ingSet.Ingredients[ing.Name]; !ok {
			ingSet.Ingredients[ing.Name] = IngredientAmountList{ing.Amount}
		} else {
			ingSet.Ingredients[ing.Name] = append(ingSet.Ingredients[ing.Name], ing.Amount)
		}
	}
}

func (ingSet *IngredientsSet) Consolidate() {
	for ing := range ingSet.Ingredients {
		units := make(map[string]IngredientAmountList)
		for _, amount := range ingSet.Ingredients[ing] {
			if _, exists := units[amount.Unit]; !exists {
				units[amount.Unit] = IngredientAmountList{}
			}
			units[amount.Unit] = append(units[amount.Unit], amount)
		}
		ingSet.Ingredients[ing] = nil
		for unit := range units {
			totalQuantity := float64(0)
			isNumeric := false
			for _, amount := range units[unit] {
				totalQuantity += amount.Quantity
				isNumeric = amount.IsNumeric
			}
			ingSet.Ingredients[ing] = append(ingSet.Ingredients[ing], cooklang.IngredientAmount{
				IsNumeric: isNumeric,
				Quantity:  totalQuantity,
				Unit:      unit,
			})
		}
	}
}

func (list IngredientAmountList) ToString() string {
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

func (ingSet *IngredientsSet) ToString() {
	for ing := range ingSet.Ingredients {
		str := ingSet.Ingredients[ing].ToString()
		if str == "" {
			fmt.Printf("%v\n", ing)
		} else {
			fmt.Printf("%v: %v\n", ing, str)
		}
	}
}

func (ingSet *IngredientsSet) ToCookLang(singleLine bool) string {
	list := []string{}
	for ingName := range ingSet.Ingredients {
		for _, ing := range ingSet.Ingredients[ingName] {
			if !ing.IsNumeric {
				list = append(list, fmt.Sprintf("@%v{}", ingName))
			} else {
				if ing.Unit == "" {
					list = append(list, fmt.Sprintf("@%v{%v}", ingName, ing.Quantity))
				} else {
					list = append(list, fmt.Sprintf("@%v{%v%%%v}", ingName, ing.Quantity, ing.Unit))
				}
			}
		}
	}
	if singleLine {
		return strings.Join(list, " ")
	}
	return strings.Join(list, "\n")
}
