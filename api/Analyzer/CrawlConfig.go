// Package Analyzer contains all structural components of the application.
package Analyzer

// The CrawlConfig struct is used to adjust the settings for the crawling module.
type CrawlConfig struct {
	// The GitHub private access token used to call the GitHub API
	GitHubAPIToken string
	// Thr User of the database
	DBUser string
	// Thr password of the previous user
	DBPassword string
	// The IP of the database server
	DBIP string
	// The port of the database server
	DBPort string
}
