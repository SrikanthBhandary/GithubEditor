// Package cmd provides command-line utilities for interacting with GitHub.
//
// ghedit - A command-line tool to edit specific files in a GitHub repository
// and commit the changes.
//
// Copyright (c) 2024 srikanth.bhandary@gmail.com
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package github

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
)

type GitHubWrapper struct {
	ClonePath string
	Token     string
	Username  string
}

func NewGitHubWrapper(token, username string) (*GitHubWrapper, error) {
	tempDir, err := os.MkdirTemp("", "gh_clone_")
	log.Println("TEMP", tempDir)
	if err != nil {
		return &GitHubWrapper{}, err
	}

	return &GitHubWrapper{
		ClonePath: tempDir,
		Token:     token,
		Username:  username,
	}, nil
}

func (ghw *GitHubWrapper) Clone(repo string) error {
	repoWithToken := fmt.Sprintf("https://%s:%s@%s", ghw.Username, ghw.Token, repo)

	fmt.Println("repoWithToken", repoWithToken)
	// Clone the repository if it doesn't exist
	if _, err := os.Stat(repo); os.IsNotExist(err) {
		// Repository does not exist, clone it
		cmd := exec.Command("git", "clone", repoWithToken, ghw.ClonePath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to clone repository: %w", err)
		}
	}
	return nil
}

// Checkout checks out the specified branch in the cloned repository.
func (ghw *GitHubWrapper) Checkout(branch string) error {
	// Change directory to the cloned repository
	if err := os.Chdir(ghw.ClonePath); err != nil {
		return fmt.Errorf("failed to change directory to %s: %w", ghw.ClonePath, err)
	}

	// Checkout the specified branch
	cmd := exec.Command("git", "checkout", branch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to checkout branch %s: %w", branch, err)
	}
	return nil
}

// GetClonePath returns the path where the repository is cloned.
func (ghw *GitHubWrapper) GetClonePath() string {
	return ghw.ClonePath
}

// GetClonePath returns the path where the repository is cloned.
func (ghw *GitHubWrapper) DeleteRepo() {
	log.Println("removing the repo: ", ghw.ClonePath)
	if err := os.RemoveAll(ghw.ClonePath); err != nil {
		log.Fatalf("Failed to remove directory: %v", err)
	}
}

// MakeRegexReplace reads a file in the cloned repository, applies regex replacement, and writes back the changes.
func (ghw *GitHubWrapper) MakeRegexReplace(file string, regEx string, valueToUpdate string) error {
	// Step 1: Read the file
	filePath := fmt.Sprintf("%s/%s", ghw.ClonePath, file)
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file %s: %v", filePath, err)
	}

	// Step 2: Compile the regex pattern
	re, err := regexp.Compile(regEx)
	if err != nil {
		return fmt.Errorf("error compiling regex %s: %v", regEx, err)
	}

	// Step 3: Perform the replacement
	newContent := re.ReplaceAllString(string(content), valueToUpdate)

	// Step 4: Write the modified content back to the file
	err = ioutil.WriteFile(filePath, []byte(newContent), 0644)
	if err != nil {
		return fmt.Errorf("error writing to file %s: %v", filePath, err)
	}

	log.Printf("Successfully replaced occurrences of '%s' in file '%s'.", regEx, filePath)
	return nil
}

// Commit commits the changes with a specified commit message.
func (ghw *GitHubWrapper) Commit(message string) error {
	if err := os.Chdir(ghw.ClonePath); err != nil {
		return fmt.Errorf("failed to change directory to %s: %w", ghw.ClonePath, err)
	}

	// Stage all changes
	cmd := exec.Command("git", "add", ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stage changes: %w", err)
	}

	// Commit changes
	commitCmd := exec.Command("git", "commit", "-m", message)
	commitCmd.Stdout = os.Stdout
	commitCmd.Stderr = os.Stderr
	if err := commitCmd.Run(); err != nil {
		return fmt.Errorf("failed to commit changes: %w", err)
	}

	log.Printf("Successfully committed changes with message: %s", message)
	return nil
}

// Push pushes the committed changes to the remote repository.
func (ghw *GitHubWrapper) Push() error {
	if err := os.Chdir(ghw.ClonePath); err != nil {
		return fmt.Errorf("failed to change directory to %s: %w", ghw.ClonePath, err)
	}

	// Push changes to the remote repository
	cmd := exec.Command("git", "push", "origin", "HEAD")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to push changes: %w", err)
	}

	log.Println("Successfully pushed changes to the remote repository.")
	return nil
}
