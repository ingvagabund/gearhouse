# Shopping list generator

Generate a list of ingredients based on provided recipes.

## Usage

Recipes are written in [cooklang](https://cooklang.org/).
A shop is currently a simple yaml file.

```bash
$ make build
$ ./_output/bin/recipes2ingredients \
  --recipe data/recipes/gnocchi-with-chicken-and-spinach.cook \
  --recipe data/recipes/houbove-rizotto.cook \
  --recipe data/recipes/sekana-z-cervene-repy.cook \
  --shop data/shops/globus-brno-ivanovice.yaml
```
