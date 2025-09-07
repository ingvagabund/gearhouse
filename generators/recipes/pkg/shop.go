package pkg

// Example struct matching the YAML structure
type Shop struct {
	Name       string     `yaml:"name"`
	Categories []Category `yaml:"categories"`
}

type Category struct {
	Name string   `yaml: "name"`
	List []string `yaml: "list"`
}

func ConstructIngredientIndex(shop *Shop) map[string]int {
	mapping := make(map[string]int)
	for idx, category := range shop.Categories {
		for _, item := range category.List {
			mapping[item] = idx
		}
	}
	return mapping
}
