package main

import (
	"os"

	"github.com/spf13/pflag"
	"k8s.io/klog/v2"
)

var (
	recipeFilenames []string
	shopFilename    string
)

func initFlags() {
	pflag.StringSliceVar(&recipeFilenames, "recipe", []string{}, "List of recipes")
	pflag.StringVar(&shopFilename, "shop", "", "Shop to visit")
	pflag.Parse()
}

func validateFlags() {
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
}
