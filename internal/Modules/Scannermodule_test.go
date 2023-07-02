// Package Modules contains all business logic modules/components of the application..
package Modules

import (
	"GitAnalyzer/api/Analyzer"
	"fmt"
	"regexp"
	"strings"
	"testing"
)

func TestRegularExpressions(t *testing.T) {
	// Initialize handlers
	config := Analyzer.Config{TemplatesPath: "../../templates"}
	repoHandler := NewRepoHandler(&GitHelper{}, config)
	th := NewTemplateHandlerWithMocks(repoHandler, &FileHandler{}, &CommandHandler{}, config)

	// Load all templates
	th.LoadTemplates(config.TemplatesPath)

	// Iterate over templates
	for _, template := range th.templates {
		// Check if the current template has no regular expression.
		if len(template.Regex) == 0 {
			continue
		}

		// Iterate over all regular expressions of the current template
		for _, regex := range template.Regex {

			// Compile regular expression
			expression := regexp.MustCompile(strings.TrimSpace(regex.Expression))

			// Iterate over test vales for the current regular expression
			for _, test := range regex.Tests {
				// Initialize map of foundString => true
				valuesFound := make(map[string]bool)
				// Execute regex against test input
				matchesFound := expression.FindAllStringSubmatch(test.Input, -1)
				// Add found values to the map
				for _, got := range matchesFound {
					value := got[regex.Group]
					value = strings.TrimSpace(value)
					valuesFound[value] = true
				}

				// Iterate over expected values
				for _, expected := range test.Want {
					if !valuesFound[expected] {
						// Check case where wanted equals ""
						if expected == "" && matchesFound == nil {
							continue
						}
						// Value was not found => Failed expression
						t.Errorf("Expression failed! Wanted value not found: %s", expected)
						continue
					} else {
						// Expected value was found
						fmt.Println(strings.TrimSpace(regex.Description), "found:", expected)
					}
				}

			}

		}
	}
}
