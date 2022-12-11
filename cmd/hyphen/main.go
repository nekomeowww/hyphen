package main

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "hyphen",
		Short: "Hyphen is a url shortener service",
	}
)

var (
	listen string
	dbPath string
)

func init() {
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(versionCmd)
}

func main() {
	runCmd.Flags().StringVarP(&listen, "listen", "l", "", "server listening address")
	runCmd.Flags().StringVar(&dbPath, "data", "", "path to store bbolt data file")

	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
