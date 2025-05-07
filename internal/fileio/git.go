// Package fileio provides functions for file input/output operations
package fileio

import (
	"fmt"
	"os/exec"

	"github.com/gumi-tsd/secret-env-manager/internal/formatting"
	"github.com/gumi-tsd/secret-env-manager/internal/functional"
)

// GitIgnoreStatus represents the status of a file in git ignore
type GitIgnoreStatus struct {
	FileName  string
	IsIgnored bool
	Error     error
}

// NewGitIgnoreStatus creates a new GitIgnoreStatus
func NewGitIgnoreStatus(fileName string, isIgnored bool, err error) GitIgnoreStatus {
	return GitIgnoreStatus{
		FileName:  fileName,
		IsIgnored: isIgnored,
		Error:     err,
	}
}

// WithError returns a new GitIgnoreStatus with the specified error
func (g GitIgnoreStatus) WithError(err error) GitIgnoreStatus {
	return GitIgnoreStatus{
		FileName:  g.FileName,
		IsIgnored: g.IsIgnored,
		Error:     err,
	}
}

// IsSuccess returns true if there is no error
func (g GitIgnoreStatus) IsSuccess() bool {
	return g.Error == nil
}

// AsResult converts GitIgnoreStatus to a Result
func (g GitIgnoreStatus) AsResult() functional.Result[bool] {
	if g.Error != nil {
		return functional.Failure[bool](g.Error)
	}
	return functional.Success(g.IsIgnored)
}

// CheckFileIgnoredByGit verifies if the output file is ignored by git and warns about security risks
func CheckFileIgnoredByGit(outputFileName string) error {
	result := CheckGitIgnoreStatus(outputFileName)

	if result.Error != nil {
		return result.Error
	}

	if !result.IsIgnored {
		// Generate the warning message using a pure function
		warning := FormatSecurityWarning(outputFileName)
		fmt.Println(warning)
		return fmt.Errorf("output file '%s' is not ignored by git", outputFileName)
	}

	return nil
}

// CheckGitIgnoreStatus checks if a file is ignored by git
func CheckGitIgnoreStatus(fileName string) GitIgnoreStatus {
	// If filename is empty, return not ignored
	if fileName == "" {
		return NewGitIgnoreStatus(fileName, false, nil)
	}

	// Get git ignore status through command execution
	return ExecuteGitCheckIgnore(fileName)
}

// ExecuteGitCheckIgnore runs the git check-ignore command
func ExecuteGitCheckIgnore(fileName string) GitIgnoreStatus {
	cmd := exec.Command("git", "check-ignore", fileName)
	output, err := cmd.CombinedOutput()

	if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() == 1 {
		// Exit code 1 means the file is NOT ignored
		return NewGitIgnoreStatus(fileName, false, nil)
	} else if err != nil {
		// Generate the error message using a pure function
		errorMsg := formatting.Warning("Error checking git ignore status: %s", err)
		fmt.Println(errorMsg)
		// Assume not ignored on error for safety
		return NewGitIgnoreStatus(fileName, false,
			fmt.Errorf("error checking git ignore status: %w", err))
	}

	// If output has content, the file is ignored
	return NewGitIgnoreStatus(fileName, len(output) > 0, nil)
}

// IsFileIgnored checks if the specified file is ignored by git
// Returns Result[bool] indicating whether the file is ignored by git
func IsFileIgnored(fileName string) functional.Result[bool] {
	status := CheckGitIgnoreStatus(fileName)
	return status.AsResult()
}

// DisplaySecurityWarning shows a detailed security warning message for non-git-ignored files
// Note: This is an I/O function with side effects
func DisplaySecurityWarning(fileName string) {
	// Generate the warning message using a pure function
	warning := FormatSecurityWarning(fileName)
	fmt.Println(warning)
}

// FormatSecurityWarning creates a security warning message for the given file
// Pure function: Always returns the same output for the same input, without side effects
func FormatSecurityWarning(fileName string) string {
	// Generate each section using pure functions
	titleSection := formatting.Error("* Security Warning")

	// Emphasize the filename
	fileSection := formatting.Error("The file \"%s\" contains sensitive information (e.g., access keys).", fileName)

	// Risk section
	riskSection := formatting.Error("** Managing this file with Git poses a SEVERE risk of information leakage! **")

	// Required actions
	actionTitle := formatting.Error("Action Required:")
	action1 := formatting.Error("1.  Open your '.gitignore' file.")
	action2 := formatting.Error("2.  Add a line containing \"%s\" to the file.", fileName)
	action3 := formatting.Error("3.  Commit the updated '.gitignore' file.")

	// Construct the warning message by concatenating strings
	return fmt.Sprintf("\n\t%s\n\n\t%s\n\n\t%s\n\n\t%s\n\t%s\n\t%s\n\t%s\n",
		titleSection, fileSection, riskSection, actionTitle, action1, action2, action3)
}
