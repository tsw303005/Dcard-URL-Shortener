package main

import (
	"log"

	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use:   "url-shortener [module]",
		Short: "url-shortener entrypoint",
	}

	// wait for api done~
	cmd.AddCommand()

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
