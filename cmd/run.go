// Package cmd contains all code used by cobra for the cli.
package cmd

import (
	"GitAnalyzer/api/Analyzer"
	"GitAnalyzer/internal/Modules"
	"github.com/spf13/cobra"
	"log"
)

// runCmd represents the run command which start the gitAnalyzer process
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Analyze the given repositories",
	Long:  `The analyzer executes the selected templates on all repositories.`,
	Run: func(cmd *cobra.Command, args []string) {
		urlFilePath, err := cmd.Flags().GetString("url-file")
		if err != nil {
			log.Fatalln("Error parsing urlFilePath flag:", err.Error())
		}
		tags, errTags := cmd.Flags().GetString("filter")
		if errTags != nil {
			log.Fatalln("Error parsing filter flag:", errTags.Error())
		}
		templatesPath, errTemplates := cmd.Flags().GetString("templates")
		if errTemplates != nil {
			log.Fatalln("Error parsing templates flag:", errTemplates.Error())
		}
		workerCount, errWorkerCount := cmd.Flags().GetInt("worker-count")
		if errWorkerCount != nil {
			log.Fatalln("Error parsing workerCount flag:", errWorkerCount.Error())
		}
		keepData, errKeepData := cmd.Flags().GetBool("keep-data")
		if errKeepData != nil {
			log.Fatalln("Error parsing keep-data flag:", errKeepData.Error())
		}
		verbose, errVerbose := cmd.Flags().GetBool("verbose")
		if errVerbose != nil {
			log.Fatalln("Error parsing verbose flag:", errVerbose.Error())
		}
		excluded, errExcluded := cmd.Flags().GetString("excluded")
		if errExcluded != nil {
			log.Fatalln("Error parsing excluded templates names flag:", errExcluded.Error())
		}
		results, errResults := cmd.Flags().GetString("results")
		if errTemplates != nil {
			log.Fatalln("Error parsing results flag:", errResults.Error())
		}

		config := Analyzer.Config{UrlFilePath: urlFilePath, Tags: tags,
			TemplatesPath: templatesPath, WorkerCount: workerCount, KeepData: keepData, Excluded: excluded,
			ResultsDir: results, Verbose: verbose}
		Modules.Run(config)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringP("filter", "f", "", "Tags to filter templates.")
	runCmd.Flags().StringP("url-file", "u", "./urls.csv", "Path of the file containing repository URLs.")
	runCmd.Flags().StringP("templates", "t", "./templates", "Path of the template directory.")
	runCmd.Flags().StringP("excluded", "e", "", "Names of excluded templates.(comma seperated)")
	runCmd.Flags().StringP("results", "r", "./results", "Path of the results directory.")
	runCmd.Flags().IntP("worker-count", "c", 5, "Number of concurrent workers.")
	runCmd.Flags().Bool("keep-data", false, "Don't delete the cloned repositories.")
	runCmd.Flags().Bool("verbose", false, "Show verbose output.")
}
