// Package cmd contains all code used by cobra for the cli.
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "GitAnalyzer",
	Short: "A template based SAST scanner.",
	Long: `A template based SAST scanner for git repositories.
GitAnalyzer is able to scan trough the whole git history of a repository.
The crawl mode can be used to fetch the URLs of all public GitHub repositories.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

}
