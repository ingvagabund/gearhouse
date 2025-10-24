# Shopping list generator

Generate a list of ingredients based on provided recipes.

## Usage

Recipes are written in [cooklang](https://cooklang.org/).
A shop is currently a simple yaml file.

```bash
$ make build
$ ./_output/bin/recipeselector \
  --history-filename examples/history.dat \
  --segment-length 3 \
  --config config/recipeselector.yaml
```
