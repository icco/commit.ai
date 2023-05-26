package main

import (
	"fmt"
	"log"
	"os"

	git "github.com/go-git/go-git/v5"
	openai "github.com/openai/openai-go/v1"
)

func main() {
	// Open the Git repository in the current directory
	repo, err := git.PlainOpen(".")
	if err != nil {
		log.Fatal("Failed to open Git repository:", err)
	}

	// Get the HEAD reference of the repository
	ref, err := repo.Head()
	if err != nil {
		log.Fatal("Failed to retrieve HEAD reference:", err)
	}

	// Get the commit object for the HEAD reference
	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		log.Fatal("Failed to retrieve commit object:", err)
	}

	// Get the diff of the commit against its parent
	parent, err := commit.Parent(0)
	if err != nil {
		log.Fatal("Failed to retrieve parent commit:", err)
	}
	patch, err := parent.Patch(commit)
	if err != nil {
		log.Fatal("Failed to retrieve commit patch:", err)
	}

	// Get the diff content as a string
	diffString := patch.String()

	// Set up the OpenAI client
	apiKey := os.Getenv("OPENAI_TOKEN")
	client := openai.NewClient(apiKey)

	// Set the prompt for the completion
	prompt := "Suggest 10 commit messages based on the following diff:\n\n%s\n\nCommit messages should:\n - follow conventional commits\n - message format should be: <type>[scope]: <description>\n\nexamples:\n - fix(authentication): add password regex pattern\n - feat(storage): add new test cases\n"

	// Generate a completion using the OpenAI API
	completion, err := client.Completions.Create(prompt, nil)
	if err != nil {
		log.Fatal("Failed to generate completion:", err)
	}

	// Retrieve the generated commit message
	commitMessage := completion.Choices[0].Text

	// Print the diff and generated commit message
	fmt.Println("Git diff:")
	fmt.Println(diffString)
	fmt.Println("Generated commit message:", commitMessage)
}
