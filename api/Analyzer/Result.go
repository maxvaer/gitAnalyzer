// Package Analyzer contains all structural components of the application.
package Analyzer

// The Result struct is used to store the result of a template
type Result struct {
	// The name of the template
	TemplateName string `csv:"-"`
	// The GitHub URL of the repository
	URL string `csv:"url"`
	// The commit hash where the result was found
	CommitHash string `csv:"commit_hash"`
	// The Timestamp displays when the result was found
	Timestamp string `csv:"timestamp"`
	// The Path of the file inside the repository where the result was found
	Path string `csv:"file_path"`
	// The Description of the found result
	Description string `csv:"description"`
	// The output of the template command or regular expression
	Output string `csv:"output"`
}
