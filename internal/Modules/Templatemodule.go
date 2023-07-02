// Package Modules contains all business logic modules/components of the application.
package Modules

import (
	"GitAnalyzer/api/Analyzer"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// The TemplateHandler struct is responsible to handle all actions
// associated with templates
type TemplateHandler struct {
	// The used Analyzer.Config containing the settings of the current scan
	Config Analyzer.Config
	// All loaded templates from the template directory
	templates []Analyzer.Template
	// The CommandHelper used to call all actions associated with the commands of a template
	CommandHelper ICommandHelper
	// RepoHelper is used to interact with a git repository
	RepoHelper IRepoHelper
	// FileHelper to interact with files from the local file system
	FileHelper IFileHelper
}

// NewTemplateHandler constructor using a IRepoHelper and a Analyzer.Config to initialize the new templateHandler struct
func NewTemplateHandler(repoHelper IRepoHelper, config Analyzer.Config) *TemplateHandler {
	templateHandler := &TemplateHandler{RepoHelper: repoHelper,
		FileHelper: &FileHandler{config}, CommandHelper: &CommandHandler{},
		Config: config}
	return templateHandler
}

// NewTemplateHandlerWithMocks is a constructor which can be called using mocks of IRepoHelper, IFileHelper,
// ICommandHelper and a Analyzer.Config during testing
func NewTemplateHandlerWithMocks(repoHelper IRepoHelper, fileHelper IFileHelper, commandHelper ICommandHelper,
	config Analyzer.Config) *TemplateHandler {
	templateHandler := &TemplateHandler{RepoHelper: repoHelper, FileHelper: fileHelper, CommandHelper: commandHelper,
		Config: config}
	return templateHandler
}

// LoadTemplates loads all templates which can be found in the given directory
func (th *TemplateHandler) LoadTemplates(path string) {
	// Walk over all files in the directory
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalln("Error loading templates:", err)
		}

		// Return if not a yaml file
		if !strings.Contains(path, ".yaml") {
			return nil
		}

		// Get filename
		filePath, _ := filepath.Abs(path)
		fileName := filepath.Base(filePath)
		//Don't load the base template
		if fileName == "template.yaml" {
			return nil
		}
		yamlFile, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatalln("Error reading template yaml:", err)
		}

		// Load template from file
		var template Analyzer.Template
		err = yaml.Unmarshal(yamlFile, &template)
		if err != nil {
			log.Fatalln("Error while unmarshalling yaml:", err.Error())
		}
		th.templates = append(th.templates, template)

		if th.Config.Verbose {
			fmt.Printf("Loaded Template: %s\n", info.Name())
		}
		return nil
	})
}

// FilterTemplates uses the provided keywords, language and excluded names, to filter the loaded templates
// and return a subset of them.
func (th *TemplateHandler) FilterTemplates(keywords, language, excludedTemplateNames string) (filteredTemplates []Analyzer.Template) {
	// If no keyword is provided, append misc templates and language specific templates
	if keywords == "" {
		keywords = "misc," + language
	}

	// Split keywords by , into words
	words := strings.Split(keywords, ",")
	// Split excluded names by ,
	excludedNames := strings.Split(excludedTemplateNames, ",")

	// Iterate over loaded templates
	for _, template := range th.templates {
		// Load and format the template name
		name := template.Name
		name = strings.TrimSpace(name)
		name = strings.ToLower(name)
		// Check if the current template is excluded
		excludedFound := false
		for _, excludedName := range excludedNames {
			excludedName = strings.TrimSpace(excludedName)
			excludedName = strings.ToLower(excludedName)
			if name == excludedName {
				// If names equal, template is excluded
				excludedFound = true
			}
		}
		if excludedFound {
			// If the template is excluded continue with the next template
			continue
		}
		// Iterate over tags of the template
		for _, tag := range template.Tags {
			// Format the current tag
			tag = strings.TrimSpace(tag)
			tag = strings.ToLower(tag)
			// Iterate over loaded keywords
			for _, word := range words {
				// Format the current keywords
				word = strings.TrimSpace(word)
				word = strings.ToLower(word)
				if word == tag {
					// If keyword and tag match, add the template to the filtered templates
					filteredTemplates = append(filteredTemplates, template)
				}
			}
		}
	}
	return filteredTemplates
}

// executeTemplate executes the given template on the given repository and provided commit hash
// A slice of Analyzer.Results is returned, containing all results from the given template
func (th *TemplateHandler) executeTemplate(template Analyzer.Template, repo *git.Repository, commitHash string) []Analyzer.Result {
	// Initialize result slice
	var results []Analyzer.Result

	// Get the path of the repository
	repoPath := th.RepoHelper.GetPathOfRepository(repo)

	if th.Config.Verbose {
		fmt.Println("Executing template:", template.Name, " for repository:", repoPath)
	}

	// Search trough files using regular expressions from the template
	regexResults := th.FileHelper.SearchFilesByRegex(repoPath, template)
	// Append results of the regex search
	for _, result := range regexResults {
		result = th.processResult(result, repo, commitHash)
		results = append(results, result)
	}

	// Get map path => Files []
	dirToFilesMap := th.FileHelper.FindFilesForCommands(repoPath, template)
	if dirToFilesMap == nil {
		return nil
	}
	// Generate map of command => path
	commandToPathMap := th.prepareCommands(template, commitHash, dirToFilesMap)

	// For each path slice of Commands
	for path := range commandToPathMap {
		// Iterate over all commands for a path
		for _, cmd := range commandToPathMap[path] {
			// Run the command
			output := th.CommandHelper.RunCommand(cmd, path, template.Script.Language)
			if output == "" {
				// If there is no output continue with the next command
				continue
			}
			if th.Config.Verbose {
				fmt.Println("Result found:", output)
			}
			// Process and format the output
			result := th.processOutput(path, template, output, repo, commitHash)
			// Append result to return slice
			results = append(results, result)
		}
	}

	return results
}

// CheckRequirements checks if all necessary tools are found to run the loaded templates
func (th *TemplateHandler) CheckRequirements() {
	// Iterate over all loaded templates
	for _, template := range th.templates {
		// Iterate over all needed tools of the current template
		for _, tool := range template.Requirements.Tools {
			// Check if the current tools is installed
			found := th.CommandHelper.CheckRequiredTools(tool)
			if !found {
				log.Fatalln("Required tool not found:", tool)
			}
		}
		// Iterate over all needed pip packages of the current template
		for _, pipPackage := range template.Requirements.Pip {
			// Check if the current pip package is installed
			found := th.CommandHelper.CheckRequiredPipPackage(pipPackage)
			if !found {
				log.Fatalln("Required pip package not found:", pipPackage)
			}
		}
		// Iterate over all needed npm packages of the current template
		for _, npmPackage := range template.Requirements.Npm {
			// Check if the current npm package is installed
			found := th.CommandHelper.CheckRequiredNPMPackage(npmPackage)
			if !found {
				log.Fatalln("Required npm package not found:", npmPackage)
			}
		}
	}
}

// prepareCommands prepares the commands of a given template by replacing the template variables
func (th *TemplateHandler) prepareCommands(template Analyzer.Template, commitHash string, pathAndFiles map[string][]string) map[string][]string {
	var result = make(map[string][]string)

	// Iterate over all paths
	for path := range pathAndFiles {
		// Get all files for the current path
		files := pathAndFiles[path]
		// Get the command from the template
		cmd := template.Script.Code
		// Replace HASH placeholder with current commit hash
		cmd = strings.Replace(cmd, "{{Hash}}", commitHash, -1)
		if files != nil && strings.Contains(cmd, "{{File}}") {
			cmdCopy := cmd
			// Replace {{file}} placeholder with path to file
			for _, file := range files {
				cmd = strings.Replace(cmdCopy, "{{File}}", "./"+file, -1)
				result[path] = append(result[path], cmd)
			}
		} else {
			result[path] = append(result[path], cmd)
		}
	}

	return result
}

// RunAllTemplates runs all steps necessary to execute the selected template by the given Analyzer.Task for
// the given repo. The cTasks is used to send updates of the progress.
func (th *TemplateHandler) RunAllTemplates(task Analyzer.Task, repo *git.Repository, cTasks chan<- Analyzer.Task) {
	// Get the timestamp of the start
	start := time.Now()
	// Filter the templates
	filteredTemplates := th.FilterTemplates(th.Config.Tags, task.Language, th.Config.Excluded)
	// Update state to running
	task.State = "running"
	cTasks <- task
	// Run the filtered templates for the repository
	results := th.runTemplatesForRepository(filteredTemplates, repo)
	// Get the finished timestamp
	elapsedTime := time.Since(start)
	// Round time to milliseconds
	elapsedTime = elapsedTime.Round(time.Millisecond)
	// Attach results and elapsed time to task object
	task.Results = results
	task.ElapsedTime = elapsedTime.String()
	// Update state to finished
	task.State = "finished"
	cTasks <- task
	if !th.Config.KeepData {
		th.RepoHelper.DeleteRepository(repo)
	}
}

// postProcess executes all postScript commands of the loaded templates
func (th *TemplateHandler) postProcess() {
	// Iterate over loaded templates
	for _, template := range th.templates {
		// Check if output must be unique
		if template.Output.Unique {
			th.FileHelper.GenerateUniqueCSV(template.Name)
		}
		// TODO: Run postScript
	}
}

// preScript executes all preScript commands of the loaded templates
func (th *TemplateHandler) preScript() {
	// Set path to preSCript folder
	preScriptPath := "pre_script"

	//Execute preScripts
	for _, template := range th.templates {
		//Delete folder if exists
		err := os.RemoveAll(preScriptPath)
		if err != nil {
			log.Fatalln("Error deleting prescript folder:", err)
			return
		}
		//Create folder
		os.MkdirAll(preScriptPath, os.ModePerm)
		// Load command
		cmd := template.PreScript.Code
		if cmd == "" {
			continue
		}
		// Load language of command
		language := template.PreScript.Language
		// Execute command
		th.CommandHelper.RunCommand(cmd, preScriptPath, language)
	}
}

// processResult formats the provided results and uses the provided repo and commitHash
// to enrich the results with a timestamp and GitHub URL
func (th *TemplateHandler) processResult(result Analyzer.Result, repo *git.Repository, commitHash string) Analyzer.Result {
	if result.URL == "" {
		url := th.RepoHelper.GetGitHubURLOfRepository(repo)
		result.URL = url
	}
	if result.Timestamp == "" {
		timeStamp := time.Now().Format("01-02-2006")
		result.Timestamp = timeStamp
	}
	if result.CommitHash == "" {
		result.CommitHash = commitHash
	}

	return result
}

// processOutput format the output of the executed commands
func (th *TemplateHandler) processOutput(path string, template Analyzer.Template, output string, repo *git.Repository, commitHash string) Analyzer.Result {
	// Format the timestamp
	timeStamp := time.Now().Format("01-02-2006")
	url := th.RepoHelper.GetGitHubURLOfRepository(repo)

	// Remove linebreaks, so they don't get included into the result CSV
	if strings.Contains(output, "\n") {
		output = strings.TrimSuffix(output, "\n")
	}

	// Return the final result
	return Analyzer.Result{TemplateName: template.Name, URL: url, CommitHash: commitHash,
		Timestamp: timeStamp, Path: path, Description: "", Output: output}
}

// runTemplatesForRepository executes the given templates on the given repository
func (th *TemplateHandler) runTemplatesForRepository(templates []Analyzer.Template, repo *git.Repository) []Analyzer.Result {
	if templates == nil {
		return nil
	}

	// Get the current HEAD commit
	headCommit, err := th.RepoHelper.GetHeadCommit(repo)
	if err != nil {
		return nil
	}

	var results []Analyzer.Result
	// Get the templateTasks
	templateTask := th.RepoHelper.GetTemplateTasks(repo, templates)
	// Iterate over all templateTasks
	for _, tTask := range templateTask {
		// Checkout the current commit
		errCheckout := th.RepoHelper.Checkout(repo, &git.CheckoutOptions{
			Hash:  plumbing.NewHash(tTask.CommitHash),
			Force: true, //Maybe remove for performance improvement https://github.com/go-git/go-git/issues/511
		})
		if errCheckout != nil {
			continue
		}
		// Execute template
		for _, template := range tTask.Templates {
			templateResults := th.executeTemplate(template, repo, tTask.CommitHash)
			results = append(results, templateResults...)
		}
	}
	// Reset the repository to the newest commit
	defer th.RepoHelper.ResetRepository(repo, headCommit)

	return results
}
