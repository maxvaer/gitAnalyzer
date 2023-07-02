// Package cmd contains all code used by cobra for the cli.
package cmd

import (
	"GitAnalyzer/internal/Modules"
	"github.com/spf13/cobra"
	"log"
)

// crawlCmd represents the crawl command
var crawlCmd = &cobra.Command{
	Use:   "crawl",
	Short: "Crawl public GitHub repositories.",
	Long: `Crawl all known public GitHub repositories.
The results will be saved in the corresponding database.`,
	Run: func(cmd *cobra.Command, args []string) {
		filePath, err := cmd.Flags().GetString("file")
		if err != nil {
			log.Fatalln("Error parsing file flag:", err.Error())
		}
		if filePath != "" {
			Modules.InsertDataFromCSV(filePath)
		}

		Modules.Crawl()
	},
}

func init() {
	rootCmd.AddCommand(crawlCmd)

	crawlCmd.Flags().StringP("file", "f", "", "Insert crawled data from a CSV file.")
}
