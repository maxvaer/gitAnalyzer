// Package Modules contains all business logic modules/components of the application.
package Modules

import (
	"GitAnalyzer/api/Analyzer"
	"GitAnalyzer/internal/mocks"
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/suite"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)

// Create test suite for the repoHandler module
type RepoModuleTestSuite struct {
	suite.Suite
	tempDir       string
	mockGitHelper *mocks.IGitHelper
	repoHandler   *RepoHandler
	repo          *git.Repository
	worktree      *git.Worktree
}

// SetupTest is run before every test of the test suite to initialize a clear state
func (suite *RepoModuleTestSuite) SetupTest() {
	// Create a temporary directory
	suite.tempDir = suite.T().TempDir()
	// Create a mocked GitHelper
	suite.mockGitHelper = mocks.NewIGitHelper(suite.T())
	// Create a repoHandler with the mocked GitHelper
	suite.repoHandler = NewRepoHandler(suite.mockGitHelper, Analyzer.Config{})
	// Initialize a test git repository
	repo, err := git.PlainInit(suite.tempDir, false)
	if err != nil {
		log.Fatalln("Error initializing in memory repository:", err.Error())
	}
	// Set the repository
	suite.repo = repo
	// Get and set the worktree of the git repository
	worktree, err := suite.repo.Worktree()
	if err != nil {
		log.Fatalln("Error getting worktree for test repository:", err.Error())
	}
	suite.worktree = worktree
}

// TestCloneOrOpenByURL_Clone tests the clone part of the cloneOrOpenByURL function
func (suite *RepoModuleTestSuite) TestCloneOrOpenByURL_Clone() {
	// Creat the GitHub URL based on the path
	subs := strings.Split(suite.tempDir, string(os.PathSeparator))
	url := "https://github.com/" + subs[len(subs)-2] + "/" + subs[len(subs)-1]

	// Set return value for clone function of the mocked GitHelper
	suite.mockGitHelper.On("Clone", suite.tempDir, url).Return(&git.Repository{}, nil)

	// Get the parent and base dir
	parentDir := filepath.Dir(suite.tempDir)
	baseDir := filepath.Dir(parentDir)
	// Remove the temporary directory
	err := os.RemoveAll(suite.tempDir)
	if err != nil {
		log.Println("Error deleting directory:", err)
	}

	// Call cloneOrOpenByURL to clone the repo into the basdir
	suite.repoHandler.cloneOrOpenByURL(url, baseDir)

	// Check that the clone function of the mock was called
	suite.mockGitHelper.AssertCalled(suite.T(), "Clone", suite.tempDir, url)
}

// TestCloneOrOpenByURL_Open tests the open part of the cloneOrOpenByURL function
func (suite *RepoModuleTestSuite) TestCloneOrOpenByURL_Open() {
	// Creat the GitHub URL based on the path
	subs := strings.Split(suite.tempDir, string(os.PathSeparator))
	url := "https://github.com/" + subs[len(subs)-2] + "/" + subs[len(subs)-1]

	// Set return value for open function of the mocked GitHelper
	suite.mockGitHelper.On("Open", suite.tempDir).Return(&git.Repository{}, nil)

	// Get the parent and base dir
	parentDir := filepath.Dir(suite.tempDir)
	baseDir := filepath.Dir(parentDir)

	// Call cloneOrOpenByURL to open the repo from the basdir
	suite.repoHandler.cloneOrOpenByURL(url, baseDir)

	// Check that the open function of the mock was called
	suite.mockGitHelper.AssertCalled(suite.T(), "Open", suite.tempDir)
}

// TestCloneOrOpenByURL_ErrorHandling_Open test the error handling of the cloneOrOpenByURL open function
func (suite *RepoModuleTestSuite) TestCloneOrOpenByURL_ErrorHandling_Open() {
	// Creat the GitHub URL based on the path
	subs := strings.Split(suite.tempDir, string(os.PathSeparator))
	url := "https://github.com/" + subs[len(subs)-2] + "/" + subs[len(subs)-1]

	// Set return value for open function of the mocked GitHelper to return an error
	suite.mockGitHelper.On("Open", suite.tempDir).Return(nil, errors.New("mocked Error"))

	// Get the parent and base dir
	parentDir := filepath.Dir(suite.tempDir)
	baseDir := filepath.Dir(parentDir)

	// Call cloneOrOpenByURL
	gotRepo := suite.repoHandler.cloneOrOpenByURL(url, baseDir)

	// Check that the repo is nil if an error occurred
	suite.Assertions.Nil(gotRepo, "Should return nil on error.")
	// Check that the open function of the mock was called
	suite.mockGitHelper.AssertCalled(suite.T(), "Open", suite.tempDir)
}

// TestCloneOrOpenByURL_ErrorHandling_Clone test the error handling of the cloneOrOpenByURL clone function
func (suite *RepoModuleTestSuite) TestCloneOrOpenByURL_ErrorHandling_Clone() {
	// Creat the GitHub URL based on the path
	subs := strings.Split(suite.tempDir, string(os.PathSeparator))
	url := "https://github.com/" + subs[len(subs)-2] + "/" + subs[len(subs)-1]

	// Set return value for clone function of the mocked GitHelper to return an error
	suite.mockGitHelper.On("Clone", suite.tempDir, url).Return(nil, errors.New("mocked Error"))

	// Get the parent and base dir
	parentDir := filepath.Dir(suite.tempDir)
	baseDir := filepath.Dir(parentDir)
	// Remove the temporary directory, so clone get called
	err := os.RemoveAll(suite.tempDir)
	if err != nil {
		log.Println("Error deleting directory:", err)
	}

	// Call cloneOrOpenByURL
	gotRepo := suite.repoHandler.cloneOrOpenByURL(url, baseDir)

	// Check that the repo is nil if an error occurred
	suite.Assertions.Nil(gotRepo, "Should return nil on error.")
	// Check that the clone function of the mock was called
	suite.mockGitHelper.AssertCalled(suite.T(), "Clone", suite.tempDir, url)
}

// TestGetPathOfRepository test if the correct path of a repository is returned
func (suite *RepoModuleTestSuite) TestGetPathOfRepository() {
	// Call GetPathOfRepository
	gotDir := suite.repoHandler.GetPathOfRepository(suite.repo)

	// Check that the got directory equals the actual directory of the repository
	suite.Assertions.Equal(suite.tempDir, gotDir, "Dirs should equal.")
	// Check that the got directory is not nil
	suite.NotNil(gotDir, "Path should not equal an empty string.")
}

// TestGetGitHubURLOfRepository tests if the correct GitHub URL of a repository is returned
func (suite *RepoModuleTestSuite) TestGetGitHubURLOfRepository() {
	// Construct the expected GitHub URL
	subs := strings.Split(suite.tempDir, string(os.PathSeparator))
	expectedURL := "https://github.com/" + subs[len(subs)-2] + "/" + subs[len(subs)-1]

	// Call GetGitHubURLOfRepository
	gotURL := suite.repoHandler.GetGitHubURLOfRepository(suite.repo)

	// Check that the expected GitHub URL equals the got GitHub URL
	suite.Assertions.Equal(expectedURL, gotURL, "URLs should be equal.")
}

// TestResetRepository tests if a repository gets reseted correctly to its initial state after cloning.
func (suite *RepoModuleTestSuite) TestResetRepository() {
	var initialCommit *object.Commit
	// Create some commits
	for i := 0; i < 3; i++ {
		_, errCommit := suite.worktree.Commit("Initial Commit.", &git.CommitOptions{
			AllowEmptyCommits: true,
			Author: &object.Signature{
				Name:  "John Doe",
				Email: "john@doe.org",
				When:  time.Now(),
			},
		})
		if errCommit != nil {
			log.Fatalln("Error creating commit for test repository:", errCommit.Error())
		}
		// Save reference to initial commit
		if i == 0 {
			initialCommit, _ = suite.repoHandler.GetHeadCommit(suite.repo)
		}
	}

	// Call ResetRepository
	suite.repoHandler.ResetRepository(suite.repo, initialCommit)

	// Get the initial commit
	gotCommit, errHead := suite.repo.Head()
	if errHead != nil {
		log.Fatalln("Error fetching Head for test repository:", errHead.Error())
	}
	// Get the hash of the initial commit
	gotHash := gotCommit.Hash().String()
	initialCommitHash := initialCommit.Hash.String()

	// Check that the hash values of the expected initial commit and the actual initial commit are the same
	suite.Assertions.Equal(initialCommitHash, gotHash, "Commits should be the same.")

}

// TestGetHeadCommit test if the correct HEAD commits is returned
func (suite *RepoModuleTestSuite) TestGetHeadCommit() {
	// Create expected HEAD commit
	commit, err := suite.worktree.Commit("Initial Commit.", &git.CommitOptions{
		AllowEmptyCommits: true,
		Author: &object.Signature{
			Name:  "John Doe",
			Email: "john@doe.org",
			When:  time.Now(),
		},
	})
	if err != nil {
		log.Fatalln("Error creating commit for test repository:", err.Error())
	}
	// Get the expected HASH value
	expectedHash := commit.String()
	// Call GetHeadCommit
	headCommit, _ := suite.repoHandler.GetHeadCommit(suite.repo)

	// Get the actual hash value
	gotHash := headCommit.Hash.String()

	// Check if the expected and actual hash values are the same
	suite.Assertions.Equal(expectedHash, gotHash, "Head commit should be same.")
}

// TestCheckout checks if the correct commits is checked out
func (suite *RepoModuleTestSuite) TestCheckout() {
	// Create some commits
	var commits []plumbing.Hash
	for i := 1; i < 3; i++ {
		commit, errCommit := suite.worktree.Commit("Commit Nr."+strconv.Itoa(i), &git.CommitOptions{
			AllowEmptyCommits: true,
			Author: &object.Signature{
				Name:  "John Doe",
				Email: "john@doe.org",
				When:  time.Now(),
			},
		})
		if errCommit != nil {
			log.Fatalln("Error creating commit for test repository:", errCommit.Error())
		}

		commits = append(commits, commit)
	}

	// Iterate over list of commits
	for _, commit := range commits {
		// Checkout the current commit
		suite.repoHandler.Checkout(suite.repo, &git.CheckoutOptions{
			Hash:  commit,
			Force: true,
		})

		current, errHead := suite.repo.Head()
		if errHead != nil {
			log.Fatalln("Error fetching Head for test repository:", errHead.Error())
		}
		// Get the hash value of the current commit
		gotHash := current.Hash().String()
		// Get the expected hash value
		expectedHash := commit.String()
		// Check if the expected and actual hash values are the same
		suite.Assertions.Equal(expectedHash, gotHash, "Commits should be the same.")
	}
}

// TestGetRemoteBranches checks if all branches of a repository are returned
func (suite *RepoModuleTestSuite) TestGetRemoteBranches() {
	// Create test references
	var testNames = []string{"refs/remotes/test", "refs/local/test"}
	for _, name := range testNames {
		newRef := plumbing.NewReferenceFromStrings(name, "testHash")
		err := suite.repo.Storer.SetReference(newRef)
		if err != nil {
			log.Fatalln("Error setting remote reference for test repository:", err)
		}
	}

	// Fetch remotes from repo
	remotes := suite.repoHandler.GetRemoteBranches(suite.repo)
	gotName := remotes[0].Name().String()
	expectedName := testNames[0]
	// Check that the expected and actual reference are the same
	suite.Assertions.Equal(expectedName, gotName, "Remotes should equal.")
	// Check if only one branch was found
	suite.Assertions.Len(remotes, 1, "Should contain one remote reference.")
	// Check that the actual value is no local reference
	suite.Assertions.NotEqualf(testNames[1], gotName, "Remote should not equal local reference.")
}

// TestGetCommitsForRepository checks if all commits of a repository are returned
func (suite *RepoModuleTestSuite) TestGetCommitsForRepository() {
	// Create commits
	var expectedCommitHashes []string
	for i := 1; i < 3; i++ {
		commit, errCommit := suite.worktree.Commit("Commit Nr."+strconv.Itoa(i), &git.CommitOptions{
			AllowEmptyCommits: true,
			Author: &object.Signature{
				Name:  "John Doe",
				Email: "john@doe.org",
				When:  time.Now(),
			},
		})
		if errCommit != nil {
			log.Fatalln("Error creating commit for test repository:", errCommit.Error())
		}
		expectedCommitHashes = append(expectedCommitHashes, commit.String())
	}

	// Call GetCommitsForBranch
	gotCommits := suite.repoHandler.GetCommitsForBranch(suite.repo)
	// Get hashes of retunred commits
	var gotCommitsHashes []string
	for _, commit := range gotCommits {
		gotCommitsHashes = append(gotCommitsHashes, commit.Hash.String())
	}

	// Sort expected and actual hashes
	sort.Strings(expectedCommitHashes)
	sort.Strings(gotCommitsHashes)

	// Count expected and actual hashes
	expectedNumberOfCommits := len(expectedCommitHashes)
	gotNumberOfCommits := len(gotCommitsHashes)

	// Check that the number of actual and expected hashes/commits are the same
	suite.Assertions.Equal(expectedNumberOfCommits, gotNumberOfCommits, "Number of commits should equal.")
	// Check that the expected hashes and actual hashes are the same.
	suite.Assertions.Equal(expectedCommitHashes, gotCommitsHashes, "Commits should equal.")
}

// TestDeleteRepository check if a repository gets removed correctly
func (suite *RepoModuleTestSuite) TestDeleteRepository() {
	// Call DeleteRepository
	suite.repoHandler.DeleteRepository(suite.repo)
	// Check that the folder is removed from the file system
	suite.Assertions.NoDirExists(suite.tempDir, "Repository should be deleted.")
}

// TestGetTemplateTasks checks that the correct templateTaks are returned for a repository and template
func (suite *RepoModuleTestSuite) TestGetTemplateTasks() {
	// Create multiple templates
	var templates []Analyzer.Template
	for i := 1; i < 4; i++ {
		var tType string
		if i%2 != 0 {
			tType = "Flat"
		} else {
			tType = "Deep"
		}

		template := Analyzer.Template{
			Name: "TestTemplate" + strconv.Itoa(i),
			Type: tType,
			Script: Analyzer.Script{
				Code: "git test {{File}} {{Hash}}",
			},
		}
		templates = append(templates, template)
	}

	// Create commits
	var commits []plumbing.Hash
	for i := 0; i < 3; i++ {
		commit, errCommit := suite.worktree.Commit("Initial Commit.", &git.CommitOptions{
			AllowEmptyCommits: true,
			Author: &object.Signature{
				Name:  "John Doe",
				Email: "john@doe.org",
				When:  time.Now(),
			},
		})
		if errCommit != nil {
			log.Fatalln("Error creating commit for test repository:", errCommit.Error())
		}
		commits = append(commits, commit)
	}

	// Create a branch
	headRef, err := suite.repo.Head()
	if err != nil {
		log.Fatalln("Error fetching test repo head:", err)
	}
	newRef := plumbing.NewReferenceFromStrings("refs/remotes/test", headRef.Hash().String())
	// Update the reference in the repository.
	err = suite.repo.Storer.SetReference(newRef)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Created expected templateTasks
	expectedTemplateTasks := []Analyzer.TemplateTask{
		{
			CommitHash: commits[2].String(),
			Templates:  templates,
		},
		{
			CommitHash: commits[1].String(),
			Templates:  []Analyzer.Template{templates[1]},
		},
		{
			CommitHash: commits[0].String(),
			Templates:  []Analyzer.Template{templates[1]},
		},
	}

	// Call GetTemplateTasks
	gotTemplateTasks := suite.repoHandler.GetTemplateTasks(suite.repo, templates)

	// Check that expected and actual templateTasks are the same
	suite.Assertions.Equal(expectedTemplateTasks, gotTemplateTasks, "TemplateTasks should equal.")
	// Check that the first templateTask contains three tempaltes
	suite.Len(gotTemplateTasks[0].Templates, 3, "Head commit should contain 3 templates")

}

// This functions runs the test suite add a 'go test' command
func TestRepoModuleTestSuite(t *testing.T) {
	suite.Run(t, new(RepoModuleTestSuite))
}
