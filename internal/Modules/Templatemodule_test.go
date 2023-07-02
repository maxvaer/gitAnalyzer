// Package Modules contains all business logic modules/components of the application.
package Modules

import (
	"GitAnalyzer/api/Analyzer"
	"GitAnalyzer/internal/mocks"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"
)

// Create test suite for the templateHandler module
type TemplateModuleTestSuite struct {
	suite.Suite
	tempDir           string
	mockFileHelper    *mocks.IFileHelper
	mockCommandHelper *mocks.ICommandHelper
	mockRepoHandler   *mocks.IRepoHelper
	templateHandler   *TemplateHandler
	repo              *git.Repository
	worktree          *git.Worktree
	commits           []*object.Commit
	wgResults         sync.WaitGroup
}

// SetupTest is run before every test of the test suite to initialize a clear state
func (suite *TemplateModuleTestSuite) SetupTest() {
	// Create a temporary directory
	suite.tempDir = suite.T().TempDir()
	// Create a mocked FileHelper
	suite.mockFileHelper = mocks.NewIFileHelper(suite.T())
	// Create a mocked CommandHelper
	suite.mockCommandHelper = mocks.NewICommandHelper(suite.T())
	// Create a mocked repoHandler
	suite.mockRepoHandler = mocks.NewIRepoHelper(suite.T())
	// Create a config
	config := Analyzer.Config{}
	// Create a templateHandler with the created mocks
	suite.templateHandler = NewTemplateHandlerWithMocks(suite.mockRepoHandler, suite.mockFileHelper, suite.mockCommandHelper, config)
	repo, err := git.PlainInit(suite.tempDir, false)
	if err != nil {
		log.Fatalln("Error initializing in test repository:", err.Error())
	}
	// Set the repository
	suite.repo = repo
	// Get and set the worktree of the git repository
	w, err := repo.Worktree()
	if err != nil {
		log.Fatalln("Error getting worktree for test repository:", err.Error())
	}
	suite.worktree = w
	// Add some commits to the repository
	for i := 0; i < 3; i++ {
		commit, errCommit := w.Commit("Commit:"+strconv.Itoa(i), &git.CommitOptions{
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
		// Add the current commit to the slice of commits
		suite.commits = append(suite.commits, &object.Commit{Hash: commit})
	}
}

// TestExecuteTemplate test if a template executions calls all commands and regexes successfully
func (suite *TemplateModuleTestSuite) TestExecuteTemplate() {
	// Create a pathToFile map
	pathsMap := make(map[string][]string)
	testFile := "001"
	pathsMap[suite.tempDir] = append(pathsMap[suite.tempDir], testFile)

	// Create a test template
	template := Analyzer.Template{
		Name: "TestTemplate 1",
		Script: Analyzer.Script{
			Code: "git test {{File}} {{Hash}}",
		},
	}

	var firstCommit = suite.commits[0]

	// Set expected values
	expectedCmd := "git test " + "./" + testFile + " " + firstCommit.Hash.String()
	expectedURL := "https://github.com/gitanalyzer/test"
	expectedOutput := "Mock Result"
	expectedTimeStamp := time.Now().Format("01-02-2006")
	expectedResult := []Analyzer.Result{
		{
			TemplateName: template.Name,
			URL:          expectedURL,
			Output:       expectedOutput,
			Path:         suite.tempDir,
			CommitHash:   firstCommit.Hash.String(),
			Timestamp:    expectedTimeStamp,
		},
	}

	// Set expected return values for mocks
	suite.mockRepoHandler.On("GetPathOfRepository", suite.repo).Return(suite.tempDir)
	suite.mockFileHelper.On("SearchFilesByRegex", suite.tempDir, template).Return(nil)
	suite.mockFileHelper.On("FindFilesForCommands", suite.tempDir, template).Return(pathsMap)
	suite.mockCommandHelper.On("RunCommand", expectedCmd, suite.tempDir, "").Return(expectedOutput)
	suite.mockRepoHandler.On("GetGitHubURLOfRepository", suite.repo).Return(expectedURL)

	// Call executeTemplate
	gotResult := suite.templateHandler.executeTemplate(template, suite.repo, firstCommit.Hash.String())

	// Check that runCommand has been called
	suite.mockCommandHelper.AssertCalled(suite.T(), "RunCommand", expectedCmd, suite.tempDir, "")

	// Check that the expected and actual results are the same
	suite.Assertions.Equal(expectedResult, gotResult, "Results should equal.")
}

// TestProcessOutput checks if the output gets set correctly
func (suite *TemplateModuleTestSuite) TestProcessOutput() {
	// Create test template
	template := Analyzer.Template{
		Name: "TestTemplate 1",
		Script: Analyzer.Script{
			Code: "git test {{File}} {{Hash}}",
		},
	}
	// Set expected values
	var firstCommit = suite.commits[0]
	expectedURL := "https://github.com/gitanalyzer/test"
	expectedOutput := "Mock Result"
	expectedTimeStamp := time.Now().Format("01-02-2006")
	expectedResult := Analyzer.Result{
		Path:         suite.tempDir,
		TemplateName: template.Name,
		URL:          expectedURL,
		Output:       expectedOutput,
		CommitHash:   firstCommit.Hash.String(),
		Timestamp:    expectedTimeStamp,
	}

	// Set the return values for mocks
	suite.mockRepoHandler.On("GetGitHubURLOfRepository", suite.repo).Return(expectedURL)

	// Call processOutput
	gotOutput := suite.templateHandler.processOutput(suite.tempDir, template, expectedOutput+"\n", suite.repo, firstCommit.Hash.String())

	// Check if the actual and expected values match
	suite.mockRepoHandler.AssertCalled(suite.T(), "GetGitHubURLOfRepository", suite.repo)
	suite.Assertions.Equal(expectedResult, gotOutput, "Outputs should equal")
}

// TestProcessResult Check if the results gets formatted correctly
func (suite *TemplateModuleTestSuite) TestProcessResult() {
	// Set expected values
	var firstCommit = suite.commits[0]
	expectedURL := "https://github.com/gitanalyzer/test"
	expectedTimeStamp := time.Now().Format("01-02-2006")
	expectedResult := Analyzer.Result{
		URL:        expectedURL,
		CommitHash: firstCommit.Hash.String(),
		Timestamp:  expectedTimeStamp,
	}

	// Set the return values for mocks
	suite.mockRepoHandler.On("GetGitHubURLOfRepository", suite.repo).Return(expectedURL)

	// Call processResult
	gotResult := suite.templateHandler.processResult(Analyzer.Result{}, suite.repo, firstCommit.Hash.String())

	// Check if the actual and expected values match
	suite.mockRepoHandler.AssertCalled(suite.T(), "GetGitHubURLOfRepository", suite.repo)
	suite.Assertions.Equal(expectedResult, gotResult, "Results should equal")
}

// TestPostProcess test if the post-processing is executed correctly.
func (suite *TemplateModuleTestSuite) TestPostProcess() {
	// Create test templates
	var templates []Analyzer.Template
	for i := 1; i < 4; i++ {
		template := Analyzer.Template{
			Name: "TestTemplate" + strconv.Itoa(i),
			Output: Analyzer.Output{
				Unique: i%2 != 0,
			},
		}
		templates = append(templates, template)
	}
	suite.templateHandler.templates = templates

	// Set the return values for mocks
	suite.mockFileHelper.On("GenerateUniqueCSV", "TestTemplate1")
	suite.mockFileHelper.On("GenerateUniqueCSV", "TestTemplate3")

	// Call postProcess
	suite.templateHandler.postProcess()

	// Check if the expected function was called
	suite.mockFileHelper.AssertNumberOfCalls(suite.T(), "GenerateUniqueCSV", 2)
}

// TestPrepareCommands checks if the commands get prepared correctly before beeing executed
func (suite *TemplateModuleTestSuite) TestPrepareCommands() {
	var firstCommit = suite.commits[0]
	// Create test template
	template := Analyzer.Template{
		Name: "TestTemplate 1",
		Script: Analyzer.Script{
			Code: "git test {{File}} {{Hash}}",
		},
	}

	// Create and set expected values
	dirs := []string{"dir1", "dir2"}
	files := []string{"file1", "file2"}
	pathAndFiles := make(map[string][]string)
	pathAndFiles[dirs[0]] = files
	pathAndFiles[dirs[1]] = []string{files[0]}
	expectedCommandMap := make(map[string][]string)
	expectedCommandMap[dirs[0]] = []string{"git test ./file1 " + firstCommit.Hash.String(),
		"git test ./file2 " + firstCommit.Hash.String()}
	expectedCommandMap[dirs[1]] = []string{"git test ./file1 " + firstCommit.Hash.String()}

	// Call prepareCommands
	gotCommandMap := suite.templateHandler.prepareCommands(template, firstCommit.Hash.String(), pathAndFiles)

	// Check if expected value and actual value match
	suite.Assertions.Equal(expectedCommandMap, gotCommandMap, "CommandMaps should equal.")
}

// TestCheckRequirements check if all requirements are tested correctly
func (suite *TemplateModuleTestSuite) TestCheckRequirements() {
	// Create a template with test requirements
	suite.templateHandler.templates = []Analyzer.Template{
		Analyzer.Template{
			Name: "TestTemplate 1",
			Requirements: Analyzer.Requirements{
				Tools: []string{"tool1", "tool2"},
				Pip:   []string{"pip1"},
				Npm:   []string{"npm1", "npm2"},
			},
		},
	}

	// Set return values for mocks
	suite.mockCommandHelper.On("CheckRequiredTools", "tool1").Return(true)
	suite.mockCommandHelper.On("CheckRequiredTools", "tool2").Return(true)
	suite.mockCommandHelper.On("CheckRequiredPipPackage", "pip1").Return(true)
	suite.mockCommandHelper.On("CheckRequiredNPMPackage", "npm1").Return(true)
	suite.mockCommandHelper.On("CheckRequiredNPMPackage", "npm2").Return(true)

	// Call CheckRequirements
	suite.templateHandler.CheckRequirements()

	// Check if the expected functions where called the expected number of times
	suite.mockCommandHelper.AssertNumberOfCalls(suite.T(), "CheckRequiredTools", 2)
	suite.mockCommandHelper.AssertNumberOfCalls(suite.T(), "CheckRequiredPipPackage", 1)
	suite.mockCommandHelper.AssertNumberOfCalls(suite.T(), "CheckRequiredNPMPackage", 2)
}

// TestFilterTemplates test if the loaded templates are filtered correctly
func (suite *TemplateModuleTestSuite) TestFilterTemplates() {
	// Create some templates
	var templates []Analyzer.Template
	for i := 1; i < 4; i++ {
		template := Analyzer.Template{
			Name: "TestTemplate" + strconv.Itoa(i),
			Tags: []string{"TestTag" + strconv.Itoa(i%2), "testLanguage"},
		}
		templates = append(templates, template)
	}
	suite.templateHandler.templates = templates

	// Filter templates by tag with no excludes
	gotFilteredTemplates := suite.templateHandler.FilterTemplates("mockTag, TestTag1", "", "")
	// Set expected templates
	expectedFilteredTemplates := []Analyzer.Template{templates[0], templates[2]}
	// Check that the expected and actual filtered templates are equal
	suite.Assertions.Equal(expectedFilteredTemplates, gotFilteredTemplates, "Filtered templates by tag should equal.")
	// Check the number of actual filtered templates are equal
	suite.Assertions.Len(gotFilteredTemplates, 2, "Should return 2 templates")

	// Filter templates by language with no excludes
	gotFilteredTemplates = suite.templateHandler.FilterTemplates("", "testLanguage", "")
	// Set expected templates
	expectedFilteredTemplates = templates
	// Check that the expected and actual filtered templates are equal
	suite.Assertions.Equal(expectedFilteredTemplates, gotFilteredTemplates, "Filtered templates by language should equal.")
	// Check the number of actual filtered templates are equal
	suite.Assertions.Len(gotFilteredTemplates, 3, "Should return 3 templates")

	// Filter templates by tag with excludes
	gotFilteredTemplates = suite.templateHandler.FilterTemplates("", "testLanguage", "TestTemplate1")
	// Set expected templates
	expectedFilteredTemplates = []Analyzer.Template{templates[1], templates[2]}
	// Check that the expected and actual filtered templates are equal
	suite.Assertions.Equal(expectedFilteredTemplates, gotFilteredTemplates, "Filtered templates by excluded name should equal.")
	// Check the number of actual filtered templates are equal
	suite.Assertions.Len(gotFilteredTemplates, 2, "Should return 2 templates")
}

// TestLoadTemplates check if all templates are loaded as expected
func (suite *TemplateModuleTestSuite) TestLoadTemplates() {
	// Create test template
	template := Analyzer.Template{
		Name: "TestTemplate 1",
		Script: Analyzer.Script{
			Code: "git test {{File}} {{Hash}}",
		},
	}

	// Serialize test template into file
	data, err := yaml.Marshal(&template)
	if err != nil {
		log.Fatalln("Error marshalling test yaml:", err.Error())
	}

	err = os.WriteFile(suite.tempDir+string(os.PathSeparator)+"testTemplate.yaml", data, 0)
	if err != nil {
		log.Fatal("Error creating test yaml file:", err)
	}

	// Load serialized template from file
	suite.templateHandler.LoadTemplates(suite.tempDir)
	gotTemplate := suite.templateHandler.templates[0]

	// Check that the loaded and created test template are equal
	suite.Assertions.Equal(template.Name, gotTemplate.Name, "Loaded template names should equal")
	suite.Assertions.Equal(template.Script.Code, gotTemplate.Script.Code, "Loaded template code should equal")
}

// This functions runs the test suite add a 'go test' command
func TestTemplateModuleTestSuite(t *testing.T) {
	suite.Run(t, new(TemplateModuleTestSuite))
}
