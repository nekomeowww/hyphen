package main

import (
	"fmt"

	"github.com/nekomeowww/hyphen/pkg/meta"
	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of hyphen",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Printf("hyphen version %s\n", meta.Version)
		},
	}
)
