package main

import (
	"os"

	"github.com/hagatasdelus/agysession/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
