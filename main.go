package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"

	openai "github.com/sashabaranov/go-openai"
)

func main() {
	ctx := context.Background()

	diffOutput, err := exec.Command("git", "diff").Output()
	if err != nil {
		log.Fatal("Failed to retrieve git diff:", err)
	}

	diffString := string(diffOutput)

	// Set up the OpenAI client
	apiKey := os.Getenv("OPENAI_KEY")
	client := openai.NewClient(apiKey)

	// Set the prompt for the completion
	prompt := fmt.Sprintf("Suggest 10 commit messages based on the following diff:\n\n%s\n\nCommit messages should:\n - follow conventional commits\n - message format should be: <type>[scope]: <description>\n\nexamples:\n - fix(authentication): add password regex pattern\n - feat(storage): add new test cases\n", diffString)

	// Generate a completion using the OpenAI API
	completion, err := client.CreateCompletion(
		ctx,
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
	commitMessages := completion.Choices[0].Text

	// Print the diff and generated commit message
	fmt.Println("Git diff:")
	fmt.Println(diffString)
	fmt.Printf("Generated commit messages:\n%+v", commitMessages)
}
