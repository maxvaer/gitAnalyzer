// Package Modules contains all business logic modules/components of the application.
package Modules

import (
	"GitAnalyzer/api/Analyzer"
	"GitAnalyzer/pkg/Utils"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// IRepoHelper the interface is used to define which actions can be called on the repoHelper
//
//go:generate mockery --name IRepoHelper
type IRepoHelper interface {
	GetPathOfRepository(repo *git.Repository) string
	GetGitHubURLOfRepository(repo *git.Repository) string
	DeleteRepository(repo *git.Repository)
	GetHeadCommit(repo *git.Repository) (*object.Commit, error)
	Checkout(repo *git.Repository, opts *git.CheckoutOptions) error
	ResetRepository(repo *git.Repository, headCommit *object.Commit)
	GetRemoteBranches(repo *git.Repository) []*plumbing.Reference
	GetTemplateTasks(repo *git.Repository, templates []Analyzer.Template) []Analyzer.TemplateTask
}

// The RepoHandler struct is responsible to handle all actions
// to control and manage repositories
type RepoHandler struct {
	// The gitHelper used to control call the git actions
	gitHelper IGitHelper
	// The used Analyzer.Config containing the settings of the current scan
	config Analyzer.Config
}

// NewRepoHandler is the constructor to create a new  repoHandler with given IGitHelper interface
// and Analyzer.Config configuration.
func NewRepoHandler(helper IGitHelper, config Analyzer.Config) *RepoHandler {
	repoHandler := &RepoHandler{gitHelper: helper, config: config}
	return repoHandler
}

// IGitHelper allows a repository to be cloned or opened
//
//go:generate mockery --name IGitHelper
type IGitHelper interface {
	// Clone a repository by its URL into a given directory
	Clone(dir, url string) (*git.Repository, error)
	// Open an already cloned repository from a given directory
	Open(dir string) (*git.Repository, error)
}

// GitHelper struct used for an implementation of the IGitHelper interface
type GitHelper struct{}

// Clone a repository by a URL into a given directory using the go-git library clone function.
func (g *GitHelper) Clone(dir, url string) (*git.Repository, error) {
	return git.PlainClone(dir, false, &git.CloneOptions{
		URL: url,
	})
}

// Open an already cloned repository from a given directory using the go-git library open function.
func (g *GitHelper) Open(dir string) (*git.Repository, error) {
	return git.PlainOpen(dir)
}

// CloneRepositories is the facade function of repoHandler to
// clone or open a repository from a given task.
// The cloned repository is then returned
func (rh *RepoHandler) CloneRepositories(task Analyzer.Task) *git.Repository {
	fmt.Println("Cloning:", task.URL)
	repo := rh.cloneOrOpenByURL(task.URL, "./repos/")
	if rh.config.Verbose {
		fmt.Println("Cloned:", task.URL)
	}
	return repo
}

// GetTemplateTasks returns the slice of Analyzer.TemplateTask for a given repository and slice of templates
// A template.task contains a commit hash and slice of templates to run on that commit.
func (rh *RepoHandler) GetTemplateTasks(repo *git.Repository, templates []Analyzer.Template) []Analyzer.TemplateTask {
	// Initialize return slice
	var templateTask []Analyzer.TemplateTask

	//Map commit => []Templates
	commitsToTemplates := Analyzer.NewCommitToTemplatesMap()

	// Iterate over all provided templates
	for _, template := range templates {
		var commits []*object.Commit
		// Switch based on the type of the current template
		switch template.Type {
		case "Full":
			// Iterate over all branches and all commits

			// Create set of all commit hashes
			commitHashset := make(map[plumbing.Hash]struct{})
			// Get all branches of the repository
			branches := rh.GetRemoteBranches(repo)
			// Iterate over branches
			for _, branch := range branches {
				// Checkout the current branch
				err := rh.Checkout(repo, &git.CheckoutOptions{Branch: branch.Name()})
				if err != nil {
					// Continue with the next branch on error e.g. broken branch
					continue
				}
				// Get the commits for the current branch
				commitsOfBranch := rh.GetCommitsForBranch(repo)
				// Iterate over all commits
				for _, commit := range commitsOfBranch {
					// Check if the current commit is already known/found
					if _, found := commitHashset[commit.Hash]; found {
						// Continue with next commit if already found
						continue
					}
					// Append unknown commit to set
					commits = append(commits, commit)
					commitHashset[commit.Hash] = struct{}{}
				}
			}
		case "Deep":
			// Get all commits for the current branch
			commits = rh.GetCommitsForBranch(repo)
		default:
			// If not typ was provided, use the flat type as default
			fallthrough
		case "Flat":
			// Just use the current HEAD
			headCommit, err := rh.GetHeadCommit(repo)
			if err == nil {
				commits = append(commits, headCommit)
			}
		}

		// Check if number of maxCommits is reached for the current template
		if template.MaxCommits > 0 && template.MaxCommits >= len(commits) {
			// If more than maxCommits are found, reslice to maxCommits
			commits = commits[:template.MaxCommits]
		}

		// Iterate over found commits
		for _, commit := range commits {
			// Append commit to commit => []Templates
			commitsToTemplates.Insert(commit.Hash.String(), template)
		}
	}

	// Iterate over all commit => []Templates
	for _, commit := range commitsToTemplates.Commits() {
		// Get the all templates to run for the current commit
		templatesForCommit := commitsToTemplates.GetTemplates(commit)
		// Create templateTask for current commit and append to result slice
		templateTask = append(templateTask, Analyzer.TemplateTask{
			CommitHash: commit,
			Templates:  templatesForCommit,
		})
	}

	return templateTask
}

// DeleteRepository removes the given repository from the file system
func (rh *RepoHandler) DeleteRepository(repo *git.Repository) {
	// Get the file path of the repository
	path := rh.GetPathOfRepository(repo)
	// Remove the repo directory
	err := os.RemoveAll(path)
	if err != nil {
		log.Fatal("Error trying to delete repo:", err)
	}

	// Get the parent directory
	parentDir := filepath.Dir(path)
	// Check if the parent directory is als empty
	isEmpty, err := Utils.IsDirEmpty(parentDir)
	if err != nil {
		fmt.Println("Error checking if parent dir is Empty:", err.Error())
	}
	if isEmpty {
		// If the parent directory is empty, remove it also
		err = os.RemoveAll(parentDir)
		if err != nil {
			log.Fatal("Error trying to delete repo parent dir:", err)
		}
	}

	// If verbose, echo that the given repository was deleted
	if rh.config.Verbose {
		fmt.Println("Removed repository:", path)
	}
}

// cloneOrOpenByURL uses the provided url and base directory to either clone or open an existing repository
func (rh *RepoHandler) cloneOrOpenByURL(url, baseDir string) *git.Repository {
	//Replace git protocol with https
	url = Utils.NormalizeGitURLToHTTPS(url)
	//Build repository path
	subs := strings.Split(url, "/")
	baseDir = strings.TrimSuffix(baseDir, "/")
	baseDir = strings.TrimSuffix(baseDir, string(os.PathSeparator))
	dir := baseDir + string(os.PathSeparator) + subs[3] + string(os.PathSeparator) + strings.Split(subs[4], ".")[0]
	//Check if repo folder does not exist
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		//Clone Repository
		repo, errClone := rh.gitHelper.Clone(dir, url)
		if errClone == nil {
			return repo
		} else {
			fmt.Println("Error cloning:", errClone)
			err = os.RemoveAll(dir)
			if err != nil {
				fmt.Println("Error removing dir of broken repository:", err.Error())
				return nil
			}
			return nil
		}
	} else {
		// If repository exists, open the repository
		repo, errOpen := rh.gitHelper.Open(dir)
		if errOpen == nil {
			return repo
		} else {
			fmt.Println("Error opening:" + errOpen.Error())
			return nil
		}
	}
}

// GetCommitsForBranch returns a list of pointers to commits of the currently checked out branch
// from the provided repository
func (rh *RepoHandler) GetCommitsForBranch(repo *git.Repository) []*object.Commit {
	// Initialize result slice
	var commits []*object.Commit

	// Retrieves the commit history
	cIter, err := repo.Log(&git.LogOptions{
		Order: git.LogOrderCommitterTime,
	})
	if err != nil {
		fmt.Println("Error fetching git Log:" + err.Error())
		return nil
	}

	// Iterate over all commits using iterator
	err = cIter.ForEach(func(c *object.Commit) error {
		// Append current commit to result
		commits = append(commits, c)
		return nil
	})
	if err != nil {
		fmt.Println("Error iterating commits:" + err.Error())
		return nil
	}

	return commits
}

// GetRemoteBranches returns all branches for the provided repository
func (rh *RepoHandler) GetRemoteBranches(repo *git.Repository) []*plumbing.Reference {
	// Initialize return slice
	var branches []*plumbing.Reference

	// Get the reference iterator
	refIter, err := repo.References()
	if err != nil {
		//fmt.Println("Error fetching references:", err.Error())
		return nil
	}

	// Iterate over all references
	refIter.ForEach(func(reference *plumbing.Reference) error {
		// Check if the current reference is a remote branch and not e.g. a tag
		if reference.Name().IsRemote() {
			// Append to result
			branches = append(branches, reference)
		}
		return nil
	})

	return branches
}

// Checkout checks out a commit provided via the CheckoutOptions for the given repository
func (rh *RepoHandler) Checkout(repo *git.Repository, options *git.CheckoutOptions) error {
	// Get the worktree
	workTree, err := repo.Worktree()
	if err != nil {
		fmt.Println("Error creating worktree:", err.Error())
	}

	// Call checkout with the provided options
	err = workTree.Checkout(options)
	return err
}

// GetHeadCommit get the HEAD commit of the given repository
func (rh *RepoHandler) GetHeadCommit(repo *git.Repository) (*object.Commit, error) {
	// Get the HEAD reference
	ref, err := repo.Head()
	if err != nil {
		//fmt.Println("Error fetching head for repository:", err)
		return nil, err
	}

	// Get the commit object for the current HEAD commit
	headCommit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		//fmt.Println("Error fetching head commit for repository:", err)
		return nil, err
	}

	// Return HEAD commit reference and no error
	return headCommit, nil
}

// ResetRepository resets the given repository to the given HEAD commit
func (rh *RepoHandler) ResetRepository(repo *git.Repository, headCommit *object.Commit) {
	// Get the worktree
	worktree, err := repo.Worktree()
	if err != nil {
		fmt.Println("Error fetching worktree for repository:", err)
	}

	// Reset the worktree to the HEAD commit
	errReset := worktree.Reset(&git.ResetOptions{
		Mode:   git.HardReset,
		Commit: headCommit.Hash,
	})
	if errReset != nil {
		fmt.Println("Error reseting repository:", errReset)
	}

}

// GetPathOfRepository returns the file system path of the given repository
func (rh *RepoHandler) GetPathOfRepository(repo *git.Repository) string {
	// Get the worktree
	worktree, err := repo.Worktree()
	if err != nil {
		fmt.Println("Error opening Worktree")
		return ""
	}
	// Get the filesystem of the worktree
	filesystem := worktree.Filesystem

	// Return the path of the filesystem
	return filesystem.Root()
}

// GetGitHubURLOfRepository returns the GitHub URL of the given repository
func (rh *RepoHandler) GetGitHubURLOfRepository(repo *git.Repository) string {
	// Get the path of the repository
	path := rh.GetPathOfRepository(repo)
	// Split the path by the path separator of the current OS
	splits := strings.Split(path, string(os.PathSeparator))

	// return URL + username + reponame
	return "https://github.com/" + splits[len(splits)-2] + "/" + splits[len(splits)-1] + ""
}
