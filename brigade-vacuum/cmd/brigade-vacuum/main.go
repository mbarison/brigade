package main

import (
	"fmt"
	"os"

	"github.com/mbarison/brigade/brigade-vacuum/cmd/brigade-vacuum/commands"
)

func main() {
	if err := commands.Root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
