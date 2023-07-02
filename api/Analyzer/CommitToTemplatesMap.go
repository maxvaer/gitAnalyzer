// Package Analyzer contains all structural components of the application.
package Analyzer

// CommitToTemplatesMap Maps a slice of commit hashes to
// a number of templates
type CommitToTemplatesMap struct {
	// Slice of commit hashes
	commitHashes []string
	// Map of commitHash => templates
	items map[string][]Template
}

// NewCommitToTemplatesMap is a constructor to create a new CommitToTemplatesMap
func NewCommitToTemplatesMap() *CommitToTemplatesMap {
	return &CommitToTemplatesMap{
		commitHashes: make([]string, 0),
		items:        make(map[string][]Template),
	}
}

// Insert adds a commit hash and template to the CommitToTemplatesMap object
func (m *CommitToTemplatesMap) Insert(commit string, template Template) {
	if _, found := m.items[commit]; !found {
		m.commitHashes = append(m.commitHashes, commit)
	}
	m.items[commit] = append(m.items[commit], template)
}

// GetTemplates returns all templates for a given commit hash
func (m *CommitToTemplatesMap) GetTemplates(commit string) []Template {
	templates, _ := m.items[commit]
	return templates
}

// Commits returns the list of all commit hashes
func (m *CommitToTemplatesMap) Commits() []string {
	return m.commitHashes
}
