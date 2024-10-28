package main

import (
	"github.com/rusinikita/goviz/internal"
)

const (
	projectDir         = "/Users/nikitarusin/Repositories/replace_me"
	pathPrefixToRemove = "github.com/rusinikita/replace_me/"
)

func main() {
	code := internal.Compile(projectDir, pathPrefixToRemove)

	internal.RenderFiles(code)
}
