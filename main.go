package main

import (
	"fmt"
	"os"

	"github.com/rusinikita/goviz/internal"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Usage: goviz <project dir>")
		os.Exit(1)
	}

	code := internal.Compile(args[1])

	internal.RenderFiles(code)
}
