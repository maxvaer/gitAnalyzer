// Package Modules contains all business logic modules/components of the application.
package Modules

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// ICommandHelper the interface is used to define which actions can be called on the commandHelper
//
//go:generate mockery --name ICommandHelper
type ICommandHelper interface {
	CheckRequiredTools(tool string) bool
	CheckRequiredPipPackage(packageName string) bool
	CheckRequiredNPMPackage(packageName string) bool
	RunCommand(command string, path string, language string) string
}

// The CommandHandler struct is responsible to handle all actions
// necessary for execution of commands provided via templates
type CommandHandler struct {
}

// runner is a function to execute a command at a given path and return the results
type runner func(command string, path string) string

// RunCommand is the facade function to run a command depending on the language at a given path and return the results.
func (ch *CommandHandler) RunCommand(command, path, language string) string {
	// Get the runner of the language
	run := ch.getRunnerForLanguage(language)
	// Run command and return result
	return run(command, path)
}

// CheckRequiredTools tests, if the provided tool was found inside the PATH of the local system
func (ch *CommandHandler) CheckRequiredTools(tool string) bool {
	// Use which to look up the tool inside the PATH
	command := "which " + tool
	// Execute command
	preparedCmd := exec.Command("bash", "-c", command)
	out, err := preparedCmd.Output()
	if err != nil {
		// If an error occurred, the tool was not found.
		return false
	}
	// Parse results as string
	output := string(out)
	// Check if the results contain the string "not found"
	return !strings.Contains(output, "not found")
}

// CheckRequiredPipPackage tests if a given Python PIP package can be found on the system.
func (ch *CommandHandler) CheckRequiredPipPackage(packageName string) bool {
	// Use pip list to look up if the tool exists
	command := "pip list | grep " + packageName
	// Execute command
	preparedCmd := exec.Command("bash", "-c", command)
	out, err := preparedCmd.Output()
	if err != nil {
		// If an error occurred, the tool was not found.
		return false
	}
	// Parse results as string
	output := string(out)
	// If an entry was found, the package exists
	return len(output) > 0
}

// CheckRequiredNPMPackage tests if a given NPM package exists on the system.
func (ch *CommandHandler) CheckRequiredNPMPackage(packageName string) bool {
	// Look up the global list of NPM packages for the desired one
	command := "npm list -g " + packageName
	// Execute command
	preparedCmd := exec.Command("bash", "-c", command)
	out, err := preparedCmd.Output()
	if err != nil {
		// If an error occurred, the tool was not found.
		return false
	}
	// Parse results as string
	output := string(out)
	// If an entry was found, the package exists
	return len(output) > 0
}

// getRunnerForLanguage matches the provided language to a fitting runner.
// If the provided language does not match an existing one, the default runner (cli) will be returned.
func (ch *CommandHandler) getRunnerForLanguage(language string) runner {
	switch language {
	// WIP: JavaScript runner
	/*case "js", "JavaScript":
	return func(command string, path string) string {
		ctx := v8.NewContext()
		val, err := ctx.RunScript("const add = (a, b) => a + b", "math.js")
		if err != nil {
			fmt.Printf("Failed to execute JavaScript command: %s", command)
			return ""
		}
		return val.String()
	}*/
	case "python":
		return func(command string, path string) string {
			// Set Python version
			pythonVersion := "python3"
			if !ch.CheckRequiredTools(pythonVersion) {
				pythonVersion = "python"
			}
			// Write the temporary script to the provided path
			ch.writeScriptFileWithCode(command, path+"/gitAnalyzerPythonScript.py")
			// Execute the script
			preparedCmd := exec.Command(pythonVersion, "./gitAnalyzerPythonScript.py")
			preparedCmd.Dir = path
			out, err := preparedCmd.Output()
			if err != nil {
				// If an error occurred output it
				fmt.Print("Failed to execute Python script:", command, err.Error())
				// Delete the temporary script
				ch.deleteScriptFile(path + "/gitAnalyzerPythonScript.py")
				return ""
			}
			// Delete the temporary script
			ch.deleteScriptFile(path + "/gitAnalyzerPythonScript.py")
			// Return result
			return string(out)
		}
	case "bash":
		return func(command string, path string) string {
			// Write the temporary script to the provided path
			ch.writeScriptFileWithCode(command, path+"/gitAnalyzerScriptFile.sh")
			// Execute the script
			preparedCmd := exec.Command("bash", "./gitAnalyzerScriptFile.sh")
			preparedCmd.Dir = path
			out, err := preparedCmd.Output()
			if err != nil {
				// If an error occurred output it
				fmt.Print("Failed to execute shell script:", command, err.Error())
				// Delete the temporary script
				ch.deleteScriptFile(path + "/gitAnalyzerPythonScript.py")
				return ""
			}
			// Delete the temporary script
			ch.deleteScriptFile(path + "/gitAnalyzerScriptFile.sh")
			// Return result
			return string(out)
		}
	case "cli":
		// Fallthrough, as cli is the default runner
		fallthrough
	default:
		return func(command string, path string) string {
			// Execute command via Bash
			preparedCmd := exec.Command("bash", "-c", command)
			preparedCmd.Dir = path
			out, err := preparedCmd.Output()
			if err != nil {
				// If an error occurred but the output still contains some values, return these values
				if len(out) > 0 {
					return string(out)
				}
				fmt.Print("Failed to execute cli command:", command, err.Error())
				return ""
			}
			// Return result
			return string(out)
		}
	}
}

// writeScriptFileWithCode writes the given code into a script file at the provided path.
func (ch *CommandHandler) writeScriptFileWithCode(code string, path string) {
	// Create script file at path
	f, err := os.Create(path)
	if err != nil {
		log.Println("Error creating script file:", err.Error())
	}
	defer f.Close()

	// Write code to script file
	_, err = f.WriteString(code)
	if err != nil {
		log.Println("Error writing script file:", err.Error())
	}
}

// deleteScriptFile removes the script file at the provided path
func (ch *CommandHandler) deleteScriptFile(path string) {
	err := os.Remove(path)
	if err != nil {
		log.Println("Error deleting script file:", err.Error())
	}
}
