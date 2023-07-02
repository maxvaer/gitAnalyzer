// Package Modules contains all business logic modules/components of the application.
package Modules

import (
	"GitAnalyzer/api/Analyzer"
	"github.com/stretchr/testify/suite"
	"log"
	"os"
	"strings"
	"testing"
)

// Create test suite for the fileHandling module
type FileHandlingModuleTestSuite struct {
	suite.Suite
	tempDir     string
	fileHandler FileHandler
}

// SetupTest is run before every test of the test suite to initialize a clear state
func (suite *FileHandlingModuleTestSuite) SetupTest() {
	// Create and set a temporary directory
	suite.tempDir = suite.T().TempDir()
	// Initialize and set the fileHandler
	suite.fileHandler.Config = Analyzer.Config{
		// Set the temporary directory as the result dir
		ResultsDir: suite.tempDir,
	}
}

// TestFilterTasksByStats test if stats (done tasks) get correctly filtered from
// the list of tasks loaded from the urls.csv
func (suite *FileHandlingModuleTestSuite) TestFilterTasksByStats() {
	// Initialize stats
	stats := []Analyzer.Stat{
		{
			URL: "testURL1",
		},
		{
			URL: "testURL2",
		},
	}

	// Initialize tasks
	tasks := []Analyzer.Task{
		{
			URL: "testURL1",
		},
		{
			URL: "testURL3",
		},
	}

	// Set expected value
	expectedTasks := []Analyzer.Task{
		{
			URL: "testURL3",
		},
	}

	// Call filterTasksByStats
	gotTasks := suite.fileHandler.filterTasksByStats(stats, tasks)

	// Check if expected value equals the got value
	suite.Assertions.Equal(expectedTasks, gotTasks, "Filtered tasks should equal.")
}

// TestGetFilePaths_Filenames test the filename filter of the getFilePaths function.
func (suite *FileHandlingModuleTestSuite) TestGetFilePaths_Filenames() {
	// Initialize test filenames
	filenames := []string{"test.txt", "testexcludedfile.txt", "lorem.json"}

	// Create files with test filenames
	for _, filename := range filenames {
		file, err := os.Create(suite.tempDir + string(os.PathSeparator) + filename)
		if err != nil {
			log.Fatalln("Error creating test files:", err)
		}
		err = file.Close()
		if err != nil {
			log.Fatalln("Error closing test files:", err)
		}
	}

	// Set expected paths from filenames
	expectedFilePaths := []string{
		suite.tempDir + string(os.PathSeparator) + filenames[2],
		suite.tempDir + string(os.PathSeparator) + filenames[0],
	}

	// Call getFilePaths with filenames and excluded values
	gotFilePaths, gotFound := suite.fileHandler.getFilePaths(suite.tempDir, filenames, []string{}, []string{"exclude"})

	// Check if the expected values match with the got values
	suite.Assertions.Equal(expectedFilePaths, gotFilePaths, "Files paths should equal.")
	suite.Assertions.True(gotFound, "Found should be true")
}

// TestGetFilePaths_FileEndings test the file endings filter of the getFilePaths function.
func (suite *FileHandlingModuleTestSuite) TestGetFilePaths_FileEndings() {
	// Initialize test filenames
	filenames := []string{"test.txt", "testexcludedfile.txt", "lorem.json"}

	// Create files with test filenames
	for _, filename := range filenames {
		file, err := os.Create(suite.tempDir + string(os.PathSeparator) + filename)
		if err != nil {
			log.Fatalln("Error creating test files:", err)
		}
		err = file.Close()
		if err != nil {
			log.Fatalln("Error closing test files:", err)
		}
	}

	// Set expected paths
	expectedFilePaths := []string{
		suite.tempDir + string(os.PathSeparator) + filenames[0],
	}

	// Call getFilePaths with file endings and excluded values
	gotFilePaths, gotFound := suite.fileHandler.getFilePaths(suite.tempDir, []string{}, []string{".txt"}, []string{"exclude"})

	// Check if the expected values match with the got values
	suite.Assertions.Equal(expectedFilePaths, gotFilePaths, "Files paths should equal.")
	suite.Assertions.True(gotFound, "Found should be true")
}

// TestFindFilesForCommands test if the correct file paths are returned for the template command
func (suite *FileHandlingModuleTestSuite) TestFindFilesForCommands() {
	// Initialize test filenames
	filenames := []string{"test.json", "testexcludedfile.json", "lorem.json"}

	// Create files with test filenames
	for _, filename := range filenames {
		file, err := os.Create(suite.tempDir + string(os.PathSeparator) + filename)
		if err != nil {
			log.Fatalln("Error creating test files:", err)
		}
		err = file.Close()
		if err != nil {
			log.Fatalln("Error closing test files:", err)
		}
	}

	// Create test template
	template := Analyzer.Template{
		Name: "TestTemplate",
		Match: Analyzer.Match{
			// Set expected file endings and excluded values
			FileEndings: []string{".json"},
			Exclude:     []string{"exclude"},
		},
	}

	// Create expected map of path => files
	expectedFileMap := make(map[string][]string)
	expectedFileMap[suite.tempDir+string(os.PathSeparator)] = []string{filenames[2], filenames[0]}

	// Call findFilesForCommands
	gotFileMap := suite.fileHandler.FindFilesForCommands(suite.tempDir, template)

	// Check if expected values equals got values
	suite.Assertions.Equal(expectedFileMap, gotFileMap, "FileMaps should equal.")
}

// TestFindFilesForCommands_NoMatchingFiles test if the correct file paths are returned
// if no matching parameters are provided
func (suite *FileHandlingModuleTestSuite) TestFindFilesForCommands_NoMatchingFiles() {
	// Create test template
	template := Analyzer.Template{
		Name: "TestTemplate",
	}

	// Create expected map of path => files
	expectedFileMap := make(map[string][]string)
	expectedFileMap[suite.tempDir] = []string{}

	// Call findFilesForCommands
	gotFileMap := suite.fileHandler.FindFilesForCommands(suite.tempDir, template)

	// Check if expected values equals got values
	suite.Assertions.Equal(expectedFileMap, gotFileMap, "FileMaps should equal.")
}

// TestSearchFilesByRegex test if file contents get searched correctly via a provided regex
func (suite *FileHandlingModuleTestSuite) TestSearchFilesByRegex() {

	// Create a secret value
	secret := "S3cr3tK3y"
	// Create test data
	data := "API_KEY: " + secret + "\n"

	// Initialize test file path
	filePath := suite.tempDir + string(os.PathSeparator) + "test.env"
	// Create test file
	err := os.WriteFile(filePath, []byte(data), 0)
	if err != nil {
		log.Fatal("Error creating test checked file:", err)
	}

	// Create test template with regular expression
	name := "TestTemplate"
	description := "Test regex"
	template := Analyzer.Template{
		Name: name,
		Regex: []Analyzer.Regex{
			{
				FileEndings: []string{".env"},
				Expression:  "(?m)^API_KEY:\\s*(.*)$",
				Group:       1,
				Description: description,
			},
		},
	}

	// Create expected results from the function call
	expectedResults := []Analyzer.Result{
		{
			TemplateName: name,
			Path:         filePath,
			Description:  description,
			Output:       secret,
		},
	}

	// Call SearchFilesByRegex
	gotResults := suite.fileHandler.SearchFilesByRegex(suite.tempDir, template)

	// Check  got value matches the expected values.
	suite.Assertions.Equal(expectedResults, gotResults, "Results should equal.")
}

// TestPrepareResultsFolder checks that the result directory get created if not already existing
func (suite *FileHandlingModuleTestSuite) TestPrepareResultsFolder() {
	// Create path to result dir
	path := suite.tempDir + string(os.PathSeparator) + "results"

	// Call PrepareResultsFolder
	suite.fileHandler.PrepareResultsFolder(path)

	// Check if the result dir exists
	suite.Assertions.DirExists(path, "Results folder should exist.")

}

// TestUniqueResultsPerURL check that the result for one URL should be unique and not doubled.
func (suite *FileHandlingModuleTestSuite) TestUniqueResultsPerURL() {
	// Create multiple result for same URL
	var results []Analyzer.Result
	for i := 0; i < 6; i++ {
		var result Analyzer.Result
		if i%2 == 0 {
			result = Analyzer.Result{
				URL:    "testURL1",
				Output: "testoutputA",
			}
		} else {
			result = Analyzer.Result{
				URL:    "testURL2",
				Output: "testoutputB",
			}
		}
		results = append(results, result)
	}

	// Create expected result map URL => results
	var expectedGotResultMap = make(map[string]map[string]Analyzer.Result)
	expectedGotResultMap["testURL1"] = make(map[string]Analyzer.Result)
	expectedGotResultMap["testURL1"]["testoutputA"] = results[0]
	expectedGotResultMap["testURL2"] = make(map[string]Analyzer.Result)
	expectedGotResultMap["testURL2"]["testoutputB"] = results[1]

	// Call uniqueResultsPerURL to generate unique results per URL
	gotResultsMap := suite.fileHandler.uniqueResultsPerURL(results)

	// Check that the expected and got values are equal
	suite.Assertions.Equal(expectedGotResultMap, gotResultsMap, "Unique result maps should equal.")
}

// TestGetCheckedReposFile checks if the checked.csv file gets returned
func (suite *FileHandlingModuleTestSuite) TestGetCheckedReposFile() {
	// Set the expected path
	expectedPath := suite.tempDir + string(os.PathSeparator) + "checked.csv"

	// Call getCheckedReposFile to return the checked.csv
	gotFile := suite.fileHandler.getCheckedReposFile()
	// Get the name of the returned file
	gotFilePath := gotFile.Name()
	err := gotFile.Close()
	if err != nil {
		log.Fatalln("Error closing test file:", err)
	}

	// Check that the file is not nil
	suite.Assertions.NotNil(gotFile, "File should not be nil.")
	// Check that the expected path equals the got path
	suite.Assertions.Equal(expectedPath, gotFilePath, "Paths should equal.")
}

// TestGetResultCSVPath check if the correct path to the result file of a template name gets returned
func (suite *FileHandlingModuleTestSuite) TestGetResultCSVPath() {
	// Create a test template name
	templateName := "testTemplate"
	// Set the expected path
	expectedPath := suite.tempDir + string(os.PathSeparator) + strings.ToLower(templateName) + ".csv"

	// Call GetResultCSVPath to get the result path for the template name
	gotPath := suite.fileHandler.GetResultCSVPath(templateName)

	// Check that the got value equals the expected value
	suite.Assertions.Equal(expectedPath, gotPath, "Path to result file should equal.")
}

// TestGetResultCSVPath check if the correct result file based on a template name gets returned
func (suite *FileHandlingModuleTestSuite) TestGetResultCSVFileForTemplate() {
	// Create a test template name
	templateName := "testTemplate"
	// Set the expected path
	expectedPath := suite.tempDir + string(os.PathSeparator) + strings.ToLower(templateName) + ".csv"

	// Call getResultCSVFileForTemplate to return the result file
	gotFile := suite.fileHandler.getResultCSVFileForTemplate(templateName)
	gotPath := gotFile.Name()
	err := gotFile.Close()
	if err != nil {
		log.Fatalln("Error closing test file:", err)
	}

	// Check that got path equals the expected path
	suite.Assertions.Equal(expectedPath, gotPath, "Paths of files should equal.")
	// Check that the file is not nil
	suite.Assertions.NotNil(gotFile, "File should not be nil.")
}

// This functions runs the test suite add a 'go test' command
func TestFileHandlingModuleTestSuite(t *testing.T) {
	suite.Run(t, new(FileHandlingModuleTestSuite))
}
