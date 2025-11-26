package main

import (
	"analiser/cmd/console"
	"analiser/cmd/web"
	"os"
)

func main() {
	osArgs := os.Args[1:]
	if len(osArgs) > 0 {
		console.Process()
	} else {
		web.Process()
	}
}
