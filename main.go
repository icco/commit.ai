package main

import (
	"context"
	"fmt"
	"log"
	"os"

	git "github.com/go-git/go-git/v5"
	openai "github.com/sashabaranov/go-openai"
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
	apiKey := os.Getenv("OPENAI_KEY")
	client := openai.NewClient(apiKey)

	// Set the prompt for the completion
	prompt := fmt.Sprintf("Suggest 10 commit messages based on the following diff:\n\n%s\n\nCommit messages should:\n - follow conventional commits\n - message format should be: <type>[scope]: <description>\n\nexamples:\n - fix(authentication): add password regex pattern\n - feat(storage): add new test cases\n", diffString)

	// Generate a completion using the OpenAI API
	completion, err := client.CreateCompletion(
		context.Background(),
		openai.CompletionRequest{
			Model:     openai.GPT3TextDavinci003,
			Prompt:    prompt,
			MaxTokens: 150,
		},
	)
	if err != nil {
		log.Fatal("Failed to generate completion:", err)
	}

	// Retrieve the generated commit message
	commitMessages := completion.Choices

	// Print the diff and generated commit message
	fmt.Println("Git diff:")
	fmt.Println(diffString)
	fmt.Printf("Generated commit messages:\n%+v", commitMessages)
}
