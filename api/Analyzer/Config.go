// Package Analyzer contains all structural components of the application.
package Analyzer

// Config is the struct used to adjust the settings for the gitAnalyzer
type Config struct {
	// Path of the csv file which contains URLs to GitHub repositories
	UrlFilePath string
	// Tags used to filter the loaded templates
	Tags string
	// Path to the directory which contains the template YAML files
	TemplatesPath string
	// Path to the directory where the results will be stored
	ResultsDir string
	// Excluded template names use to filter all loaded templates
	Excluded string
	// The WorkerCount is used to adjust the number of workers in the worker-pool
	WorkerCount int
	// If KeepData is set to true, the repositories will not be deleted after the scan
	KeepData bool
	// Verbose can be used to get a more detailed output
	Verbose bool
}
