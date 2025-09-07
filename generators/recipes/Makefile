.PHONY: build
build:
	CGO_ENABLED=0 go build -o _output/bin/ingredients ./cmd

.PHONY: example
example: build
	./_output/bin/ingredients \
		--recipe data/gnocchi-with-chicken-and-spinach.cook \
		--recipe data/houbove-rizotto.cook \
		--recipe data/sekana-z-cervene-repy.cook \
		--shop data/globus-brno-ivanovice.yaml
