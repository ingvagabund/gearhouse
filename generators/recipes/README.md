# Shopping list generator

Generate a list of ingredients based on provided recipes.

## Usage

Recipes are written in [cooklang](https://cooklang.org/).
A shop is currently a simple yaml file.

```bash
$ make build
$ ./_output/bin/ingredients \
  --recipe data/gnocchi-with-chicken-and-spinach.cook \
  --recipe data/houbove-rizotto.cook \
  --recipe data/sekana-z-cervene-repy.cook \
  --shop data/globus-brno-ivanovice.yaml
```
