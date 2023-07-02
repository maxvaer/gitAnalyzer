// Package Modules contains all business logic modules/components of the application.
package Modules

import (
	"GitAnalyzer/api/Analyzer"
	"GitAnalyzer/pkg/Utils"
	"bufio"
	"fmt"
	"github.com/gocarina/gocsv"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
)

// IFileHelper the interface is used to define which actions can be called on the fileHelper
//
//go:generate mockery --name IFileHelper
type IFileHelper interface {
	SearchFilesByRegex(rootDir string, template Analyzer.Template) (result []Analyzer.Result)
	FindFilesForCommands(rootDir string, template Analyzer.Template) map[string][]string
	GetResultCSVPath(templateName string) string
	GenerateUniqueCSV(templateName string)
	GetTasks(urlsCSVPath string, checkedCSVPath string) (tasks []Analyzer.Task)
	PrepareResultsFolder(path string)
	MarshalSingleResult(result Analyzer.Result)
	MarshalMultipleResults(results []Analyzer.Result)
	MarshalStat(stat Analyzer.Stat)
}

// The FileHandler struct is responsible to handle all actions
// regarding file management and manipulation
type FileHandler struct {
	Config Analyzer.Config
}

// getResultCSVFileForTemplate returns the result file for the provided template name.
func (fh *FileHandler) getResultCSVFileForTemplate(templateName string) *os.File {
	// Get the path of the result file
	path := fh.GetResultCSVPath(templateName)
	// Open the result file
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Println("Error Opening File: ", err)
		return nil
	}
	return file
}

// GetResultCSVPath returns the path to the result file for a template by its name.
func (fh *FileHandler) GetResultCSVPath(templateName string) string {
	// Format name
	name := strings.ToLower(templateName)
	// Construct path
	relativePath := fh.Config.ResultsDir + string(os.PathSeparator) + name + ".csv"
	// Get the absolut path
	absolutePath, err := filepath.Abs(relativePath)
	if err != nil {
		fmt.Println("Error resolving path for result CSV file.")
	}
	return absolutePath
}

// getCheckedReposFile returns the file containing all repos which are already checked
func (fh *FileHandler) getCheckedReposFile() *os.File {
	// Get path of checked.csv
	path := fh.GetResultCSVPath("checked")
	// Open the file
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Println("Error Opening File: ", err)
	}

	return file
}

// GenerateUniqueCSV Generates a CSV file containing unique outputs per repository for a given template.
// The file will be named <templateName>_unique.csv
func (fh *FileHandler) GenerateUniqueCSV(templateName string) {
	// Get path of original CSV result file
	originalPath := fh.GetResultCSVPath(templateName)
	// Check if the original result file exists
	if !Utils.FileExists(originalPath) {
		// Return if the original CSV file does not exist.
		return
	}
	// Check if the original CSV file has some content
	fileInfo, err := os.Stat(originalPath)
	if err != nil || fileInfo.Size() == 0 {
		//FileInfo could not be created or file is empty
		return
	}

	// Construct the path for the unique CSV file.
	uniquePath := strings.TrimSuffix(originalPath, ".csv") + "_unique.csv"
	uniqueCSVFile, err := os.Create(uniquePath)
	if err != nil {
		log.Fatalln("error trying to create unique file", err)
	}

	// Load original results
	results := fh.UnMarshallResults(templateName)
	if results == nil {
		return
	}

	// Slice which contains the unique results.
	var uniqueResults []Analyzer.Result

	// Get the unique results per repository/URL
	resultsPerURL := fh.uniqueResultsPerURL(results)

	// Iterate over the map of unique results
	for _, outputs := range resultsPerURL {
		for _, result := range outputs {
			// Add unique results to slice.
			uniqueResults = append(uniqueResults, result)
		}
	}

	// Write unique results into the CSV file.
	fh.MarshalResults(uniqueResults, uniqueCSVFile)
}

// uniqueResultsPerURL maps the given results to a map like: RepositoryURL => Output => Result
// this generates a unique map of results
func (fh *FileHandler) uniqueResultsPerURL(results []Analyzer.Result) map[string]map[string]Analyzer.Result {
	// Map RepositoryURL => Output => Result
	urlsToResult := make(map[string]map[string]Analyzer.Result)

	// Iterate over Results
	for _, result := range results {
		// Check if URL of result was already seen.
		if _, exists := urlsToResult[result.URL]; !exists {
			// If no URL was found, create a new Map Output => Result under the Map entry of the URL
			urlsToResult[result.URL] = make(map[string]Analyzer.Result)
		}
		// Check if the output was already found for the given repository URL
		if _, exists := urlsToResult[result.URL][result.Output]; !exists {
			// If the output was not found, add the current result under the output entry.
			urlsToResult[result.URL][result.Output] = result
		}
	}

	return urlsToResult
}

// PrepareResultsFolder checks if the provided result directory is created.
// If the directory is not present it will be created.
func (fh *FileHandler) PrepareResultsFolder(path string) {
	// Check if dir is exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// If the dir does not exist, create it
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Println("Error creating result dir:", err)
		}
	}
}

// MarshalSingleResult marshals/persists the given result into the corresponding csv file.
func (fh *FileHandler) MarshalSingleResult(result Analyzer.Result) {
	wrapSlice := []Analyzer.Result{result}
	resultCSVFile := fh.getResultCSVFileForTemplate(result.TemplateName)
	fh.MarshalResults(wrapSlice, resultCSVFile)
}

// MarshalMultipleResults marshals/persists the result slice into the corresponding csv file.
func (fh *FileHandler) MarshalMultipleResults(results []Analyzer.Result) {
	for _, result := range results {
		fh.MarshalSingleResult(result)
	}
}

// MarshalResults arshals/persists the result slice into the provided file
func (fh *FileHandler) MarshalResults(results []Analyzer.Result, resultCSVFile *os.File) {
	// Get fileInfo for size
	fileInfo, err := resultCSVFile.Stat()
	if err != nil {
		log.Fatalln("Error getting fileInfo!", err.Error())
		return
	}

	// Check if file has some content
	// If file has content, marshall without headers
	if fileInfo.Size() == 0 {
		err = gocsv.MarshalFile(results, resultCSVFile)
	} else {
		err = gocsv.MarshalWithoutHeaders(results, resultCSVFile)
	}
	if err != nil {
		log.Fatalln("Error marshalling result struct to file!", err.Error())
		return
	}
	defer resultCSVFile.Close()
}

// UnMarshallResults loads all results for the given template name
func (fh *FileHandler) UnMarshallResults(templateName string) []Analyzer.Result {
	var results []Analyzer.Result
	// Get the result file by template name
	resultsCSVFile := fh.getResultCSVFileForTemplate(templateName)
	if resultsCSVFile == nil {
		return nil
	}

	// Load the results from the file
	if err := gocsv.UnmarshalFile(resultsCSVFile, &results); err != nil {
		log.Fatalln("Error Unmarshalling Results:", err.Error())
		return nil
	}

	defer resultsCSVFile.Close()

	return results
}

// MarshalStat marshals/persists the given stats into the checked.csv file
func (fh *FileHandler) MarshalStat(stat Analyzer.Stat) {
	wrapSlice := []Analyzer.Stat{stat}
	// Get the checked.csv file
	checkedCSVFile := fh.getCheckedReposFile()
	fileInfo, err := checkedCSVFile.Stat()
	if err != nil {
		log.Fatalln("Error getting fileInfo!", err.Error())
		return
	}

	// Check if file has some content
	// If file has content, marshall without headers
	if fileInfo.Size() == 0 {
		err = gocsv.MarshalFile(&wrapSlice, checkedCSVFile)
	} else {
		err = gocsv.MarshalWithoutHeaders(&wrapSlice, checkedCSVFile)
	}
	if err != nil {
		log.Fatalln("Error marshalling stat struct to file!")
		return
	}
	defer checkedCSVFile.Close()
}

// SearchFilesByRegex This method applies all regular expressions of the given template in
// the given root directory. Any matches will be returned inside slice
// of Analyzer.Result.
func (fh *FileHandler) SearchFilesByRegex(rootDir string, template Analyzer.Template) (result []Analyzer.Result) {
	if len(template.Regex) == 0 {
		// Return nil if template has no regular expressions
		return nil
	}

	// Iterate over all regular expressions of the template
	for _, regex := range template.Regex {
		// Compile the regular expression
		expression := regexp.MustCompile(strings.TrimSpace(regex.Expression))
		// Save filtering parameters
		excludes := regex.Exclude
		fileEndings := regex.FileEndings
		fileNames := regex.Filename
		// Fetch the matching file paths
		filePaths, found := fh.getFilePaths(rootDir, fileNames, fileEndings, excludes)
		if !found {
			// Continue with next regex if no matching file was found.
			continue
		}

		// Iterate over matching file paths
		for _, path := range filePaths {
			// Open the current file
			file, err := os.Open(path)
			if err != nil {
				log.Printf("Error opening file %s: %v\n", path, err.Error())
				errClose := file.Close()
				if errClose != nil {
					log.Println("Error closing file:", errClose.Error())
				}
				continue
			}

			// Read the file line by line
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()

				// Check if line contains matches.
				// (-1) matches all occurrences in the line.
				matchesSlice := expression.FindAllStringSubmatch(line, -1)

				// Check if some matches were found
				if len(matchesSlice) == 0 {
					continue
				}
				var output string
				// Iterate over found matches
				for _, matches := range matchesSlice {
					// Get the correct match via the regex group
					match := matches[regex.Group]
					// Check if the match is a false positive
					falsePositiveFound := false
					for _, falsePositive := range regex.FalsePositives {
						if strings.Contains(match, falsePositive) {
							falsePositiveFound = true
						}

					}
					if falsePositiveFound {
						//Continue/skip if false positive is found in matched string.
						continue
					}

					// Add the result to the output
					if output == "" {
						output = output + match
					} else {
						output = output + " " + match
					}
				}
				if len(output) == 0 {
					// If no output was found, continue with the next regex
					continue
				}
				// Create new result and append it to the slice of results
				result = append(result, Analyzer.Result{TemplateName: template.Name, Path: file.Name(), Description: regex.Description, Output: output})

			}
			// Close the file
			errClose := file.Close()
			if errClose != nil {
				log.Println("Error closing file:", errClose.Error())
				return nil
			}
		}

	}
	// Return the slice of results
	return result
}

// FindFilesForCommands creates a unique map of path => files for files needed for the command.
// Based on the provided template, all files which match inside the provided root dir
// will be inserted into the map
func (fh *FileHandler) FindFilesForCommands(rootDir string, template Analyzer.Template) map[string][]string {
	// Initialize result map path => files
	result := make(map[string][]string)
	// Check if the provided template match is empty
	if reflect.DeepEqual(template.Match, Analyzer.Match{}) {
		// As no match was provided, run the command at the root dir
		result[rootDir] = []string{}
	}

	//keywords := template.Match.Keywords
	// TODO: Add Keywords
	// Find the file paths
	filePaths, found := fh.getFilePaths(rootDir, template.Match.Filename, template.Match.FileEndings, template.Match.Exclude)
	if !found {
		// Return nil if no wanted filenames where found in the repository
		return result
	}

	// Iterate over found paths
	for _, path := range filePaths {
		// Get fileInfo
		fInfo, errStat := os.Stat(path)
		if errStat != nil {
			log.Print("Error converting path to fileInfo:", errStat)
		}
		// Check if the file/path is not a folder
		if !fInfo.IsDir() {
			//Split into path + file
			dir, file := filepath.Split(path)
			result[dir] = append(result[dir], file)
		} else {
			result[path] = nil
		}
	}

	return result
}

// getFilePaths searches all paths inside the provided root dir which
// match one of the provided filters.
func (fh *FileHandler) getFilePaths(rootDir string, filenames []string,
	fileEndings []string, excludes []string) (filePaths []string, found bool) {
	found = false

	// Iterate over all paths inside the given root dir
	err := filepath.Walk(rootDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Check if path contains a wanted filename
			for _, filename := range filenames {
				if excludes != nil {
					for _, exclude := range excludes {
						// Check if excluded keyword is not found
						if !strings.Contains(path, exclude) {
							// Add filepath if filename is found
							if strings.Contains(path, filename) {
								filePaths = append(filePaths, path)
								found = true
							}
						}
					}
				} else {
					// Add filepath if filename is found
					if strings.Contains(path, filename) {
						filePaths = append(filePaths, path)
						found = true
					}
				}
			}

			// Check if path ends on a wanted fileEnding
			fileExtension := filepath.Ext(path)
			if excludes != nil {
				for _, exclude := range excludes {
					// Check if excluded keyword is not found
					if !strings.Contains(path, exclude) {
						// Add filepath if fileEnding is found
						if Utils.Contains(fileEndings, fileExtension) {
							filePaths = append(filePaths, path)
							found = true
						}
					}
				}
			} else {
				// Add filepath if fileEnding is found
				if Utils.Contains(fileEndings, fileExtension) {
					filePaths = append(filePaths, path)
					found = true
				}
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
	return filePaths, found
}

// GetTasks loads all tasks from the provided checked.csv (path) and removes the allready
// checked repos from the loaded repos of the urls.csv file (path)
// Returns the list of none checked repos/tasks.
func (fh *FileHandler) GetTasks(urlsCSVPath string, checkedCSVPath string) (tasks []Analyzer.Task) {
	// Open the file of all tasks (url.csv)
	tasksFile, err := os.Open(urlsCSVPath)
	if err != nil {
		log.Println("Error opening CSV file:", err.Error())
		return nil
	}

	// Load tasks from urls.csv
	if err = gocsv.UnmarshalWithoutHeaders(tasksFile, &tasks); err != nil {
		log.Fatalln("Error Unmarshalling Tasks:", err.Error())
		return nil
	}

	err = tasksFile.Close()
	if err != nil {
		log.Println("Error closing URL file!", err.Error())
		return nil
	}

	// If no checked.csv file exists, return all tasks
	if !Utils.FileExists(checkedCSVPath) {
		return tasks
	}

	// Load checked stats (done tasks)
	stats := fh.UnMarshallStats(checkedCSVPath)

	// No stats loaded
	if len(stats) == 0 {
		return tasks
	}

	// Remove checked stats from Tasks
	tasks = fh.filterTasksByStats(stats, tasks)

	if err != nil {
		log.Println("Error closing checked file!", err.Error())
		return nil
	}

	return tasks
}

// UnMarshallStats loads all stats (done tasks) from the provide file
// and returns them as a slice
func (fh *FileHandler) UnMarshallStats(checkedCSVPath string) []Analyzer.Stat {
	var stats []Analyzer.Stat
	// Open file (checked.csv)
	checkedCSVFile, err := os.Open(checkedCSVPath)
	// Load stats
	if err = gocsv.UnmarshalFile(checkedCSVFile, &stats); err != nil {
		log.Fatalln("Error Unmarshalling Stats:", err.Error())
		return nil
	}

	defer checkedCSVFile.Close()

	return stats
}

// filterTasksByStats removes all stats (done tasks) from the list of tasks which need to be scanned.
func (fh *FileHandler) filterTasksByStats(stats []Analyzer.Stat, tasks []Analyzer.Task) []Analyzer.Task {
	// Iterate over all stats
	for _, stat := range stats {
		statURL := stat.URL
		// Iterate over tasks
		for i, task := range tasks {
			// If URL of stat and task match, the repo was already checked.
			if statURL == task.URL {
				// Reslice tasks, to exclude the already checked repo
				tasks = append(tasks[:i], tasks[i+1:]...)
				// Continue with the next stat
				break
			}
		}
	}
	return tasks
}
