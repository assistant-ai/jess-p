package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/user"
	"strings"

	"github.com/assistant-ai/llmchat-client/db"
	"github.com/assistant-ai/llmchat-client/gpt"
	prompttools "github.com/assistant-ai/prompt-tools"
)

func expandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		usr, err := user.Current()
		if err != nil {
			return "", err
		}
		return strings.Replace(path, "~", usr.HomeDir, 1), nil
	}
	return path, nil
}

func main() {
	var param string
	flag.StringVar(&param, "p", "", "prompt")
	flag.Parse()

	var buffer bytes.Buffer
	_, err := io.Copy(&buffer, os.Stdin)
	if err != nil {
		panic(err)
	}

	// Print the content read from stdin
	printedContent := buffer.String()
	// print(printedContent)

	models := gpt.GetLlmClientGptModels()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	modelName := "gpt4"
	gptModel := models[modelName]
	path, err := expandPath("~/.jess/open-ai.key")
	llmClient, err := gpt.NewGptClientFromFile(path, 3, (*gpt.GPTModel)(gptModel), db.RandomContextId, models[modelName].MaxTokens, nil)
	if llmClient == nil {
		fmt.Println("ololo")
		return
	}
	prompt, _ := prompttools.CreateInitialPrompt(`User will provide instruction and context where this instruction will be applied. Output will be send to a bash pipe so your output should NOT have any explanations or anything, jsut required operations.
	if you asked to update something most likely you need to show same data stracture but updated, in the same format. Do not show HOW to update data structure, update it and output the result.
	Here is how you might be used by user: cat file | you -p "prompt" >> result
	
	`).
		AddTextToPrompt("user instructions: " + param).
		StartOfAdditionalInformationSection().
		AddTextToPrompt("Piped input: " + printedContent).
		EndOfAdditionalInformationSection().
		GenerateFinalPrompt()
	answer, err := llmClient.SendRandomContextMessage(prompt)

	io.WriteString(os.Stdout, answer)
}
