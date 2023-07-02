// Package Analyzer contains all structural components of the application.
package Analyzer

// The Stat struct is used to store the information
// which repository was scanned and how long did it took.
type Stat struct {
	// URL of the scanned repository
	URL string `csv:"url"`
	// Time it took to scan the repository
	ElapsedTime string `csv:"elapsed_time"`
}
