// Package Analyzer contains all structural components of the application.
package Analyzer

// The Meta struct is used to provide some meta information for a template
type Meta struct {
	// References like links to websites etc.
	References []string `yaml:"references"`
	// The Impact of the result/vulnerability
	Impact string `yaml:"impact"`
	// The CWE category of the result/vulnerability
	CWE        string `yaml:"CWE"`
	CVSS       string `yaml:"CVSS"`
	CVE        string `yaml:"CVE"`
	Mitigation string `yaml:"mitigation"`
}

// The Script struct describes the command which can be executed
type Script struct {
	// Language of the provided Code, can be: bash, python, cli (default)
	Language string `yaml:"language"`
	// The Code which will be executed
	Code string `yaml:"code"`
}

// The Match struct is used to provide information to find and filter files inside the repository
type Match struct {
	// Filenames to search for
	Filename []string `yaml:"filename"`
	// FileEndings to search for
	FileEndings []string `yaml:"file_endings"`
	// Keywords which can be a part of the path of a file
	Keywords []string `yaml:"keywords"`
	// If an Exclude is found and matched, the file/path will not be returned
	Exclude []string `yaml:"exclude"`
}

// The Test struct is used within a Validation to unit-test a regular expression
type Test struct {
	// The input string for the test
	Input string `yaml:"input"`
	// The expected output of the test
	Want []string `yaml:"want"`
}

// The Validation struct is used to provide functionality to test a regular expression
type Validation struct {
	// A Repository which is vulnerable to the given regular expression
	Repository string `yaml:"repository"`
	// A number of unit-tests to verify the regular expression
	Tests []Test `yaml:"tests"`
}

// The Regex struct defines a regular expression and can be used to identify vulnerabilities
// or misconfigurations in files.
type Regex struct {
	// The Description of the regular expression
	Description string `yaml:"description"`
	// The Expression used as a regular expression
	Expression string `yaml:"expression"`
	// The Group is used to select one of the matches from the results of the Expression
	Group int `yaml:"group"`
	// FileEndings to search for
	FileEndings []string `yaml:"file_endings"`
	// FalsePositives can be used to filter the result for known non-valid outputs.
	FalsePositives []string `yaml:"false_positives"`
	// Filenames to search for
	Filename []string `yaml:"filename"`
	// If an Exclude is found in a matched, the file/path will not be returned
	Exclude []string `yaml:"exclude"`
	// References like links to websites etc.
	References []string `yaml:"references"`
	// A Validation struct to verify that the regular expression is working correctly
	Validation `yaml:"validation"`
}

// The Output struct can be used to define further processing of the output
type Output struct {
	// If Unique is set, a new file will be generated for the given template,
	// which will contain only unique values for the results.
	Unique bool `yaml:"uniq"`
}

// The Requirements struct is used to define which tools/packages are needed by a template
type Requirements struct {
	// Required Tools
	Tools []string `yaml:"tools"`
	// Required Pip-packages
	Pip []string `yaml:"pip"`
	// Required Npm-packages
	Npm []string `yaml:"npm"`
}

// The Template struct defines a Template used by the gitAnalyzer application
// to scan a GitHub Repository
type Template struct {
	// The Name of the Template
	Name string `yaml:"name"`
	// The Description of what the Template checks
	Description string `yaml:"description"`
	// The Requirements struct to define needed tools/packages etc.
	Requirements `yaml:"requirements"`
	// Tags is a slice of strings used to filter templates
	Tags []string `yaml:"tags"`
	// Type defines which commits will be checked by the template
	// Valid options are: "Full" to check all commits on all branches, Flat (default) to check the HEAD commit on the main branch
	// and Deep to check all commits on the main branch
	Type string `yaml:"type"`
	// MaxCommits can be used to abort the template after a set number of searched commits
	// This can prevent a tasks from scanning "endlessly" if a repository contains > 100.000 commits
	MaxCommits int `yaml:"max_commits"`
	// A list of Regex is used to identify vulnerabilities or misconfiguration inside the repository
	Regex []Regex `yaml:"regex"`
	// The Script struct can be used to execute commands inside the repository
	Script `yaml:"script"`
	// The Output struct can be used to define further processing of the output
	Output `yaml:"output"`
	// The Match struct is used to provide information to find and filter files inside the repository
	Match `yaml:"match"`
	// The Meta struct is used to provide some meta information for the template
	Meta `yaml:"meta"`
	// The PreScript can contain a Script which will be executed once before the scans starts.
	// This can be used to set up docker container etc.
	PreScript Script `yaml:"pre_script"`
}
