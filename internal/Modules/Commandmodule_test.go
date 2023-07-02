// Package Modules contains all business logic modules/components of the application.
package Modules

import (
	"github.com/stretchr/testify/suite"
	"log"
	"os"
	"testing"
)

// Create test suite for the command module
type CommandModuleTestSuite struct {
	suite.Suite
	tempDir        string
	commandHandler CommandHandler
}

// SetupTest is run before every test of the test suite to initialize a clear state
func (suite *CommandModuleTestSuite) SetupTest() {
	// Create and set a temporary directory
	suite.tempDir = suite.T().TempDir()
	// Initialize and set the commandHandler
	suite.commandHandler = CommandHandler{}
}

// TestDeleteScriptFile test if a script file is deleted as expected
func (suite *CommandModuleTestSuite) TestDeleteScriptFile() {
	// Create script file
	path := suite.tempDir + string(os.PathSeparator) + "script.sh"
	file, err := os.Create(path)
	if err != nil {
		log.Fatalln("Error creating test files:", err)
	}
	err = file.Close()
	if err != nil {
		log.Fatalln("Error closing test file:", err)
	}

	// Call deleteScriptFile to delete teh file
	suite.commandHandler.deleteScriptFile(path)

	// Check if the file is not existing anymore
	suite.Assertions.NoFileExists(path, "Script file should be deleted.")
}

// TestWriteScriptFileWithCode test if a script file with code get created correctly
func (suite *CommandModuleTestSuite) TestWriteScriptFileWithCode() {
	// Initialize test code
	expectedCode := "gitanalyzer -h"

	// Call writeScriptFileWithCode
	path := suite.tempDir + string(os.PathSeparator) + "test.sh"
	suite.commandHandler.writeScriptFileWithCode(expectedCode, path)

	// Check if the file was created
	suite.Assertions.FileExists(path, "Test script file should exist.")

	// Open and read the file
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalln("Error reading test script file.")
	}

	// Parse the output as string
	gotCode := string(data)

	// Check if the expected and got code are the same.
	suite.Assertions.Equal(expectedCode, gotCode, "Code snippet should equal.")
}

// This functions runs the test suite add a 'go test' command
func TestCommandModuleTestSuite(t *testing.T) {
	suite.Run(t, new(CommandModuleTestSuite))
}
