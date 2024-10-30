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

package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/srikanthbhandary/github-editor/pkg/github"
)

// ANSI escape codes for colors
const (
	Reset      = "\033[0m"
	Green      = "\033[32m"
	Blue       = "\033[34m"
	BoldYellow = "\033[1;33m"
)

var (
	githubRepo    string
	branch        string
	file          string
	regEx         string
	valueToUpdate string
	token         string
	rootCmd       = &cobra.Command{
		Use:   "ghedit",
		Short: "ghe",
		Long:  `Use this command to edit the specif file in the github and commit it`,
		Run: func(cmd *cobra.Command, args []string) {
			gh, err := github.NewGitHubWrapper(token, "SrikanthBhandary")
			if err != nil {
				log.Fatal(err)
			}

			defer gh.DeleteRepo()

			// Clone the repository
			if err := gh.Clone(githubRepo); err != nil {
				log.Fatal(err)
			}

			// Optionally checkout a specific branch
			if err := gh.Checkout(branch); err != nil {
				log.Fatal(err)
			}

			// Print the clone path
			fmt.Printf("Cloned repository to: %s\n", gh.GetClonePath())
			err = gh.MakeRegexReplace(file, regEx, valueToUpdate) // `// linter:\d+`, "// linter:245555")
			if err != nil {
				log.Fatal(err)
			}

			// Commit changes
			err = gh.Commit("Updated comments in file.go")
			if err != nil {
				log.Fatal(err)
			}

			// Push changes
			err = gh.Push()
			if err != nil {
				log.Fatal(err)
			}

		},
	}
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&githubRepo, "repo", "r", "", "repository to clone")
	rootCmd.PersistentFlags().StringVarP(&branch, "branch", "b", "", "branch to checkout")
	rootCmd.PersistentFlags().StringVarP(&file, "file", "f", "", "fileToModify")
	rootCmd.PersistentFlags().StringVarP(&regEx, "regEx", "e", "", "regex to find and replace")
	rootCmd.PersistentFlags().StringVarP(&valueToUpdate, "val", "v", "", "value to update")
	rootCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "value to update")

}

func initConfig() {
	// Display configuration settings with colored labels

	ValidateConfig()

	log.Printf(Green + "Initializing the config" + Reset)

	// Print separator line
	log.Printf("%s%-12s%s", BoldYellow, "----------------------------------", Reset)

	// Print each configuration item with colour-coded labels

	log.Printf("%s%-12s%s: %-20s", Blue, "Github Repo", Reset, githubRepo)
	log.Printf("%s%-12s%s: %-20s", Blue, "Branch", Reset, branch)
	log.Printf("%s%-12s%s: %-20s", Blue, "File", Reset, file)
	log.Printf("%s%-12s%s: %-20s", Blue, "Regex", Reset, regEx)
	log.Printf("%s%-12s%s: %-20s", Blue, "Value", Reset, valueToUpdate)

	// Print separator line
	log.Printf("%s%-12s%s", BoldYellow, "----------------------------------", Reset)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Command execution failed: %v", err)
	}
}

// ValidateConfig ensures that the required configuration values are provided
func ValidateConfig() {

	// Check if required flags are provided
	if githubRepo == "" {
		log.Fatalf("Error: Github repository (--repo) is required.")
	}
	if branch == "" {
		log.Fatalf("Error: Branch (--branch) is required.")
	}
	if file == "" {
		log.Fatalf("Error: File (--file) to modify is required.")
	}
	if regEx == "" {
		log.Printf("Warning: No regex pattern provided (--regEx). Proceeding without regex substitution.")
	}
	if valueToUpdate == "" {
		log.Printf("Warning: No value to update provided (--val). Proceeding without replacing values.")
	}

	log.Println("Configuration validation completed successfully.")
}
