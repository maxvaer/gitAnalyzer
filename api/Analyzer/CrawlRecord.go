// Package Analyzer contains all structural components of the application.
package Analyzer

// The CrawlRecord is a struct, which is used to pass the loaded
// data from the GitHub API to the database
type CrawlRecord struct {
	// URL of the GitHub repository
	URL string
	// ID of the GitHub repository
	ID string
	// Creation date of the GitHub repository
	CreationDate string
	// The number of time the GitHub repository was forked
	ForkCount string
	// The Size of the GitHub repository in Bytes
	Size string
	// The number of Stars of the GitHub repository
	Stars string
	// The last time the GitHub repository was updated
	UpdateDate string
	// The primary programming language used inside the GitHub repository
	Language string
	// The number of commits from the GitHub repository
	CommitCount string
}
